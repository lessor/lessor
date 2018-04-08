# Developer Guide

## Initial Setup

Check out the repository to the appropriate location in your `$GOPATH`:

```
git clone git@github.com:lessor/lessor.git $GOPATH/src/github.com/lessor/lessor
```

## Managing Dependencies

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

## Running Tests

Use `go test` to run tests:

```
go test -cover -race -v ./...
```

## Builing The Code

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

## Running Locally

To run the Lessor controller locally against the Kubernetes API server that is your currently configured `kubectl` context, run the following:

```
lessor run controller --local --debug
```

## Generating Clientset

To generate the [clientset](https://github.com/kubernetes/community/blob/master/contributors/devel/generating-clientset.md) for the Lessor API types, run the following:

```
./vendor/k8s.io/code-generator/generate-groups.sh all \
  github.com/lessor/lessor/pkg/client \
  github.com/lessor/lessor/pkg/apis \
  lessor.io:v1 \
  --go-header-file /dev/null
```

## Downloading Containers

Lessor containers are published on [Google Container Registry](https://cloud.google.com/container-registry/). The Lessor containers are not public so you may need to run the following to configure local access:

```
gcloud config configurations activate <configuration>
gcloud docker --authorize-only
```

The `latest` tag is continuously built from the `master` via [Google Container Builder](https://cloud.google.com/container-builder/) and published on [Google Container Registry](https://cloud.google.com/container-registry/):

```
docker pull gcr.io/lessor-io/lessor:latest
```

Each commit to each branch of the Lessor repository (`git@github.com:lessor/lessor.git`) also builds a container with the following naming scheme:`gcr.io/lessor-io/lessor:branch-commitsha`. For example:

```
docker pull gcr.io/lessor-io/lessor:master-81ea9bf9c8672a3c07be338dd6e2e8fd10d6cfba
```

Development branches (and their containers) are usually deleted as soon as possible, but the master containers should stay around for at least a few releases.
