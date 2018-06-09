# Hard Multi-Tenancy Design and Implementation Discussion

## Authors

- Mike Arpaia ([@marpaia](https://github.com/marpaia))

## Context

In [Hard Multi-Tenancy in Kubernetes](https://docs.google.com/document/d/1mNL5oCIqtVwXI9piTPMuGArdZH8CA2UFaxHtM5Myp6M/edit#), Jessie Frazelle wrote about a design for "hard multi-tenancy" and how it may work in Kubernetes. In this document, I’d like to agree on an objective design and set some initial development objectives.

David Oppenheimer said the following rather succinctly in a comment on Jessie’s doc:

> I think there are two ways to look at this architecture. One is that it's a way to provide strong isolation between tenants within a cluster, by nesting per-tenant clusters within a host cluster. (That's the way it's described here--start with a cluster and add isolation boundaries.)
>
> The other way to look at it is that you start with independent per-tenant clusters (the typical way people do strong isolation between tenants today) and ask how can you make that more resource-efficient. When you start from that direction I think you end up in a very similar place with respect to the control plane components, i.e. hosting the control planes for all of the clusters on a shared cluster (see for example https://kubernetes.io/blog/2018/05/17/gardener/). Of course that still leaves you with separate physical nodes, so it doesn't get you all the way to the third approach described here.

Jessie's doc outlines 3 levels of organization that have varying degrees of isolation and resource efficiency. For the sake of keeping things interesting, I'm going to assume that we've all decided that we should go with the third design (with some modifications based on comments).

The key modification is based on a comment by David, where he added that API server and system services should be able to take advantage of whatever isolation is being provided to other nodes.

The resultant distribution of workloads thus looks something like this:

```
+------------+ +------------+  +------------+ +-------------+  +------------+ +------------+
| Tenant Foo | | Tenant Bar |  | Tenant Foo | | Tenant Bar  |  | Tenant Foo | | Tenant Bar |
| API Server | |    Pod     |  |    Pod     | | API Server  |  |    Pod     | |    Pod     |
+------------+ +------------+  +------------+ +-------------+  +------------+ +------------+
+------------+ +------------+  +------------+ +-------------+  +------------+ +------------+
| Tenant Bar | | Tenant Bar |  | Tenant Bar | | Tenant Foo  |  | Tenant Bar | | Tenant Bar |
|    Pod     | |    Pod     |  |    Pod     | |    Pod      |  |    Pod     | |    Pod     |
+------------+ +------------+  +------------+ +-------------+  +------------+ +------------+
+------------+ +------------+  +------------+ +-------------+  +------------+ +------------+
| Tenant Foo | |   Ring0    |  | Tenant Bar | |    Ring0    |  | Tenant Foo | |    Ring0   |
|    Pod     | |  Kubelet   |  |    Pod     | |   Kubelet   |  |    Pod     | |   Kubelet  |
+------------+ +------------+  +------------+ +-------------+  +------------+ +------------+
+---------------------------+  +----------------------------+  +---------------------------+
|          Node 1           |  |          Node 2            |  |          Node 3           |
+---------------------------+  +----------------------------+  +---------------------------+
+------------------------------------------------------------------------------------------+
|                                     Ring0 Kubernetes                                     |
+------------------------------------------------------------------------------------------+
```

## Proposal

I propose that we consider implementing a controller which manages an API Server (including all required components) per "tenant" as described in the Context section above using the Kubernetes Operator pattern.

As with any "operator", the role of the CRD is to define objective state. So if we think that we should do the whole Kubernetes in Kubernetes thing, the role of the controller is then to:

- Automate the management of a Kubernetes cluster per tenant, running on Kubernetes
  - API Server, Etcd, etc
- Effectively configure authentication and authorization to the appropriate API server per-user, per-tenant

The main runloop of the Lessor controller can be found in [`pkg/controller/synchronize.go`](../../pkg/controller/synchronize.go) in the `resolveTenantState` function. Whenever a tenant resource is updated, this function processes the tenant resource and endeavors to converge desired and actual state.

```go
// resolveTenantState compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the tenant resource
// with the current status of the resource.
func (c *Controller) resolveTenantState(key string) error {
	tenant, ok, err := c.tenantForCacheKey(key)
	// check err and ok

	// given the tenant variable which contains the full Kubernetes resource,
	// do stuff like:
	// - ensure all API components are running properly
	// - ensure authnz is configured correctly
```

## Consequences

The biggest consequence of this decision as far as I can tell is the Kubelet. This design calls for a single "Ring0 Kubelet" on each Ring0 Node to be communicating with multiple API servers (n+1, where n is the number of tenants). This is challenging because, as enumerated by David on Jessie's doc, there are a whole mess of problems with how various parts of Kubernetes work today that would complicate a design like this.
