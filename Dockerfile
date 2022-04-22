FROM golang:alpine AS builder

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# 移动到工作目录：/build
WORKDIR /build

# 复制项目中的 go.mod 和 go.sum文件并下载依赖信息
COPY go.mod .
COPY go.sum .
RUN go mod tidy

# 将我们的代码编译成二进制可执行文件 Star-Travels
RUN go build -o web_app .

# 将代码复制到容器中
COPY . .

# 运行程序
ENTRYPOINT web_app
