# Getting Started

<p align="center">
  <img src="./images/gophers/kubernetes_word.png" width="600">
</p>

## Installation Steps

To install Lessor in a Kubernetes cluster, you'll need:

- `kubectl` configured with admin access to a Kubernetes cluster

Run the following from the root of the repository to install the latest Service Catalog and Lessor components:

```
kubectl apply -f ./lessor.yaml
```

## Verifying The Installation

Verify that the `lessor-controller` deployment is running and healthy in `lessor-system`:

```
kubectl get pods --namespace lessor-system
```

Verify that the `service-catalog-apiserver` and `service-catalog-controller-manager` deployments are running and healthy in `kube-catalog`:

```
kubectl get pods --namespace kube-catalog
```

## Deploy An Application

Try creating an example tenant:

```
kubectl apply -f ./examples/crd.yaml
```

Watch the components that make up the tenant start up:

```
kubectl get pods --namespace acme-labs
```

Expose the webserver locally from a pod in the `kuard` deployment:

```
kubectl port-forward --namespace acme-labs kuard-7f79b5c84d-q7ktv 8080:8080
```

You should see be able to navigate to http://localhost:8080 and see the following:

![kuard](./images/screenshots/kuard.png)

## Uninstalling

To delete all Lessor deployment and services, you can delete the resources you created earlier `lessor.yaml`:

```
kubectl delete -f ./lessor.yaml
```
