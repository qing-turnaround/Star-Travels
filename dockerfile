FROM busybox

workdir /app

# 运行程序（先提前编译一下 make build）
ENTRYPOINT ["./web_app"]

