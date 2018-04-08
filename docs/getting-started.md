# Getting Started

> *Note: Parts of this section are currently aspirational and may not work as described.*

## Installation Steps

To install Lessor in a Kubernetes cluster, you'll need:

- A Kubernetes cluster which supports Service objects of type LoadBalancer ([GKE](https://cloud.google.com/kubernetes-engine/), [Minikube](https://kubernetes.io/docs/getting-started-guides/minikube/#quickstart), etc.)
- `kubectl` configured with admin access to your cluster
- The `lessor` CLI tool (optional)

With your `kubectl` context configured to use your cluster, you can use the `lessor` CLI to "adopt" a cluster:

```
lessor adopt cluster
```

You can also create the required deployments via a hosted manifest:

```
curl -L https://lessor.io/k8s/adopt.yaml | kubectl apply -f -
```

## Verifying The Installation

First ensure the following Kubernetes services are deployed:

- lessor-ingress

```
kubectl get services -n lessor-system
```

Next, ensure the corresponding Kubernetes pods are deployed and all containers are up and running:

- lessor-ingress
- lessor-controller

```
kubectl get pods -n lessor-system
```

## Deploy An Application

Try creating an example tenant:

```
kubectl apply -f ./examples/tenant.yaml
```

## Uninstalling

To delete all tenants and all Lessor deployment and services, you can run the following:

```
lessor eject cluster
```

If you don't have access to the `lessor` biniary, you can run a one-time job via a hosted manifest:

```
curl -L https://lessor.io/k8s/eject.yaml | kubectl apply -f -
```
