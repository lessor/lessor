#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname "${BASH_SOURCE}")/../..
DIFFROOT="${SCRIPT_ROOT}/k8s"
TMP_DIFFROOT="${SCRIPT_ROOT}/_tmp/k8s"
_tmp="${SCRIPT_ROOT}/_tmp"

cleanup() {
    rm -rf "${_tmp}"
}
trap "cleanup" EXIT SIGINT

cleanup

mkdir -p "${TMP_DIFFROOT}"
cp -a "${DIFFROOT}"/* "${TMP_DIFFROOT}"

"${SCRIPT_ROOT}/tools/codegen/update-k8s-codegen.sh"
echo "diffing ${DIFFROOT} against freshly generated codegen"
ret=0
diff -Naupr "${DIFFROOT}" "${TMP_DIFFROOT}" || ret=$?
cp -a "${TMP_DIFFROOT}"/* "${DIFFROOT}"
if [[ $ret -eq 0 ]]
then
    echo "${DIFFROOT} up to date."
else
    echo "${DIFFROOT} is out of date. Please run make generate"
    exit 1
fi
