# 启动编译环境
FROM golang:1.16-alpine AS builder

# 配置编译环境
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct

# 拷贝源代码到镜像中
COPY . /go/src/api/goods_web

# 编译
WORKDIR /go/src/api/goods_web
RUN go install ./goods_web/...

FROM alpine:3.13
COPY --from=builder /go/bin/goods_web /bin/goods_web
ENV ADDR=:8080

# 申明暴露的端口
EXPOSE 8080

# 设置服务入口
ENTRYPOINT [ "/bin/goods_web" ]
