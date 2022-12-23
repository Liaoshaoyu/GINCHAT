FROM golang:1.17

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    GOPROXY=goproxy.io \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

ADD ./GINCHAT /code

WORKDIR /code

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

RUN echo 'Asia/Shanghai' >/etc/timezone

RUN go mod download

RUN go build -o /code/build/myapp .

EXPOSE 8081

ENTRYPOINT [ "/code/build/myapp" ]
