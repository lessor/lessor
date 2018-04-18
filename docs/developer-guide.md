# Developer Guide

<p align="center">
  <img src="./images/gophers/gophers_working.png" width="400">
</p>

The developer guide has the following sections:

- [Setup](#setup)
- [Build](#build)
- [Test](#test)
- [Run](#run)
- [Develop](#develop)

## Setup

Check out the repository to the appropriate location in your `$GOPATH`:

```
git clone git@github.com:lessor/lessor.git $GOPATH/src/github.com/lessor/lessor
```

## Build

### Managing Dependencies

Lessor uses [Dep](https://github.com/golang/dep) to manage Go dependencies.

To install the latest build of `dep`, run the following:

```
go get -u github.com/golang/dep/cmd/dep
```

To download the project dependences, run:

```
dep ensure -vendor-only
```

If you've added new code that requires on a new dependency, you must run the following for your dependency to be added:

```
dep ensure
```

If you'd like to update all dependencies, run `dep ensure -update`. If you'd like to just update one dependency, you're out of luck because Dep doesn't support that.

### Building The Code

Use `go build` to build the code:

```
go build
```

This will produce a `lessor` binary in the root of the repository.

```
./lessor --help
```

You can use the Go toolchain to install the binary to `$GOPATH/bin/lessor`:

```
go install
```

This is a static binary, so feel free to use `cp` to install it to your desired location:

```
cp ./lessor /usr/local/bin/
```

## Test

Use `go test` to run tests:

```
go test -cover -race -v ./...
```

### Running a CI Build Locally

The CircleCI configuration for Lessor includes a number of lint and test steps. If you'd like to run a complete, representative CI build locally, download the `circleci` CLI tool. See the [official installation instructions](https://circleci.com/docs/2.0/local-cli/#installing-the-circleci-local-cli-on-macos-and-linux-distros) for download information.

Once you have the tool installed in your path, run the following from the root of the repository:

```
circleci build
```

## Run

### Running The Controller Locally

To run the Lessor controller locally against the Kubernetes API server that is your currently configured `kubectl` context, run the following:

```
# create the namespace, CRD, etc
kubectl apply -f ./examples/development.yaml

# run the controller locally
lessor run controller --local --debug
```

### Downloading and Running Containers

#### Authentication

Lessor containers are published on [Google Container Registry](https://cloud.google.com/container-registry/). The Lessor containers are not public so you will need to install the [`gcloud`](https://cloud.google.com/sdk/downloads) command-line tool and run the following to configure local access:

```
gcloud auth configure-docker
```

If you'd like to configure your Kubernetes cluster to use these containers, see the Heptio documentation on [how to pull from private registries with Kubernetes](http://docs.heptio.com/content/private-registries/pr-gcr.html).

#### Tags

The `latest` tag is continuously built from the `master` via [Google Container Builder](https://cloud.google.com/container-builder/) and published on [Google Container Registry](https://cloud.google.com/container-registry/):

```
docker pull gcr.io/lessor-io/lessor:latest
```

Each commit to each branch of the Lessor repository (`git@github.com:lessor/lessor.git`) also builds a container with the following naming scheme:`gcr.io/lessor-io/lessor:branch-commitsha`. For example:

```
docker pull gcr.io/lessor-io/lessor:master-81ea9bf9c8672a3c07be338dd6e2e8fd10d6cfba
```

Development branches (and their containers) are usually deleted as soon as possible, but the master containers should stay around for at least a few releases.

#### Running The Container

You can use the `docker` CLI to run the container:

```
docker run -it gcr.io/lessor-io/lessor lessor --help
```

## Develop

### Generating Clientset

To generate the [clientset](https://github.com/kubernetes/community/blob/master/contributors/devel/generating-clientset.md) for the Lessor API types, run the following:

```
./vendor/k8s.io/code-generator/generate-groups.sh all \
  github.com/lessor/lessor/pkg/client \
  github.com/lessor/lessor/pkg/apis \
  lessor.io:v1 \
  --go-header-file /dev/null
```

### Import Organization

When importing Go dependences in a file, they should be formatted in the following order, each section separated by a newline:

1. Standard Library Imports
2. Non-standard Library Imports (anything on GitHub, etc.)
3. Upstream Kubernetes Imports (`k8s.io/...`)

Consider the following example:

```go
import (
	"fmt"
	"net/http"

	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
)
```

### Generating Service Catalog Manifests

Service Catalog suggests using Helm to install the required components. See the [Service Catalog Documentation](https://github.com/kubernetes-incubator/service-catalog/blob/master/docs/install.md) for the latest information on installing Service Catalog. The security model of Helm does not effectively account for the fact that the Kubernetes cluster may be a hostile environment. To install Service Catalog, the installation instructions also suggest giving the Tiller (Helm 2's server-side component) even more privileges.

The [Helm 3 Security Considerations Design Doc](https://github.com/kubernetes-helm/community/blob/master/helm-v3/007-security.md) explains that since Helm 3 will be purely client-side, the users access control properties will be used to install components in a way that ties into Kubernetes RBAC. Because of this future, [Helm V2 Access Control is never going to exist](https://github.com/kubernetes/helm/issues/1918#issuecomment-378700824). For this reason, I do not think it is reasonable to suggest using Helm to install the Service Catalog components.

Until Service Catalog can reliably be installed via some other means, we should generate a Helm manifest via the `helm install --dry-run --debug` in a one-off cluster and include the resultant resources in our bundled manifest. The following is a set of commands to generate a manifest.

```bash
# install the Tiller into the kube-system namespace
helm init

# add the Service Catalog Helm repository
helm repo add svc-cat https://svc-catalog-charts.storage.googleapis.com

# do a dry-run to generate the yaml (along with other logs)
helm install svc-cat/catalog \
  --name service \
  --namespace kube-catalog \
  --dry-run --debug > examples/service-catalog.yaml

# edit out the parts you don't want, make sure namespaces are correct, etc
vim examples/service-catalog.yaml

# uninstall Helm
kubectl delete deployment tiller-deploy -n kube-system
kubectl delete service tiller-deploy -n kube-system
```
