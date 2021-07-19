# 构造可执行文件
FROM --platform=$TARGETPLATFORM golang:1.16-alpine as build_bin
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.io
# 停用cgo
ENV CGO_ENABLED=0
# 安装grpc_health_probe
RUN go get github.com/grpc-ecosystem/grpc-health-probe
WORKDIR /code
COPY go.mod /code/go.mod
COPY go.sum /code/go.sum
# 添加源文件
COPY echo_pb /code/echo_pb/
COPY echo_sdk /code/echo_sdk/
COPY echo_serv /code/echo_serv/
COPY main.go /code/main.go
RUN go build -ldflags "-s -w" -o echoserv-go main.go

# 使用upx压缩可执行文件
FROM --platform=$TARGETPLATFORM alpine:3.11 as upx
WORKDIR /code
# 安装upx
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN apk update && apk add --no-cache upx && rm -rf /var/cache/apk/*
COPY --from=build_bin /code/echoserv-go .
RUN upx --best --lzma -o echoserv echoserv-go


# 使用压缩过的可执行文件构造镜像
FROM --platform=$TARGETPLATFORM alpine:3.12.2 as build_upximg
# 打包镜像
COPY --from=build_bin /go/bin/grpc-health-probe .
RUN chmod +x /grpc-health-probe
COPY --from=upx /code/echoserv .
RUN chmod +x /echoserv
EXPOSE 5000
# HEALTHCHECK --interval=30s --timeout=30s --start-period=30s --retries=3 CMD [ "/grpc-health-probe","-addr=:5000" ]
ENTRYPOINT [ "/echoserv"]