FROM golang:latest AS builder

WORKDIR /src

COPY . /src

# 设置 Go module 代理
ENV GOPROXY=https://goproxy.cn,direct

# 下载并安装依赖
RUN go mod download && go mod verify


# 编译二进制文件
RUN go build -o /app/vmq

# 编译前端
FROM node:latest AS frontend

WORKDIR /src/frontend

COPY ./frontend /src/frontend

RUN npm install && npm run build


# 阶段二
FROM ubuntu:latest

# 从 builder 阶段拷贝二进制文件
COPY --from=builder /app/vmq /app/vmq

# 从 frontend 阶段拷贝前端文件
COPY --from=frontend /src/frontend/dist /app/web
