package controller

import (
	"fmt"

	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
)

// runWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message on the
// workqueue.
func (c *Controller) runTenantWorker() {
	for {
		obj, shutdown := c.tenantWorkqueue.Get()

		if shutdown {
			continue
		}

		// We wrap this block in a func so we can defer c.tenantWorkqueue.Done.
		err := func(obj interface{}) error {
			// We call Done here so the workqueue knows we have finished
			// processing this item. We also must remember to call Forget if we
			// do not want this work item being re-queued. For example, we do
			// not call Forget if a transient error occurs, instead the item is
			// put back on the workqueue and attempted again after a back-off
			// period.
			defer c.tenantWorkqueue.Done(obj)
			var key string
			var ok bool
			// We expect strings to come off the workqueue. These are of the
			// form namespace/name. We do this as the delayed nature of the
			// workqueue means the items in the informer cache may actually be
			// more up to date that when the item was initially put onto the
			// workqueue.
			if key, ok = obj.(string); !ok {
				// As the item in the workqueue is actually invalid, we call
				// Forget here else we'd go into a loop of attempting to
				// process a work item that is invalid.
				c.tenantWorkqueue.Forget(obj)
				level.Info(c.logger).Log(
					"err", fmt.Sprintf("expected string in workqueue but got %#v", obj),
				)
				return nil
			}
			// Run resolveTenantState, passing it the namespace/name string of the tenant
			// resource to be resolved.
			if err := c.resolveTenantState(key); err != nil {
				return errors.Wrapf(err, "error resolving tenant state for key: %s", key)
			}
			// Finally, if no error occurs we Forget this item so it does not
			// get queued again until another change happens.
			c.tenantWorkqueue.Forget(obj)
			level.Debug(c.logger).Log("msg", "successfully synced Tenant resource", "key", key)
			return nil
		}(obj)

		if err != nil {
			level.Info(c.logger).Log(
				"err", err,
				"msg", "error in processing workqueue job",
			)
			continue
		}
	}
}

// enqueueTenant takes a tenant resource and converts it into a namespace/name
// string which is then put onto the work queue. This method should *not* be
// passed resources of any type other than tenant.
func (c *Controller) enqueueTenant(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		level.Info(c.logger).Log(
			"err", err,
			"msg", "getting cache key for object while enqueueing tenant",
		)
		return
	}
	c.tenantWorkqueue.AddRateLimited(key)
}

// handleTenantOwner will take any resource implementing metav1.Object and attempt
// to find the tenant resource that 'owns' it. It does this by looking at the
// objects metadata.ownerReferences field for an appropriate OwnerReference.
// If the object is "owned" by a Tenant resource, we will enqueue that tenant
// resource to be processed. If the object does not have an appropriate
// OwnerReference, it will simply be skipped.
func (c *Controller) handleTenantOwner(obj interface{}) {
	var object metav1.Object
	var ok bool
	if object, ok = obj.(metav1.Object); !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			level.Info(c.logger).Log("err", "error decoding object, invalid type")
			return
		}
		object, ok = tombstone.Obj.(metav1.Object)
		if !ok {
			level.Info(c.logger).Log("err", "error decoding object tombstone, invalid type")
			return
		}
		level.Debug(c.logger).Log("msg", "recovered deleted object from tombstone", "name", object.GetName())
	}
	level.Debug(c.logger).Log("msg", "processing object", "name", object.GetName())
	if ownerRef := metav1.GetControllerOf(object); ownerRef != nil {
		// If this object is not owned by a tenant, we should not do anything more
		// with it.
		if ownerRef.Kind != "Tenant" {
			return
		}

		tenant, err := c.tenantsLister.Tenants(object.GetNamespace()).Get(ownerRef.Name)
		if err != nil {
			level.Debug(c.logger).Log(
				"msg", "ignoring orphaned object",
				"object", object.GetSelfLink(),
				"tenant", ownerRef.Name,
			)
			return
		}

		c.enqueueTenant(tenant)
		return
	}
}
