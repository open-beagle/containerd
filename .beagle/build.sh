#!/usr/bin/env bash

set -ex

git config --global --add safe.directory $PWD

apt-get update
apt-get install -y --no-install-recommends libbtrfs-dev

export GO111MODULE=off
export STATIC=true
export BUILDTAGS="seccomp"

export GOARCH=amd64
make
mkdir -p _output/linux-$GOARCH/
mv bin/* _output/linux-$GOARCH/

export GOARCH=arm64
export CC=aarch64-linux-gnu-gcc
make
mkdir -p _output/linux-$GOARCH/
mv bin/* _output/linux-$GOARCH/
