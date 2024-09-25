#!/usr/bin/env bash

set -ex

git config --global --add safe.directory $PWD

apt-get update
apt-get install -y --no-install-recommends libbtrfs-dev

export GO111MODULE=off
export STATIC=true
export BUILDTAGS="seccomp"

git apply .beagle/v1.7-add-loong64-support.patch
export GOARCH=loong64
export CC=loongarch64-linux-gnu-gcc
make
mkdir -p _output/linux-$GOARCH/
mv bin/* _output/linux-$GOARCH/
git apply -R .beagle/v1.7-add-loong64-support.patch
