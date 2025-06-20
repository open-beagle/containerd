#!/usr/bin/env bash

set -ex

git config --global --add safe.directory $PWD

export CGO_ENABLED=0
export GO111MODULE=off
export STATIC=true
export BUILDTAGS="seccomp"
export VERSION=${VERSION:-2.0.0}
export REVISION=$(git rev-parse --short HEAD)

export TARGETPLATFORM=linux/amd64
xx-apk add musl-dev gcc btrfs-progs-dev
xx-go --wrap
make
mkdir -p _output/$TARGETPLATFORM/
mv bin/* _output/$TARGETPLATFORM/

export TARGETPLATFORM=linux/arm64
xx-apk add musl-dev gcc btrfs-progs-dev
xx-go --wrap
make
mkdir -p _output/$TARGETPLATFORM/
mv bin/* _output/$TARGETPLATFORM/
