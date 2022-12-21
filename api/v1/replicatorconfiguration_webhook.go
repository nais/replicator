/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var replicatorconfigurationlog = logf.Log.WithName("replicatorconfiguration-resource")

func (r *ReplicatorConfiguration) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/validate-nais-io-v1-replicatorconfiguration,mutating=false,failurePolicy=fail,sideEffects=None,groups=nais.io,resources=replicatorconfigurations,verbs=create;update,versions=v1,name=replicatorconfiguration.nais.io,admissionReviewVersions=v1

var _ webhook.Validator = &ReplicatorConfiguration{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *ReplicatorConfiguration) ValidateCreate() error {
	replicatorconfigurationlog.Info("validate create", "name", r.Name)

	return r.validateReplicationConfiguration()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *ReplicatorConfiguration) ValidateUpdate(old runtime.Object) error {
	replicatorconfigurationlog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *ReplicatorConfiguration) ValidateDelete() error {
	replicatorconfigurationlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}
func (r *ReplicatorConfiguration) validateReplicationConfiguration() error {
	fmt.Printf("Validating ReplicationConfiguration %v", r.Spec.Resources)
	if len(r.Spec.Resources) == 0 {
		return fmt.Errorf("no resources specified")
	}

	for _, resource := range r.Spec.Resources {
		if resource.Template == "" {
			return fmt.Errorf("template is empty")
		}
		//resource, err := parser.RenderTemplate(unstructured.Unstructured{}, resource.Template)
		//if err != nil {
		//	return fmt.Errorf("failed to render template: %w", err)
		//}
		//if resource.GetKind() == "" {
		//	return fmt.Errorf("kind is empty")
		//}
		//if resource.GetAPIVersion() == "" {
		//	return fmt.Errorf("apiVersion is empty")
		//}
		//if resource.GetName() == "" {
		//	return fmt.Errorf("name is empty")
		//}
	}

	if len(r.Spec.Values.Secrets) == 0 && len(r.Spec.Values.ConfigMaps) == 0 {
		return fmt.Errorf("no values specified")
	}

	return nil
}
