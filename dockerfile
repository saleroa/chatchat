# 设置基础镜像
FROM golang:1.19-alpine AS builder

ENV GOOS linux
ENV GOPROXY https://goproxy.cn,direct

# 创建工作目录
WORKDIR /app
ENV TZ Asia/Shanghai
# 将 YAML 配置文件添加到工作目录中

# 将应用程序代码添加到工作目录中
#COPY config.yaml .
#COPY sentinel.yaml .
COPY . .

# 构建应用程序
RUN go build -o main ./cmd/main.go

# 第二阶段：生成最终镜像
FROM alpine:latest

# 安装所需的运行时依赖（如果有的话）
RUN apk --no-cache add ca-certificates
RUN apk update && apk add tzdata
ENV TZ Asia/Shanghai

WORKDIR /app

# 从第一阶段中复制编译好的应用程序
COPY --from=builder /app/cmd/main .
COPY --from=builder /app/config.yaml .
COPY --from=builder /app/sentinel.yaml .


# 设置容器的入口命令
EXPOSE 8088
CMD ["./main"]