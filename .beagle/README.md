# git

```bash
git remote add upstream git@github.com:containerd/containerd.git

git fetch upstream

git merge v1.6.9
```

## build

```bash
docker run -it --rm \
-v $PWD/:/go/src/github.com/containerd/containerd \
-w /go/src/github.com/containerd/containerd \
-e VERSION=v1.6.9-beagle \
registry.cn-qingdao.aliyuncs.com/wod/golang:1.19 \
bash .beagle/build.sh
```

## test

```bash
file _output/linux-amd64/containerd

# amd64-test
docker run -it --rm \
-v $PWD/:/go/src/github.com/containerd/containerd \
-w /go/src/github.com/containerd/containerd \
registry.cn-qingdao.aliyuncs.com/wod/debian:bullseye-amd64 \
./_output/linux-amd64/containerd -v

docker run -it --rm \
-v $PWD/:/go/src/github.com/containerd/containerd \
-w /go/src/github.com/containerd/containerd \
registry.cn-qingdao.aliyuncs.com/wod/alpine:3-amd64 \
./_output/linux-amd64/containerd -v

# arm64-test
docker run -it --rm \
-v $PWD/:/go/src/github.com/containerd/containerd \
-w /go/src/github.com/containerd/containerd \
registry.cn-qingdao.aliyuncs.com/wod/debian:bullseye-arm64 \
./_output/linux-arm64/containerd -v

docker run -it --rm \
-v $PWD/:/go/src/github.com/containerd/containerd \
-w /go/src/github.com/containerd/containerd \
registry.cn-qingdao.aliyuncs.com/wod/alpine:3-arm64 \
./_output/linux-arm64/containerd -v

# ppc64le-test
docker run -it --rm \
-v $PWD/:/go/src/github.com/containerd/containerd \
-w /go/src/github.com/containerd/containerd \
registry.cn-qingdao.aliyuncs.com/wod/debian:bullseye-ppc64le \
./_output/linux-ppc64le/containerd -v

docker run -it --rm \
-v $PWD/:/go/src/github.com/containerd/containerd \
-w /go/src/github.com/containerd/containerd \
registry.cn-qingdao.aliyuncs.com/wod/alpine:3-ppc64le \
./_output/linux-ppc64le/containerd -v

# mips64le-test
docker run -it --rm \
-v $PWD/:/go/src/github.com/containerd/containerd \
-w /go/src/github.com/containerd/containerd \
registry.cn-qingdao.aliyuncs.com/wod/debian:bullseye-mips64le \
./_output/linux-mips64le/containerd -v

docker run -it --rm \
-v $PWD/:/go/src/github.com/containerd/containerd \
-w /go/src/github.com/containerd/containerd \
registry.cn-qingdao.aliyuncs.com/wod/alpine:3-mips64le \
./_output/linux-mips64le/containerd -v
```

## cache

```bash
# ????????????-->????????????????????????
docker run --rm \
  -e PLUGIN_REBUILD=true \
  -e PLUGIN_ENDPOINT=$PLUGIN_ENDPOINT \
  -e PLUGIN_ACCESS_KEY=$PLUGIN_ACCESS_KEY \
  -e PLUGIN_SECRET_KEY=$PLUGIN_SECRET_KEY \
  -e PLUGIN_PATH="/cache/open-beagle/containerd" \
  -e PLUGIN_MOUNT="./.git" \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  registry.cn-qingdao.aliyuncs.com/wod/devops-s3-cache:1.0

# ????????????-->????????????????????????????????????
docker run --rm \
  -e PLUGIN_RESTORE=true \
  -e PLUGIN_ENDPOINT=$PLUGIN_ENDPOINT \
  -e PLUGIN_ACCESS_KEY=$PLUGIN_ACCESS_KEY \
  -e PLUGIN_SECRET_KEY=$PLUGIN_SECRET_KEY \
  -e PLUGIN_PATH="/cache/open-beagle/containerd" \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  registry.cn-qingdao.aliyuncs.com/wod/devops-s3-cache:1.0
```
