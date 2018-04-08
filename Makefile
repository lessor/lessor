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
	./vendor/k8s.io/code-generator/generate-groups.sh all \
		github.com/lessor/lessor/pkg/client github.com/lessor/lessor/pkg/apis \
		lessor.io:v1 \
		--go-header-file /dev/null

lessor: .pre-build
	go build -i -o build/lessor ./cmd/lessor
