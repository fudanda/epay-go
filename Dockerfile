# epay-go/Dockerfile
# 前端构建阶段
FROM node:22-alpine AS web-builder

WORKDIR /app/web

COPY web/package.json web/package-lock.json ./
RUN npm ci

COPY web/ ./
RUN npm run build

# 后端构建阶段
FROM golang:1.25-alpine AS builder

WORKDIR /app

ARG GOPROXY=https://goproxy.cn,direct
ENV GOPROXY=${GOPROXY}

# 安装依赖
RUN apk add --no-cache git

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 注入前端构建产物，供 go:embed 打包
COPY --from=web-builder /app/web/dist ./web/dist

# 构建
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o epay-server ./cmd/server

# 运行阶段
FROM alpine:3.19

WORKDIR /app

# 安装必要工具
RUN apk --no-cache add ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 从构建阶段复制二进制文件
COPY --from=builder /app/epay-server .
COPY --from=builder /app/config.example.yaml ./config.yaml

# 暴露端口
EXPOSE 8080

# 运行
CMD ["./epay-server"]
