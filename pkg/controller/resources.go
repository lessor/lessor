package controller

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	lessorv1 "github.com/lessor/lessor/pkg/apis/lessor.io/v1"
)

// generator is a type which can be used to create parameterized deployable
// Kubernetes resources.
type generator struct {
	tenant          *lessorv1.Tenant
	ownerReferences []metav1.OwnerReference
}

// NewGenerator accepts a tenant instance and returns a generator which can
// be used to generate deployable resources for that tenant.
func newGenerator(tenant *lessorv1.Tenant) *generator {
	return &generator{
		tenant:          tenant,
		ownerReferences: ownerReferencesForTenant(tenant),
	}
}

// OwnerReferencesForTenant accepts a teenant and returns a set of owner references
// that indicate the the resource is managed by a controller on behalf of the given
// tenant resource
func ownerReferencesForTenant(tenant *lessorv1.Tenant) []metav1.OwnerReference {
	return []metav1.OwnerReference{
		*metav1.NewControllerRef(tenant, schema.GroupVersionKind{
			Group:   lessorv1.SchemeGroupVersion.Group,
			Version: lessorv1.SchemeGroupVersion.Version,
			Kind:    "tenant",
		}),
	}
}

func (g *generator) Namespace() *corev1.Namespace {
	return &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:            g.tenant.Name,
			OwnerReferences: g.ownerReferences,
		},
	}
}
