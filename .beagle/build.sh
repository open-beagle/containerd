#!/usr/bin/env bash

set -ex

apt-get update
apt-get install -y --no-install-recommends libbtrfs-dev

export GO111MODULE=off
export STATIC=true

export GOARCH=amd64
make
mkdir -p _output/linux-$GOARCH/
mv bin/* _output/linux-$GOARCH/

export GOARCH=arm64
export CC=aarch64-linux-gnu-gcc
make
mkdir -p _output/linux-$GOARCH/
mv bin/* _output/linux-$GOARCH/

export GOARCH=ppc64le
export CC=powerpc64le-linux-gnu-gcc
make
mkdir -p _output/linux-$GOARCH/
mv bin/* _output/linux-$GOARCH/

export GOARCH=mips64le
export CC=mips64el-linux-gnuabi64-gcc
make
mkdir -p _output/linux-$GOARCH/
mv bin/* _output/linux-$GOARCH/

# export GOARCH=loong64
# export CC=loongarch64-linux-gnu-gcc
# make
# mkdir -p _output/linux-$GOARCH/
# mv bin/* _output/linux-$GOARCH/