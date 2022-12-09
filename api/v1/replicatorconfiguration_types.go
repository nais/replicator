package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// ReplicatorConfigurationSpec defines the desired state of ReplicatorConfiguration
type ReplicatorConfigurationSpec struct {
	NamespaceSelector metav1.LabelSelector `json:"namespaceSelector,omitempty"`

	// +kubebuilder:pruning:PreserveUnknownFields
	Resources []unstructured.Unstructured `json:"resources,omitempty"`
}

// ReplicatorConfigurationStatus defines the observed state of ReplicatorConfiguration
type ReplicatorConfigurationStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ReplicatorConfiguration is the Schema for the replicatorconfigurations API
type ReplicatorConfiguration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ReplicatorConfigurationSpec   `json:"spec,omitempty"`
	Status ReplicatorConfigurationStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ReplicatorConfigurationList contains a list of ReplicatorConfiguration
type ReplicatorConfigurationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ReplicatorConfiguration `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ReplicatorConfiguration{}, &ReplicatorConfigurationList{})
}
