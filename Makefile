PATH := $(GOPATH)/bin:$(PATH)

all: lessor

.pre-build:
	mkdir -p build

deps:
	go get -u github.com/golang/dep/cmd/dep
	dep ensure -vendor-only

test:
	go test -cover -race -v ./...

generate:
	./tools/codegen/update-k8s-codegen.sh

lessor: .pre-build
	go build -i -o build/lessor ./cmd/lessor
