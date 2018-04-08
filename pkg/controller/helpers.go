package controller

import (
	"fmt"

	lessorv1 "github.com/lessor/lessor/pkg/apis/lessor.io/v1"
	"github.com/pkg/errors"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/tools/cache"
)

const (
	// ControllerAgentName is the event source component name when using the
	// Kubernetes Event Recorder
	ControllerAgentName = "lessor"

	// SuccessSynced is used as part of the Event 'reason' when a tenant is synced
	SuccessSynced = "Synced"

	// ErrResourceExists is used as part of the Event 'reason' when a tenant fails
	// to sync due to a Deployment of the same name already existing.
	ErrResourceExists = "ErrResourceExists"

	// MessageResourceExists is the message used for Events when a resource
	// fails to sync due to a Deployment already existing
	MessageResourceExists = "Resource %q already exists and is not managed by tenant"

	// MessageResourceSynced is the message used for an Event fired when a tenant
	// is synced successfully
	MessageResourceSynced = "Tenant synced successfully"
)

func (c *Controller) tenantForCacheKey(key string) (*lessorv1.Tenant, bool, error) {
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return nil, false, errors.Wrapf(err, "invalid resource key: %s", key)
	}

	// Get the Tenant resource with this namespace/name
	tenant, err := c.tenantsLister.Tenants(namespace).Get(name)
	if err != nil {
		switch {
		case apierrors.IsNotFound(err):
			// The Tenant resource may no longer exist, which may not be an error
			return nil, false, nil
		default:
			return nil, false, err
		}
	}

	return tenant, true, nil
}

func (c *Controller) validateTenant(tenant *lessorv1.Tenant) error {
	switch tenant.Name {
	case "":
		return errors.New("tenant name cannot be empty")
	case "shared", "edge", "kube-system", "kube-public":
		return fmt.Errorf("%s is a reserved tenant name", tenant.Name)
	}
	return nil
}
