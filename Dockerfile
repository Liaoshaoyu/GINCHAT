FROM golang:1.17

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    GOPROXY=goproxy.io \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# 移动到工作目录：/build
WORKDIR /build

# 将代码复制到容器中
COPY ./* /build

RUN go env -w GOPROXY=https://goproxy.cn,direct

RUN go mod download

# 声明服务端口
EXPOSE 8081

# 启动容器时运行的命令
CMD ["go", "run", "main.go"]