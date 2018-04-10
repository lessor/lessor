# Getting Started

> *Note: Parts of this section are currently aspirational and may not work as described.*

## Installation Steps

To install Lessor in a Kubernetes cluster, you'll need:

- A Kubernetes cluster which supports Service objects of type LoadBalancer ([GKE](https://cloud.google.com/kubernetes-engine/), [Minikube](https://kubernetes.io/docs/getting-started-guides/minikube/#quickstart), etc.)
- `kubectl` configured with admin access to your cluster
- The `lessor` CLI tool (optional)

With your `kubectl` context configured to use your cluster, you can use the `lessor` CLI to "adopt" a cluster:

```
kubectl apply -f ./example/install.yaml
```

## Verifying The Installation

Ensure the Kubernetes pods for the following Deployments are deployed and all containers are up and running:

- lessor-controller

```
kubectl get pods -n lessor-system
```

## Deploy An Application

Try creating an example tenant:

```
kubectl apply -f ./examples/tenant.yaml
```

Watch the tenant start up:

```
kubectl get pods -n acme-labs
```

## Uninstalling

To delete all tenants and all Lessor deployment and services, you can run a one-time job via a hosted manifest:

```
kubectl delete -f ./examples/install.yaml
```
