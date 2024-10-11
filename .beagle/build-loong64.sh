#!/usr/bin/env bash

set -ex

git config --global --add safe.directory $PWD

apt-get update
apt-get install -y --no-install-recommends libbtrfs-dev

export CGO_ENABLED=0
export GO111MODULE=off
export STATIC=true
export BUILDTAGS="seccomp"
export VERSION=${VERSION:-v2.0.0}
export REVISION=$(git rev-parse --short HEAD)

export GOARCH=loong64
export CC=loongarch64-linux-gnu-gcc
make
mkdir -p _output/linux/$GOARCH/
mv bin/* _output/linux/$GOARCH/
