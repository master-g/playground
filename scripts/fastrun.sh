#!/usr/bin/env bash

echo "==> Building..."

PACKAGE=github.com/master-g/playground/cmd/playground
COMMIT_HASH=$(git rev-parse --short HEAD)
BUILD_DATE=$(date +%Y-%m-%dT%TZ%z)

LD_FLAGS="-X ${PACKAGE}/buildinfo.CommitHash=${COMMIT_HASH} -X ${PACKAGE}/buildinfo.BuildDate=${BUILD_DATE}"

echo "${LD_FLAGS}"
go build -ldflags "${LD_FLAGS}" -o ../bin/playground ../cmd/playground

echo "==> Starting..."
exec ../bin/playground
