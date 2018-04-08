# Lessor [![CircleCI](https://circleci.com/gh/lessor/lessor/tree/master.svg?style=svg&circle-token=6df998c0f2085edbc4bfeaf38e5114f990204c36)](https://circleci.com/gh/lessor/lessor/tree/master)

Lessor is a set of tools for deploying, managing, and securing many instances of single-tenant applications on [Kubernetes](https://kubernetes.io/). Lessor allows you to proxy to and independently scale each tenant with network and data isolation by default. This approach makes application development simpler and more secure.

Lessor is a permissively licensed open source project that was created by [Kolide](https://github.com/kolide) to deploy, manage, observe, and secure the Kolide Cloud product.

## Motivation

Companies that create products for other companies or teams often have to reason about how to deal with the tenancy of each team. There are generally two paths:

- Deploy one monolithic application that handles multi-tenant data isolation via application logic
- Deploy and proxy to many instances of smaller, more isolated single-tenant applications

When faced with these two options, most companies choose to build the multi-tenant monolith. While the second path results in simpler, more secure software, many single-tenant applications are much more difficult to operate and observe. Large multi-tenant monoliths, however, have a habit of becoming difficult to operate and observe as well though.

Lessor aims to make it easier to choose to deploy and proxy to many instances of a single-tenant application by providing tools, services, and libraries that are purpose-built for this kind of deployment strategy.

## Downloads

### Binaries

To download the latest `lessor` binary, you can use `go get`:

```
go get -u github.com/lessor/lessor/cmd/lessor
```

### Containers

Lessor containers are published on [Google Container Registry](https://cloud.google.com/container-registry/). These containers are not public though, so you may need to run the following to configure local access via the `docker` CLI tool:

```
gcloud config configurations activate <configuration>
gcloud docker --authorize-only
```

#### Latest Build

The `latest` tag is continuously built from the `master` branch and published to [Google Container Builder](https://cloud.google.com/container-builder/):

```
docker pull gcr.io/lessor-io/lessor:latest
```

#### Development Builds

Each commit to each branch of the Lessor repository (`git@github.com:lessor/lessor.git`) also builds a container with the following naming scheme:`gcr.io/lessor-io/lessor:branch-commitsha`. For example:

```
docker pull gcr.io/lessor-io/lessor:master-81ea9bf9c8672a3c07be338dd6e2e8fd10d6cfba
```

Development branches (and their containers) are usually deleted as soon as possible, but the master containers should stay around for at least a few releases.
