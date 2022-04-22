FROM busybox

# 工作目录
WORKDIR /app/

# 将代码复制到容器中
COPY . .

# 运行程序（先提前编译一下 make build）
ENTRYPOINT web_app
