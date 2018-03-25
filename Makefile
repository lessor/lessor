PATH := $(GOPATH)/bin:$(PATH)
VERSION = $(shell git describe --tags --always --dirty)
BRANCH = $(shell git rev-parse --abbrev-ref HEAD)
REVISION = $(shell git rev-parse HEAD)
REVSHORT = $(shell git rev-parse --short HEAD)
USER = $(shell whoami)
GOVERSION = $(shell go version | awk '{print $$3}')
NOW	= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

KIT_VERSION = "\
	-X github.com/lessor/lessor/vendor/github.com/kolide/kit/version.appName=${APP_NAME} \
	-X github.com/lessor/lessor/vendor/github.com/kolide/kit/version.version=${VERSION} \
	-X github.com/lessor/lessor/vendor/github.com/kolide/kit/version.branch=${BRANCH} \
	-X github.com/lessor/lessor/vendor/github.com/kolide/kit/version.revision=${REVISION} \
	-X github.com/lessor/lessor/vendor/github.com/kolide/kit/version.buildDate=${NOW} \
	-X github.com/lessor/lessor/vendor/github.com/kolide/kit/version.buildUser=${USER} \
	-X github.com/lessor/lessor/vendor/github.com/kolide/kit/version.goVersion=${GOVERSION}"

DOCKER_IMAGE = gcr.io/lessor-io/lessor:${VERSION}

ifeq ($(shell uname), Darwin)
	SHELL := /bin/bash
endif

all: build

.pre-build:
	mkdir -p build

deps:
	go get -u github.com/golang/dep/cmd/dep
	dep ensure -vendor-only
	cd ./website && yarn

.PHONY: build
build: lessor

build-linux: lessor-linux

test:
	go test -cover -race -v $(shell go list ./... | grep -v /vendor/)

container: build-linux .pre-lessor
	docker build -t ${DOCKER_IMAGE} .

container-push:
	gcloud docker -- push ${DOCKER_IMAGE}

.pre-lessor:
	$(eval APP_NAME = lessor)

lessor: .pre-build .pre-lessor
	go build -i -o build/lessor -ldflags ${KIT_VERSION} ./cmd/lessor

lessor-linux: .pre-build .pre-lessor
	GOOS=linux go build -i -o build/lessor-linux-amd64 -ldflags ${KIT_VERSION} ./cmd/lessor

.PHONY: website
website:
	cd ./website && yarn start
