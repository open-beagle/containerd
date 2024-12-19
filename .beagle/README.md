# containerd

<https://github.com/containerd/containerd>

```bash
git remote add upstream git@github.com:containerd/containerd.git

git fetch upstream

git merge v2.0.1
```

## build

```bash
# golang build cross
docker run -it --rm \
  -v $PWD/:/go/src/github.com/containerd/containerd/v2 \
  -w /go/src/github.com/containerd/containerd/v2 \
  -e VERSION=2.0.1-beagle \
  registry.cn-qingdao.aliyuncs.com/wod/golang:1.22-alpine \
  bash .beagle/build-cross.sh

# golang build
docker run -it --rm \
  -v $PWD/:/go/src/github.com/containerd/containerd/v2 \
  -w /go/src/github.com/containerd/containerd/v2 \
  -e VERSION=2.0.1-beagle \
  registry.cn-qingdao.aliyuncs.com/wod/golang:1.22 \
  bash .beagle/build.sh

# golang build loong64
docker run -it --rm \
  -v $PWD/:/go/src/github.com/containerd/containerd/v2 \
  -w /go/src/github.com/containerd/containerd/v2 \
  -e VERSION=2.0.1-beagle \
  registry.cn-qingdao.aliyuncs.com/wod/golang:1.22-loongnix \
  bash .beagle/build-loong64.sh
```

## test

```bash
file _output/linux-amd64/containerd

# amd64-test
docker run -it --rm \
-v $PWD/:/go/src/github.com/containerd/containerd/v2 \
-w /go/src/github.com/containerd/containerd/v2 \
registry.cn-qingdao.aliyuncs.com/wod/debian:bullseye-amd64 \
./_output/linux-amd64/containerd -v

docker run -it --rm \
-v $PWD/:/go/src/github.com/containerd/containerd/v2 \
-w /go/src/github.com/containerd/containerd/v2 \
registry.cn-qingdao.aliyuncs.com/wod/alpine:3-amd64 \
./_output/linux-amd64/containerd -v

# arm64-test
docker run -it --rm \
-v $PWD/:/go/src/github.com/containerd/containerd/v2 \
-w /go/src/github.com/containerd/containerd/v2 \
registry.cn-qingdao.aliyuncs.com/wod/debian:bullseye-arm64 \
./_output/linux-arm64/containerd -v

docker run -it --rm \
-v $PWD/:/go/src/github.com/containerd/containerd/v2 \
-w /go/src/github.com/containerd/containerd/v2 \
registry.cn-qingdao.aliyuncs.com/wod/alpine:3-arm64 \
./_output/linux-arm64/containerd -v

# loong64-test
docker run -it --rm \
-v $PWD/:/go/src/github.com/containerd/containerd/v2 \
-w /go/src/github.com/containerd/containerd/v2 \
registry.cn-qingdao.aliyuncs.com/wod/alpine:3-loong64 \
./_output/linux-loong64/containerd -v
```

## cache

```bash
# 构建缓存-->推送缓存至服务器
docker run --rm \
  -e PLUGIN_REBUILD=true \
  -e PLUGIN_ENDPOINT=$S3_ENDPOINT_ALIYUN \
  -e PLUGIN_ACCESS_KEY=$S3_ACCESS_KEY_ALIYUN \
  -e PLUGIN_SECRET_KEY=$S3_SECRET_KEY_ALIYUN \
  -e DRONE_REPO_OWNER="open-beagle" \
  -e DRONE_REPO_NAME="containerd" \
  -e PLUGIN_MOUNT="./.git" \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  registry.cn-qingdao.aliyuncs.com/wod/devops-s3-cache:1.0

# 读取缓存-->将缓存从服务器拉取到本地
docker run --rm \
  -e PLUGIN_RESTORE=true \
  -e PLUGIN_ENDPOINT=$S3_ENDPOINT_ALIYUN \
  -e PLUGIN_ACCESS_KEY=$S3_ACCESS_KEY_ALIYUN \
  -e PLUGIN_SECRET_KEY=$S3_SECRET_KEY_ALIYUN \
  -e DRONE_REPO_OWNER="open-beagle" \
  -e DRONE_REPO_NAME="containerd" \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  registry.cn-qingdao.aliyuncs.com/wod/devops-s3-cache:1.0
```
