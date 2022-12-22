package controllers

import (
	"context"
	"fmt"
	"net/http"

	naisiov1 "nais/replicator/api/v1"
	"nais/replicator/internal/template"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

//+kubebuilder:webhook:path=/validate-replicatorconfig,mutating=false,failurePolicy=fail,sideEffects=None,groups=nais.io,resources=replicatorconfigurations,verbs=create;update,versions=v1,name=replicatorconfiguration.nais.io,admissionReviewVersions=v1

type ReplicatorValidator struct {
	Client  client.Client
	decoder *admission.Decoder
}

func (v *ReplicatorValidator) Handle(ctx context.Context, req admission.Request) admission.Response {
	rc := &naisiov1.ReplicatorConfiguration{}

	println("Validating ReplicatorConfiguration...")
	err := v.decoder.Decode(req, rc)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}
	if err := v.validateReplicatorConfiguration(rc); err != nil {
		return admission.Denied(err.Error())
	}

	return admission.Allowed("")
}

// replicatorValidator implements admission.DecoderInjector.
// A decoder will be automatically injected.

// InjectDecoder injects the decoder.
func (v *ReplicatorValidator) InjectDecoder(d *admission.Decoder) error {
	v.decoder = d
	return nil
}

func (v *ReplicatorValidator) validateReplicatorConfiguration(rc *naisiov1.ReplicatorConfiguration) error {
	if len(rc.Spec.Resources) == 0 {
		return fmt.Errorf("no resources specified")
	}

	for _, resource := range rc.Spec.Resources {
		if resource.Template == "" {
			return fmt.Errorf("template is empty")
		}
		resource, err := template.RenderTemplate(map[string]string{}, resource.Template, template.WithOption("missingkey=invalid"))
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

	return nil
}
