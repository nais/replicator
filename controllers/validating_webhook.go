package controllers

import (
	"context"
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	apierrors "k8s.io/apimachinery/pkg/api/errors"

	naisiov1 "nais/replicator/api/v1"
	"nais/replicator/internal/replicator"
	"nais/replicator/internal/template"

	v1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

//+kubebuilder:webhook:path=/validate-replicationconfig,mutating=false,failurePolicy=fail,sideEffects=None,groups=nais.io,resources=replicationconfigs,verbs=create;update,versions=v1,name=replicationconfig.nais.io,admissionReviewVersions=v1

type ReplicatorValidator struct {
	Client  client.Client
	decoder admission.Decoder
}

func NewReplicatorValidator(mgr ctrl.Manager) *ReplicatorValidator {
	return &ReplicatorValidator{Client: mgr.GetClient(), decoder: admission.NewDecoder(mgr.GetScheme())}
}

func (v *ReplicatorValidator) Handle(ctx context.Context, req admission.Request) admission.Response {
	rc := &naisiov1.ReplicationConfig{}

	log.Info("Validating ReplicationConfig...")
	err := v.decoder.Decode(req, rc)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}
	if err := v.validateReplicationConfig(rc); err != nil {
		return admission.Denied(err.Error())
	}

	return admission.Allowed("")
}

func (v *ReplicatorValidator) validateReplicationConfig(rc *naisiov1.ReplicationConfig) error {
	if len(rc.Spec.Resources) == 0 {
		return fmt.Errorf("no resources specified")
	}

	for _, resource := range rc.Spec.Resources {
		if resource.Template == "" {
			return fmt.Errorf("template is empty")
		}
		resource, err := template.RenderTemplate(replicator.TemplateValues{Values: map[string]string{}}, resource.Template, template.WithOption("missingkey=invalid"))
		if err != nil {
			return fmt.Errorf("failed to render template: %w", err)
		}
		if resource.GetKind() == "" {
			return fmt.Errorf("kind is empty")
		}
		if resource.GetAPIVersion() == "" {
			return fmt.Errorf("apiVersion is empty")
		}
		if resource.GetName() == "" {
			return fmt.Errorf("name is empty")
		}
	}

	if err := v.validateValuesExists(context.Background(), rc); err != nil {
		return err
	}

	return nil
}

func (v *ReplicatorValidator) validateValuesExists(ctx context.Context, rc *naisiov1.ReplicationConfig) error {
	for _, s := range rc.Spec.TemplateValues.Secrets {
		var secret v1.Secret
		err := v.Client.Get(ctx, client.ObjectKey{Name: s.Name, Namespace: os.Getenv("POD_NAMESPACE")}, &secret)
		if err == nil {
			continue
		}

		if apierrors.IsNotFound(err) {
			if s.Validate {
				return fmt.Errorf("values references non-existing secret '%s'", s.Name)
			}

			log.Debugf("secret '%s' not found; ignoring error...", s.Name)
			continue
		}

		return fmt.Errorf("getting secret '%s': %w", s.Name, err)
	}

	return nil
}
