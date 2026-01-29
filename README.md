# Go-Code-Runner

Go-Code-Runner là một project **online code runner** đơn giản, mục đích chính là **chạy code** do người dùng gửi lên.  
Dự án **không tập trung vào chấm điểm hay quản lý bài thi**, mà đóng vai trò backend để compile và execute code trong môi trường tách biệt.

---

## Mục tiêu dự án

- Cung cấp service để **chạy source code theo yêu cầu**
- Làm backend cho các hệ thống luyện code / online IDE
- Đơn giản, dễ đọc, dễ mở rộng

---

## Công nghệ sử dụng

- Golang
- Redis
- Docker (tuỳ chọn, để sandbox code)
- Docker Compose (tuỳ chọn)

---

## Yêu cầu môi trường

- Go >= 1.21
- Redis **(bắt buộc)**
- Docker (tuỳ chọn)

---

## Bắt buộc: khởi chạy Redis trước

Ứng dụng **không tự khởi động Redis**.  
Người dùng cần đảm bảo Redis đang chạy **trước khi chạy server**.

### Chạy Redis bằng Docker (khuyến nghị)

```bash
docker run -d \
  --name redis \
  -p 6379:6379 \
  redis:7
```

### Hoặc chạy Redis cài sẵn trên máy

```bash
redis-server
```

Kiểm tra Redis hoạt động:

```bash
redis-cli ping
# PONG
```

---

## Cấu trúc thư mục (tham khảo)

```
.
├── main.go            # Entry point
├── internal/          # Business logic
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
└── README.md
```

---

## Cách chạy dự án

### 1. Đảm bảo Redis đang chạy

Redis phải chạy ở `localhost:6379` (hoặc đúng cấu hình trong code).

---

### 2. Chạy trực tiếp bằng Go

```bash
go run main.go
```

Hoặc build trước rồi chạy:

```bash
go build -o runner
./runner
```

---


## Nguyên lý hoạt động

1. Client gửi request chứa source code và thông tin ngôn ngữ
2. Server lưu code vào file tạm
3. Redis được dùng để lưu job / trạng thái / kết quả chạy
4. Server compile và execute code (có thể thông qua Docker)
5. Trả về output, error và exit status

---

## Đóng góp

1. Fork repository
2. Tạo branch mới
3. Commit thay đổi
4. Tạo Pull Request

