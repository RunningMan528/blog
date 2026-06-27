# 构建阶段
FROM golang:1.26.4-alpine AS builder

WORKDIR /app

ENV GOPROXY=https://goproxy.cn,direct

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o blog ./cmd/main.go

# 运行阶段
FROM alpine:3.22

WORKDIR /app

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && apk add --no-cache ca-certificates tzdata \
    && mkdir -p /app/uploads

COPY --from=builder /app/blog /app/blog

EXPOSE 8080

CMD [ "/app/blog" ]