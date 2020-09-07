# Copyright 2018 The KubeSphere Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

# Copyright 2018 The KubeSphere Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.
FROM golang:1.12 as ks-apiserver-builder

COPY / /go/src/kubesphere.io/kubesphere

WORKDIR /go/src/kubesphere.io/kubesphere
RUN CGO_ENABLED=0 GO111MODULE=on GOOS=linux GOARCH=amd64 GOFLAGS=-mod=vendor go build -i -ldflags '-w -s' -o ks-apiserver cmd/ks-apiserver/apiserver.go

FROM alpine:3.9
# 替换 Alpine 默认的源地址
# 防止出现 ERROR: http://dl-cdn.alpinelinux.org/alpine/v3.9/main: BAD signature
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN apk add --update ca-certificates && update-ca-certificates
COPY --from=ks-apiserver-builder /go/src/kubesphere.io/kubesphere/ks-apiserver /usr/local/bin/
CMD ["sh"]
