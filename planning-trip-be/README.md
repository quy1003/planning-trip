 # Planning Trip Backend

## Mục tiêu
Backend Go này phục vụ API cho ứng dụng Planning Trip trong monorepo. Nó dùng router Gin để xử lý HTTP, và cấu trúc chia rõ nhiệm vụ giữa entrypoint, business logic, handler và các tiện ích dùng chung.

## Cấu trúc thư mục
- `cmd/server/main.go`: entrypoint duy nhất, chịu trách nhiệm khởi tạo dịch vụ, middleware, router, và cổng lắng nghe. Giữ phần này mỏng để dễ mở rộng route mới mà không cần chỉnh lại service/transport.
- `internal/database`: quản lý kết nối database bằng GORM (đọc `DATABASE_URL`, mở kết nối, ping, cấu hình pool).
- `internal/repository`: tầng thao tác DB (GORM query), tách riêng khỏi service.
- `internal/service`: mỗi thư mục con là một domain/service (hiện có `health`). Service chỉ tập trung vào logic nghiệp vụ và trả DTO, không xử lý HTTP.
- `internal/transport/http`: chứa các handler chuyên biệt cho HTTP. Mỗi handler chỉ lo việc đọc request, gọi service, và trả response chuẩn (qua helper `internal/response`).
- `internal/middleware`: các middleware dùng chung như logging, request ID, recovery, hoặc auth sẽ nằm ở đây và được thêm vào Gin router với `router.Use(...)`.
- `internal/response`: helper để chuẩn hóa payload trả cho client: `{success,data,message,timestamp}` và hỗ trợ APIError để gói lỗi/capture exception.

## Các file chính hiện có
- `go.mod`/`go.sum`: mô-đun Go, quản lý dependency.
- `cmd/server/main.go`: load `.env`, mở kết nối DB bằng GORM, đăng ký route qua Gin.
- `cmd/migrate/main.go`: entrypoint chạy migration GORM để tạo/cập nhật schema database.
- `internal/config/env.go`: load và validate env tập trung (fail-fast nếu thiếu `DATABASE_URL`).
- `internal/database/postgres.go`: helper mở kết nối PostgreSQL với GORM.
- `internal/database/migrate.go`: chạy `AutoMigrate` cho toàn bộ model.
- `internal/model/models.go`: định nghĩa model GORM cho các bảng theo thiết kế README tổng.
- `internal/repository/user/repository.go`: CRUD user với GORM.
- `internal/service/user/service.go`: business logic user (validate input + gọi repository).
- `internal/transport/http/user/handler.go`: route/controller CRUD user.
- `internal/service/health/service.go`: trả trạng thái `ok` cùng timestamp.
- `internal/transport/http/health/handler.go`: gọi service và sử dụng `response.Write`/`WriteError` để gửi JSON.
- `internal/middleware/logging.go`: log mỗi request và thời gian xử lý.
- `internal/response/response.go`: định nghĩa `APIResponse` và `APIError`, kèm hàm `Write`/`WriteError` dùng chung.

## Làm việc tiếp theo
1. Để thêm route mới: tạo service trong `internal/service/<domain>`, handler tương ứng `internal/transport/http/<domain>`, rồi đăng ký trong `cmd/server/main.go`.
2. Nếu cần middleware mới (auth, metrics…), đặt file trong `internal/middleware` rồi wrap router.
3. Dùng `response` helper cho mọi handler để response đồng nhất.

## Chạy server
Đảm bảo có biến môi trường:
```
DATABASE_URL=postgresql://postgres:postgres@localhost:5432/planning_trip
```

```
go run ./cmd/server
```
Mở `http://localhost:8080/health` để kiểm tra.

## API user (CRUD)
- `GET /users`: danh sách user.
- `GET /users/:id`: lấy chi tiết user.
- `POST /users`: tạo user.
- `PUT /users/:id`: cập nhật user.
- `DELETE /users/:id`: xóa user.

## Migrate database (GORM)
Chạy migration để tạo bảng từ model:
```
go run ./cmd/migrate
```
Sau khi migrate xong mới chạy API:
```
go run ./cmd/server
```

## Seed data giả lập
Sinh dữ liệu mẫu cho local database (users/trips/places/comments/reactions/photos...):
```
go run ./cmd/seed
```
Seed script sẽ tự chạy migrate trước khi insert dữ liệu.

## Live reload (air)
1. Cài đặt `air` (chỉ cần chạy 1 lần):
   ```
   go install github.com/air-verse/air@latest
   ```
   (hoặc dùng `brew install air` nếu trên macOS)
2. Khởi động backend với `air -c .air.toml`. Air sẽ theo dõi `cmd/` và `internal/` để tự động rebuild/restart mỗi khi bạn lưu file, cho cảm giác giống Nodemon.
