# Lessor [![CircleCI](https://circleci.com/gh/lessor/lessor/tree/master.svg?style=svg&circle-token=6df998c0f2085edbc4bfeaf38e5114f990204c36)](https://circleci.com/gh/lessor/lessor/tree/master)

<p align="center">
  <img src="./docs/images/gophers/gophers_working.png" width="400">
</p>

Lessor is a [Kubernetes Operator](https://coreos.com/blog/introducing-operators.html) for deploying, managing, and securing multi-tenant workloads.

- [Introduction](#introduction)
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

Lessor is a [Kubernetes Operator](https://coreos.com/blog/introducing-operators.html) which aims to help manage the lifecycle of multi-tenant workloads. This repo currently contains a very high-level `Tenant` Custom Resource as well as some Controller functionality. This codebase is mostly being used to experiment with various approaches to multi-tenancy on Kubernetes.

If you're looking to contribute to this project, check out the [GitHub Issues](https://github.com/lessor/lessor/issues) and join the [#wg-multitenancy](https://kubernetes.slack.com/messages/C8E6YA9S7/) channel on the Kubernetes Slack. You can get an invite to Kubernetes Slack [here](http://slack.k8s.io/).

## How Does It Work?

### `Tenant` Custom Resource

Each complete tenant in your environment is represented by the `Tenant` Kubernetes custom resource. See an [example Custom Resource](./examples/tenant.yaml) for a more complete example of the configurable attributes of a tenant.

The following is a minimal example:

```yaml
apiVersion: lessor.io/v1
kind: Tenant
metadata:
  name: acme-labs
  labels:
    name: acme-labs
spec:
  namespaces:
  - acme-labs
  - acme-labs-dev
  - acme-labs-skunkworks
```

### Controller

Lessor uses the [Operator](https://coreos.com/blog/introducing-operators.html) pattern to encode domain-specific operational knowledge into software. The Operator pattern describes using a Kubernetes [Custom Resource Definition](https://kubernetes.io/docs/concepts/api-extension/custom-resources/) and a [Controller](https://github.com/kubernetes/community/blob/master/contributors/devel/controllers.md) to provide a declarative configuration interface to a self-healing system.

See the [Developer Guide](./docs/developer-guide.md) for information on building the controller and see the [Getting Started Guide](./docs/getting-started.md) for information on binary distributions.
