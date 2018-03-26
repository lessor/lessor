package fake

import (
	lessor_io_v1 "github.com/lessor/lessor/pkg/apis/lessor.io/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeTenants implements TenantInterface
type FakeTenants struct {
	Fake *FakeLessorV1
	ns   string
}

var tenantsResource = schema.GroupVersionResource{Group: "lessor.io", Version: "v1", Resource: "tenants"}

var tenantsKind = schema.GroupVersionKind{Group: "lessor.io", Version: "v1", Kind: "Tenant"}

// Get takes name of the tenant, and returns the corresponding tenant object, and an error if there is any.
func (c *FakeTenants) Get(name string, options v1.GetOptions) (result *lessor_io_v1.Tenant, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(tenantsResource, c.ns, name), &lessor_io_v1.Tenant{})

	if obj == nil {
		return nil, err
	}
	return obj.(*lessor_io_v1.Tenant), err
}

// List takes label and field selectors, and returns the list of Tenants that match those selectors.
func (c *FakeTenants) List(opts v1.ListOptions) (result *lessor_io_v1.TenantList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(tenantsResource, tenantsKind, c.ns, opts), &lessor_io_v1.TenantList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &lessor_io_v1.TenantList{}
	for _, item := range obj.(*lessor_io_v1.TenantList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested tenants.
func (c *FakeTenants) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(tenantsResource, c.ns, opts))

}

// Create takes the representation of a tenant and creates it.  Returns the server's representation of the tenant, and an error, if there is any.
func (c *FakeTenants) Create(tenant *lessor_io_v1.Tenant) (result *lessor_io_v1.Tenant, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(tenantsResource, c.ns, tenant), &lessor_io_v1.Tenant{})

	if obj == nil {
		return nil, err
	}
	return obj.(*lessor_io_v1.Tenant), err
}

// Update takes the representation of a tenant and updates it. Returns the server's representation of the tenant, and an error, if there is any.
func (c *FakeTenants) Update(tenant *lessor_io_v1.Tenant) (result *lessor_io_v1.Tenant, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(tenantsResource, c.ns, tenant), &lessor_io_v1.Tenant{})

	if obj == nil {
		return nil, err
	}
	return obj.(*lessor_io_v1.Tenant), err
}

// Delete takes name of the tenant and deletes it. Returns an error if one occurs.
func (c *FakeTenants) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(tenantsResource, c.ns, name), &lessor_io_v1.Tenant{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeTenants) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(tenantsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &lessor_io_v1.TenantList{})
	return err
}

// Patch applies the patch and returns the patched tenant.
func (c *FakeTenants) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *lessor_io_v1.Tenant, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(tenantsResource, c.ns, name, data, subresources...), &lessor_io_v1.Tenant{})

	if obj == nil {
		return nil, err
	}
	return obj.(*lessor_io_v1.Tenant), err
}
