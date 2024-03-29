package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ReplicationConfigSpec defines the desired state of ReplicationConfig
type ReplicationConfigSpec struct {
	NamespaceSelector metav1.LabelSelector `json:"namespaceSelector,omitempty"`
	TemplateValues    TemplateValues       `json:"templateValues,omitempty"`
	Resources         []Resource           `json:"resources,omitempty"`
}

type Secret struct {
	Name string `json:"name,omitempty"`
	// Validate checks that the secret exists before the ReplicationConfig is accepted.
	// Setting this to false explicitly marks the secret as eventually consistent during reconciliation for retry.
	// +kubebuilder:default=true
	// +kubebuilder:validation:Optional
	Validate bool `json:"validate,omitempty"`
}

type Resource struct {
	Template string `json:"template,omitempty"`
}

// ReplicationConfigStatus defines the observed state of ReplicationConfig
type ReplicationConfigStatus struct {
	SynchronizationTimestamp metav1.Time `json:"synchronizationTimestamp,omitempty"`
	SynchronizationHash      string      `json:"synchronizationHash,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster,shortName=repconf

// ReplicationConfig is the Schema for the replicationconfigs API
type ReplicationConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ReplicationConfigSpec   `json:"spec,omitempty"`
	Status ReplicationConfigStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ReplicationConfigList contains a list of ReplicationConfig
type ReplicationConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ReplicationConfig `json:"items"`
}

type TemplateValues struct {
	Values    map[string]string `json:"values,omitempty"`
	Secrets   []Secret          `json:"secrets,omitempty"`
	Namespace Namespace         `json:"namespace,omitempty"`
}

type Namespace struct {
	Labels      []string `json:"labels,omitempty"`
	Annotations []string `json:"annotations,omitempty"`
}

func init() {
	SchemeBuilder.Register(&ReplicationConfig{}, &ReplicationConfigList{})
}
