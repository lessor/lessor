package controller

// resolveTenantState compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the tenant resource
// with the current status of the resource.
func (c *Controller) resolveTenantState(key string) error {
	return nil
}

// resolveControlPlaneState compares the actual state with the desired, and
// attempts to converge the two. It then updates the Status block of the
// control plane resource with the current status of the resource.
func (c *Controller) resolveControlPlaneState(key string) error {
	return nil
}
