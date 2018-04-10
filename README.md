# Lessor [![CircleCI](https://circleci.com/gh/lessor/lessor/tree/master.svg?style=svg&circle-token=6df998c0f2085edbc4bfeaf38e5114f990204c36)](https://circleci.com/gh/lessor/lessor/tree/master)

<p align="center">
  <img src="./docs/gophers/boxes.png" width="400">
</p>

Deploy, manage, and secure single-tenant applications on [Kubernetes](https://kubernetes.io/).

- [Introduction](#introduction)
- [Motivation](#motivation)
- [How Does It Work?](#how-does-it-work)

In addition, here are some other documents that may be helpful:

- [Getting Started](./docs/getting-started.md)
- [Developer Guide](./docs/developer-guide.md)

## Introduction

> **les·sor** <br>
> *noun*
>
> a person or company that leases a good or service to an entity according to an agreement

Lessor is an open platform for deploying, managing, and securing many instances of single-tenant applications on [Kubernetes](https://kubernetes.io/). Lessor allows you to proxy to and independently scale each tenant with network and data isolation by default. This approach makes application development simpler and more secure.

## Motivation

Companies that create products for other companies or teams often have to reason about how to deal with the tenancy of each team. There are generally two paths:

- Deploy one monolithic application that handles multi-tenant data isolation via application logic
- Deploy and proxy to many instances of smaller, more isolated single-tenant applications

When faced with these two options, most companies choose to build the multi-tenant monolith. While the second path results in simpler, more secure software, many single-tenant applications are much more difficult to operate and observe. Large multi-tenant monoliths, however, have a habit of becoming difficult to operate and observe as well though.

Lessor aims to make it easier to choose to deploy and proxy to many instances of a single-tenant application by providing tools, services, and libraries that are purpose-built for this kind of deployment strategy.

## How Does It Work?

Lessor uses the [Operator](https://coreos.com/blog/introducing-operators.html) pattern to encode domain-specific operational knowledge into software. The Operator pattern describes using a Kubernetes [Custom Resorce Definition](https://kubernetes.io/docs/concepts/api-extension/custom-resources/) and a [Controller](https://github.com/kubernetes/community/blob/master/contributors/devel/controllers.md) to provide a declarative configuration interface to a self-healing system.

While most Kubernetes Operators deal with the administration of a single service, Lessor aims to automate larger application and cluster SRE objectives such as:

- tenant provisioning
- high-velocity deployments
- external resource acquisition
- secret distribution

Each "tenant" in your environment is represented as a Kubernetes resource. Each Tenant resource contains the metadata that describes how to configure, deploy, and connect the microservices which make up the tenant. See one of the [example tenants](./examples/tenant.yaml) for an overview of the configurable attributes of a tenant.
