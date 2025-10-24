# Stage 1: Build Go binary
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o runner main.go

# Stage 2: Runtime (C++ only)
FROM alpine:3.20

RUN apk add --no-cache \
    g++ \
    bash \
    && rm -rf /var/cache/apk/*

ENV TZ=Asia/Ho_Chi_Minh
ENV CALLBACK_API_URL=http://localhost:8080/submission/handleCallBackApiHeheheheasd
ENV REDIS_ADDR=localhost:6379

WORKDIR /app
COPY --from=builder /app/runner .

CMD ["./runner"]
