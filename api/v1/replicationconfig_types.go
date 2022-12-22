package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ReplicationConfigSpec defines the desired state of ReplicationConfig
type ReplicationConfigSpec struct {
	NamespaceSelector metav1.LabelSelector `json:"namespaceSelector,omitempty"`
	Values            Values               `json:"values,omitempty"`
	Resources         []Resource           `json:"resources,omitempty"`
}

type Values struct {
	Secrets    []ConfigResource `json:"secrets,omitempty"`
	ConfigMaps []ConfigResource `json:"configMaps,omitempty"`
}

type ConfigResource struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

type Resource struct {
	Template string `json:"template,omitempty"`
}

// ReplicationConfigStatus defines the observed state of ReplicationConfig
type ReplicationConfigStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster,shortName=repconf

// ReplicationConfig is the Schema for the replicatorconfigurations API
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

func init() {
	SchemeBuilder.Register(&ReplicationConfig{}, &ReplicationConfigList{})
}
