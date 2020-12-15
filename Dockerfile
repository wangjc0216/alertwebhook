FROM alpine:3.10 AS final

ENV APP_PATH="/app/log_webhook"
WORKDIR "/app"

# 拷贝程序，如有必要另外拷贝其他文件
COPY  log_webhook ${APP_PATH}

# 运行程序
ENTRYPOINT ["/app/log_webhook"]
