#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname ${BASH_SOURCE})/../..
CODEGEN_PKG=${CODEGEN_PKG:-$(cd ${SCRIPT_ROOT}; ls -d -1 ./vendor/k8s.io/code-generator 2>/dev/null || echo ${GOPATH}/src/k8s.io/code-generator)}

vendor/k8s.io/code-generator/generate-groups.sh all \
  github.com/lessor/lessor/pkg/client github.com/lessor/lessor/pkg/apis \
  lessor.io:v1 \
  --go-header-file ${SCRIPT_ROOT}/tools/codegen/go-header-file
