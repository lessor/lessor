# Lessor [![CircleCI](https://circleci.com/gh/lessor/lessor/tree/master.svg?style=svg&circle-token=6df998c0f2085edbc4bfeaf38e5114f990204c36)](https://circleci.com/gh/lessor/lessor/tree/master)

<p align="center">
  <img src="./docs/images/gophers/boxes.png" width="400">
</p>

Deploy, manage, and secure applications on [Kubernetes](https://kubernetes.io/).

- [Introduction](#introduction)
- [Use Cases](#use-cases)
- [How Does It Work?](#how-does-it-work)

In addition, here are some other documents that may be helpful:

- [Documentation](./docs/README.md)
- [Getting Started](./docs/getting-started.md)
- [Developer Guide](./docs/developer-guide.md)

## Introduction

> **les·sor** <br>
> *noun*
>
> a person or company that leases a good or service to an entity according to an agreement

Lessor is an open platform for deploying, managing, and securing many instances of tenanted applications on [Kubernetes](https://kubernetes.io/). Lessor allows you to independently operate and scale each tenant with network and data isolation by default. This approach makes application development simpler and more secure.

## Use Cases

### B2B SaaS

Companies that create products for other companies or teams often have to reason about how to deal with the tenancy of each team. There are generally two paths:

- Deploy one monolithic application that handles multi-tenant data isolation via application logic
- Deploy and proxy to many instances of smaller, more isolated single-tenant applications

When faced with these two options, most companies choose to build the multi-tenant monolith. While the second path results in simpler, more secure software, many single-tenant applications are much more difficult to operate and observe. Large multi-tenant monoliths, however, have a habit of becoming difficult to operate and observe as well though.

Lessor aims to make it easier to choose to deploy and proxy to many instances of a single-tenant application by providing tools, services, and libraries that are purpose-built for this kind of deployment strategy.

### Staging Environments

Often, during the development process, developers need a quick (but reliable) way to deply an instance of an application (with various versions of it's components). Perhaps this is apart of a CI system that auto-deploys every PR to a repo or perhaps you need to create a demo instance of an app to perform user research.

Lessor aims to make this process easier by providing a typed API that can be managed via source-controlled files and RBAC for deploying application instances.

## How Does It Work?

Lessor uses the [Operator](https://coreos.com/blog/introducing-operators.html) pattern to encode domain-specific operational knowledge into software. The Operator pattern describes using a Kubernetes [Custom Resorce Definition](https://kubernetes.io/docs/concepts/api-extension/custom-resources/) and a [Controller](https://github.com/kubernetes/community/blob/master/contributors/devel/controllers.md) to provide a declarative configuration interface to a self-healing system.

While most Kubernetes Operators deal with the administration of a single service, Lessor aims to automate larger application and cluster SRE objectives such as:

- provisioning
- high-velocity deployments
- external resource acquisition
- secret distribution

### `Tenant` Resource

Each complete application instance in your environment is represented by the "Tenant" Kubernetes custom resource. Each Tenant resource contains the metadata that describes how to configure, deploy, and connect the microservices which make up the tenant. See an [example CRD](./examples/crd.yaml) for a more complete example of the configurable attributes of a tenant.

The following is a more minimal tenant that shows how to use the Open Service Broker API to bind to a MySQL server in Azure and deploy a set of minimally templated Kubernetes resources.

```yaml
apiVersion: lessor.io/v1
kind: Tenant
metadata:
  name: acme-labs
  labels:
    name: acme-labs

spec:
  # External services can be bound to using the Open Service Broker API.
  #
  # Lessor allows you to define the Service Instance and will automatically create the
  # appropriate binding. The stateless applications should be aware of things like what
  # format secrets will be in when bound, etc.
  catalog:
    serviceInstances:
      - clusterServiceClassExternalName: azure-mysql
        clusterServicePlanExternalName: basic50
        parameters:
          location: eastus
          resourceGroup: demo

  # Deployable resources can be generated via simple templates. A number of
  # template formats are supported.
  apps:
    templates:
      - name: kuard-handlebars
        type: handlebars
        url: https://lessor.io/latest/examples/templates/kuard-handlebars.yaml
        values:
          image: gcr.io/kuar-demo/kuard-amd64:1
```
