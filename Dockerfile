# 阶段 1: 构建 Vue 前端
FROM node:20-alpine AS frontend-builder
WORKDIR /app
COPY zharchiver-frontend/package*.json ./
RUN npm install
COPY zharchiver-frontend/ .
RUN npm run build

# 阶段 2: 构建 Go 后端
# 使用 Debian 基础镜像 (bookworm) 以兼容最终的 Rod 镜像 (基于 Ubuntu/Debian) 和 CGO (sqlite3)
FROM golang:1.26-bookworm AS backend-builder
WORKDIR /app
COPY zharchiver/go.mod zharchiver/go.sum ./
RUN go mod download
COPY zharchiver/ .
# 必须开启 CGO 以编译 github.com/mattn/go-sqlite3
ENV CGO_ENABLED=1
RUN go build -o zharchiver .

# 阶段 3: 最终运行镜像 (自带 Chromium 和中文字体)
FROM ghcr.io/go-rod/rod:latest
WORKDIR /app

# 拷贝后端可执行文件
COPY --from=backend-builder /app/zharchiver .
# 拷贝前端打包后的纯静态网页文件，放到 dist 目录下供 Go 代理
COPY --from=frontend-builder /app/dist ./dist

# 创建持久化存储目录
RUN mkdir -p /app/storage

# 暴露后端端口
EXPOSE 8080

# 启动容器时执行后端程序
CMD ["./zharchiver"]
