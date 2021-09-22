#!/bin/bash

set -euo pipefail

required_args="
DOCKER_BUILD_DIR
DOCKER_OUTPUT_DIR
TARGET_GOOS
TARGET_GOARCH
"

for env_var in ${required_args}; do
  if [ -z "$(echo ${env_var})" ]; then
    echo "Missing ${env_var}, required env vars:"
    echo ${required_args}
    exit 1
  fi
done

cd ${DOCKER_BUILD_DIR}

go env -w GO111MODULE=off

GOPATH=$PWD go get -u -v github.com/jteeuwen/go-bindata/...
GOPATH=$PWD bin/go-bindata -o src/go-cowsay/bindata.go src/go-cowsay/cows
GOPATH=$PWD GOOS=$TARGET_GOOS GOARCH=$TARGET_GOARCH \
  go build -o ${DOCKER_OUTPUT_DIR}/cowsay go-cowsay

