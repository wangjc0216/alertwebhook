FROM golang:1.13 AS builder
WORKDIR /alertwebhook/
COPY main.go .
COPY vendor/github.com   /go/src/github.com
COPY vendor/gorm.io   /go/src/gorm.io

RUN CGO_ENABLED=0 GOOS=linux go build -o alertwebhook .

FROM alpine:3.10 AS final

ENV APP_PATH="/app/alertwebhook"
WORKDIR "/app"

# 拷贝程序，如有必要另外拷贝其他文件
COPY  --from=builder /alertwebhook/alertwebhook  ${APP_PATH}

# 运行程序
ENTRYPOINT ["/app/alertwebhook"]
