# 使用官方的Go镜像作为构建环境
FROM golang:1.19 as builder

# 安装git
RUN apt-get update && apt-get install -y git

# 设置工作目录
WORKDIR /app

# 复制当前项目的源代码
COPY ./main.go /app/main.go

# 克隆funcaptcha库的最新代码
RUN git clone https://github.com/acheong08/funcaptcha.git

# COPY ./funcaptcha /app/funcaptcha
COPY ./api.go /app/funcaptcha/api.go
COPY ./funcaptcha.go /app/funcaptcha/funcaptcha.go

# 初始化Go模块并添加funcaptcha依赖
RUN go mod init arkose && go mod edit -require=github.com/acheong08/funcaptcha@latest

# 使用replace指令引用本地funcaptcha库
RUN echo "replace github.com/acheong08/funcaptcha => ./funcaptcha" >> go.mod

# 清除模块缓存
RUN go clean -modcache
# 处理所有依赖
RUN go mod tidy

# 添加所需的依赖
RUN go get github.com/gin-gonic/gin

# 构建应用
RUN CGO_ENABLED=0 go build -o /arkose-token-service


# 使用alpine作为基础镜像来创建一个最终镜像
FROM alpine

# 从构建器镜像中复制构建好的应用
COPY --from=builder /arkose-token-service /arkose-token-service

# 声明服务运行在哪个端口
EXPOSE 8080

# 启动服务
ENTRYPOINT ["/arkose-token-service"]
