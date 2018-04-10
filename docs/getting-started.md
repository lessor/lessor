# Getting Started

> *Note: Parts of this section are currently aspirational and may not work as described.*

## Installation Steps

To install Lessor in a Kubernetes cluster, you'll need:

- `kubectl` configured with admin access to a Kubernetes cluster

With your `kubectl` context configured to use your cluster, run the following from the root of the repository to install the Lessor components:

```
kubectl apply -f ./lessor.yaml
```

## Verifying The Installation

Ensure the Kubernetes pods for the following Deployments are deployed and all containers are up and running:

- `lessor-controller`

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

To delete all Lessor deployment and services, you can delete the resources you created earlier `lessor.yaml`:

```
kubectl delete -f ./lessor.yaml
```

To delete all tenants and the Custom Resource Definition, run:

```
kubectl delete crd tenants.lessor.io
```
