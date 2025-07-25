FROM golang:1.21-alpine AS builder

WORKDIR /app

# 安装依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# 运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

# 从构建阶段复制二进制文件和模板
COPY --from=builder /app/main .
COPY --from=builder /app/templates ./templates

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./main"]
