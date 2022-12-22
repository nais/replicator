package v1

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"nais/replicator/internal/template"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var webhookLog = logf.Log.WithName("replicatorconfiguration-resource")

func (r *ReplicatorConfiguration) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/validate-nais-io-v1-replicatorconfiguration,mutating=false,failurePolicy=fail,sideEffects=None,groups=nais.io,resources=replicatorconfigurations,verbs=create;update,versions=v1,name=replicatorconfiguration.nais.io,admissionReviewVersions=v1

var _ webhook.Validator = &ReplicatorConfiguration{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *ReplicatorConfiguration) ValidateCreate() error {
	webhookLog.Info("validate create", "name", r.Name)
	return r.validateReplicatorConfiguration()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *ReplicatorConfiguration) ValidateUpdate(old runtime.Object) error {
	webhookLog.Info("validate update", "name", r.Name)
	return r.validateReplicatorConfiguration()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *ReplicatorConfiguration) ValidateDelete() error {
	panic("no-op")
}

func (r *ReplicatorConfiguration) validateReplicatorConfiguration() error {
	if len(r.Spec.Resources) == 0 {
		return fmt.Errorf("no resources specified")
	}

	for _, resource := range r.Spec.Resources {
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
