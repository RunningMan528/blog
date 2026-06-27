# 构建阶段
FROM golang:1.26.4-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o blog ./cmd/main.go

# 运行阶段
FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata && mkdir -p /app/uploads

COPY --from=builder /app/blog /app/blog

EXPOSE 8080

CMD [ "/app/blog" ]