package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Tenant struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TenantSpec   `json:"spec"`
	Status TenantStatus `json:"status"`
}

type TenantSpec struct {
	Databases    *DatabasesReference `json:"databases"`
	Repos        []DeployableRepo    `json:"repos"`
	Organization string              `json:"organization"`
	Email        string              `json:"email"`
}

// TenantStatus is the status for a Tenant resource
type TenantStatus struct {
	AvailableReplicas int32 `json:"availableReplicas"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type TenantList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Tenant `json:"items"`
}

type DeployableRepo struct {
	Name      string              `json:"name"`
	Container *ContainerReference `json:"container"`
	Varz      *VarzReference      `json:"varz"`
}

type DatabasesReference struct {
	Postgres []*PostgresReference `json:"postgres"`
}

type PostgresReference struct {
	Name string `json:"name"`
}

type ContainerReference struct {
	Name    *string `json:"name"`
	Version *string `json:"version"`
}

type VarzReference struct {
	Org      *string `json:"org"`
	Repo     *string `json:"repo"`
	Ref      *string `json:"ref"`
	Varz     string  `json:"varz"`
	Template string  `json:"template"`
}
