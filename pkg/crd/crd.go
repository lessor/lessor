package crd

import (
	"github.com/imdario/mergo"
	lessorv1 "github.com/lessor/lessor/pkg/apis/lessor.io/v1"
	"github.com/pkg/errors"

	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextcs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateOrUpdateCRDs will attempt to create or update all Kubernetes
// Custom Resource Definitions that are used in Lessor
func CreateOrUpdateCRDs(clientset apiextcs.Interface) error {
	crds := []*apiextv1beta1.CustomResourceDefinition{
		&apiextv1beta1.CustomResourceDefinition{
			ObjectMeta: metav1.ObjectMeta{
				Name: "tenants.lessor.io",
			},
			Spec: apiextv1beta1.CustomResourceDefinitionSpec{
				Group:   lessorv1.SchemeGroupVersion.Group,
				Version: lessorv1.SchemeGroupVersion.Version,
				Scope:   apiextv1beta1.NamespaceScoped,
				Names: apiextv1beta1.CustomResourceDefinitionNames{
					Plural: "tenants",
					Kind:   "Tenant",
				},
			},
		},
	}

	for _, crd := range crds {
		existing, err := clientset.ApiextensionsV1beta1().CustomResourceDefinitions().Get(crd.ObjectMeta.Name, metav1.GetOptions{})
		switch {
		case apierrors.IsNotFound(err):
			if _, createErr := clientset.ApiextensionsV1beta1().CustomResourceDefinitions().Create(crd); err != nil {
				return errors.Wrap(createErr, "error creating CRD")
			}
		case err != nil:
			return errors.Wrap(err, "error getting crd")
		default:
			merged := existing.DeepCopy()
			mergo.Merge(&merged, crd)
			if _, updateErr := clientset.ApiextensionsV1beta1().CustomResourceDefinitions().Update(merged); err != nil {
				return errors.Wrap(updateErr, "errors updating CRD")
			}
		}
	}
	return nil
}
