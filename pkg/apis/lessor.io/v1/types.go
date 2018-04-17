package v1

import (
	servicecatalog "github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
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
	Namespace string      `json:"namespace"`
	Catalog   CatalogSpec `json:"catalog"`
	Apps      AppSpec     `json:"apps"`
}

type AppSpec struct {
	Templates []TemplateSpec `json:"templates"`
}

type TemplateSpec struct {
	Name   string            `json:"name"`
	Type   string            `json:"type"`
	Url    string            `json:"url"`
	Values map[string]string `json:"values"`
}

type CatalogSpec struct {
	ServiceInstances []*servicecatalog.ServiceInstanceSpec `json:"serviceInstances"`
}

func (t *Tenant) Namespace() string {
	if t.Spec.Namespace == "" {
		return t.Name
	}

	return t.Spec.Namespace
}

func (t *Tenant) NamespaceResource() *corev1.Namespace {
	return &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:            t.Namespace(),
			OwnerReferences: t.ownerReferences(),
		},
	}
}

func (t *Tenant) ownerReferences() []metav1.OwnerReference {
	return []metav1.OwnerReference{
		*metav1.NewControllerRef(t, schema.GroupVersionKind{
			Group:   SchemeGroupVersion.Group,
			Version: SchemeGroupVersion.Version,
			Kind:    "tenant",
		}),
	}
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
