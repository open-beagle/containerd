# containerd

<https://github.com/containerd/containerd>

```bash
git -C ansible-docker-containerd remote add upstream git@github.com:containerd/containerd.git

git -C ansible-docker-containerd fetch upstream

git -C ansible-docker-containerd merge v2.1.9
```

## git

```bash
# ./.github/workflows/build-2.1.yml
git -C ansible-docker-containerd checkout release-v2.1 && \
git -C ansible-docker-containerd merge main && \
git -C ansible-docker-containerd push origin release-v2.1 && \
git -C ansible-docker-containerd checkout main
```

## debug

```bash
# loong64
docker run -it --rm \
-v $PWD/:/go/src/github.com/containerd/containerd/ \
-v $PWD/ansible-docker-containerd:/go/src/github.com/containerd/containerd/v2 \
-w /go/src/github.com/containerd/containerd/v2 \
-e VERSION=2.1.9-beagle \
registry.cn-qingdao.aliyuncs.com/wod/golang:1.24-loongnix \
bash .beagle/build-loong64.sh

# amd64&arm64
docker run -it --rm \
-v $PWD/:/go/src/github.com/containerd/containerd/ \
-v $PWD/ansible-docker-containerd:/go/src/github.com/containerd/containerd/v2 \
-w /go/src/github.com/containerd/containerd/v2 \
-e VERSION=2.1.9-beagle \
registry.cn-qingdao.aliyuncs.com/wod/golang:1.24-alpine \
bash .beagle/build-cross.sh
```

## test

```bash
file ansible-docker-containerd/_output/linux/amd64/containerd

# amd64
docker run -it --rm \
-v $PWD/:/go/src/github.com/containerd/containerd/ \
-v $PWD/ansible-docker-containerd:/go/src/github.com/containerd/containerd/v2 \
-w /go/src/github.com/containerd/containerd/v2 \
registry.cn-qingdao.aliyuncs.com/wod/debian:bullseye-amd64 \
./_output/linux/amd64/containerd -v

docker run -it --rm \
-v $PWD/:/go/src/github.com/containerd/containerd/ \
-v $PWD/ansible-docker-containerd:/go/src/github.com/containerd/containerd/v2 \
-w /go/src/github.com/containerd/containerd/v2 \
registry.cn-qingdao.aliyuncs.com/wod/alpine:3-amd64 \
./_output/linux/amd64/containerd -v

# arm64
docker run -it --rm \
-v $PWD/:/go/src/github.com/containerd/containerd/ \
-v $PWD/ansible-docker-containerd:/go/src/github.com/containerd/containerd/v2 \
-w /go/src/github.com/containerd/containerd/v2 \
registry.cn-qingdao.aliyuncs.com/wod/debian:bullseye-arm64 \
./_output/linux/arm64/containerd -v

docker run -it --rm \
-v $PWD/:/go/src/github.com/containerd/containerd/ \
-v $PWD/ansible-docker-containerd:/go/src/github.com/containerd/containerd/v2 \
-w /go/src/github.com/containerd/containerd/v2 \
registry.cn-qingdao.aliyuncs.com/wod/alpine:3-arm64 \
./_output/linux/arm64/containerd -v
```
