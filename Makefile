PATH := $(GOPATH)/bin:$(PATH)
GO111MODULE=off

all: build

deps:
	go get -u github.com/golang/dep/cmd/dep
	dep ensure -vendor-only
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install

build:
	go build .

test:
	go test -cover -race -v ./...

lint:
	gometalinter --disable-all \
	--enable=vet \
	--enable=golint \
	--enable=misspell \
	--skip=client \
	--skip=apis \
	./pkg/...

generate:
	./vendor/k8s.io/code-generator/generate-groups.sh all \
		github.com/lessor/lessor/pkg/client \
		github.com/lessor/lessor/pkg/apis \
		lessor.io:v1 \
		--go-header-file /dev/null
