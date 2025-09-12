# Stage 1: Build
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy và cache go.mod, go.sum trước để tận dụng layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy mã nguồn
COPY . .

# Build binary với tối ưu hóa
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o runner main.go

# Stage 2: Runtime
FROM alpine:3.20

# Cài đặt các runtime cần thiết
RUN apk add --no-cache \
    python3 \
    openjdk17 \
    g++ \
    bash \
    && rm -rf /var/cache/apk/*

# Thiết lập biến môi trường
ENV TZ=Asia/Ho_Chi_Minh
ENV CALLBACK_API_URL=http://localhost:8080/submission/handleCallBackApiHeheheheasd
ENV REDIS_ADDR=localhost:6379

WORKDIR /app
COPY --from=builder /app/runner .

CMD ["./runner"]