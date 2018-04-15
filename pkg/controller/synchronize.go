package controller

import (
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

// resolveTenantState compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the tenant resource
// with the current status of the resource.
func (c *Controller) resolveTenantState(key string) error {
	tenant, ok, err := c.tenantForCacheKey(key)
	if err != nil {
		return errors.Wrap(err, "couldn't find tenant given cache key")
	}
	if !ok {
		level.Info(c.logger).Log("err", "attempted to process tenant but tenant no longer exists", "tenant", key)
	}

	if err := c.validateTenant(tenant); err != nil {
		// We choose to absorb the error here as the worker would requeue the
		// resource otherwise. Since the tenant is invalid, requeueing the
		// tenant won't fix this problem.
		level.Info(c.logger).Log("msg", "tenant is invalid", "err", err, "key", key)
		return nil
	}

	_, err = c.namespacesLister.Get(tenant.Namespace())
	switch {
	case apierrors.IsNotFound(err):
		if _, err := c.kubeClient.CoreV1().Namespaces().Create(tenant.NamespaceResource()); err != nil {
			return errors.Wrapf(err, "error creating namespace for tenant %s", tenant.Name)
		}
	case err != nil:
		errors.Wrapf(err, "error getting namespace for tenant %s", tenant.Name)
	default:
		// the namespace already exists
	}

	if err := c.rehydrateSecrets(c.templateNamespace, tenant.Name); err != nil {
		return errors.Wrap(err, "error creating secrets for tenant")
	}

	c.recorder.Event(tenant, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)

	return nil
}
