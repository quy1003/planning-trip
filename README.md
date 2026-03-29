# Planning Trip Monorepo

## 1) Mô tả dự án
Planning Trip là ứng dụng web hỗ trợ nhóm lập kế hoạch chuyến đi theo lộ trình. Mỗi chuyến đi được thể hiện như một bài đăng cho phép:
- Tạo danh sách địa điểm và gợi ý tuyến đường.
- Hiển thị các địa điểm lên bản đồ tương tác.
- Lên timetable/schedule chi tiết theo từng ngày và giờ.
- Thành viên khác bình luận, thả cảm xúc cho từng điểm hoặc toàn bộ chuyến đi.
- Đóng góp ảnh vào album chung kèm chú thích và cảm nhận cá nhân.

Ví dụ: Quý Nguyễn tạo plan Nha Trang, thêm Tháp Bà, Vinpearl Land, xếp lịch theo thứ tự ngày. Các thành viên khác thấy bản đồ, timetable, bình luận, thả cảm xúc và tải ảnh lên album chung.

## 2) Mục tiêu sản phẩm
- Giúp nhóm lên kế hoạch chuyến đi trực quan, có bản đồ và lịch trình rõ ràng.
- Tạo không gian tương tác giữa thành viên bằng bình luận, reaction, góp ý.
- Lưu giữ kỷ niệm chuyến đi qua album ảnh đóng góp chung.

## 3) Phạm vi MVP
- Quản lý người dùng, thành viên của một trip.
- Tạo/chỉnh sửa/xem thông tin trip.
- Thêm địa điểm và hiển thị trên map.
- Quản lý timeline/schedule cho từng ngày đi.
- Bình luận và reaction cho trip, địa điểm, ảnh.
- Album ảnh chung, người dùng có thể upload và bày tỏ cảm nghĩ.

## 4) Các vai trò và quyền hạn
- **Owner**: tạo trip, mời/thu hồi thành viên, chỉnh sửa toàn bộ nội dung.
- **Editor**: chỉnh sửa địa điểm, lịch trình, nội dung trip.
- **Viewer**: xem trip, bình luận, reaction, đóng góp ảnh nếu được cấp quyền.

## 5) Danh sách thực thể nghiệp vụ
- `User`
- `Trip`
- `TripMember`
- `Place`
- `TripPlace`
- `TripScheduleItem`
- `Comment`
- `Reaction`
- `Album`
- `Photo`
- `PhotoTag` (tuỳ chọn)
- `Notification` (phase 2)

## 6) Đề xuất bảng database (PostgreSQL)

### `users`
- `id` (PK, uuid)
- `full_name`
- `email` (unique)
- `password_hash`
- `avatar_url`
- `bio`
- `created_at`, `updated_at`

### `trips`
- `id` (PK, uuid)
- `creator_id` (FK -> users.id)
- `title`
- `description`
- `start_date`, `end_date`
- `visibility` (`private`, `group`, `public`)
- `status` (`draft`, `published`, `archived`)
- `cover_image_url`
- `created_at`, `updated_at`

### `trip_members`
- `id` (PK, uuid)
- `trip_id` (FK -> trips.id)
- `user_id` (FK -> users.id)
- `role` (`owner`, `editor`, `viewer`)
- `joined_at`
- Unique (`trip_id`, `user_id`)

### `places`
- `id` (PK, uuid)
- `name`
- `address`
- `lat`, `lng`
- `google_place_id` (nullable)
- `created_by` (FK -> users.id, nullable)
- `created_at`, `updated_at`

### `trip_places`
- `id` (PK, uuid)
- `trip_id` (FK -> trips.id)
- `place_id` (FK -> places.id)
- `title` (nullable, dùng để ghi tên tuỳ chỉnh)
- `note`
- `day_index`
- `order_index`
- `estimated_cost` (nullable)
- `created_by` (FK -> users.id)
- `created_at`, `updated_at`

### `trip_schedule_items`
- `id` (PK, uuid)
- `trip_id` (FK -> trips.id)
- `trip_place_id` (FK -> trip_places.id, nullable)
- `title`
- `description`
- `start_time`, `end_time`
- `day_index`
- `order_index`
- `created_by` (FK -> users.id)
- `created_at`, `updated_at`

### `albums`
- `id` (PK, uuid)
- `trip_id` (FK -> trips.id)
- `title`
- `description`
- `created_by` (FK -> users.id)
- `created_at`, `updated_at`

### `photos`
- `id` (PK, uuid)
- `album_id` (FK -> albums.id)
- `trip_id` (FK -> trips.id)
- `uploader_id` (FK -> users.id)
- `url`
- `thumbnail_url`
- `caption`
- `taken_at` (nullable)
- `trip_place_id` (FK -> trip_places.id, nullable)
- `created_at`, `updated_at`

### `photo_tags` (tuỳ chọn)
- `id` (PK, uuid)
- `photo_id` (FK -> photos.id)
- `user_id` (FK -> users.id)
- `tagged_by` (FK -> users.id)
- `created_at`
- Unique (`photo_id`, `user_id`)

### `comments` (đối tượng đa hình)
- `id` (PK, uuid)
- `author_id` (FK -> users.id)
- `target_type` (`trip`, `trip_place`, `photo`)
- `target_id` (uuid)
- `parent_id` (FK -> comments.id, nullable)
- `content`
- `created_at`, `updated_at`, `deleted_at`

### `reactions` (đối tượng đa hình)
- `id` (PK, uuid)
- `user_id` (FK -> users.id)
- `target_type` (`trip`, `trip_place`, `comment`, `photo`)
- `target_id` (uuid)
- `reaction_type` (`like`, `love`, `wow`, `sad`, `angry`)
- `created_at`
- Unique (`user_id`, `target_type`, `target_id`, `reaction_type`)

### `notifications` (phase 2)
- `id` (PK, uuid)
- `user_id` (FK -> users.id)
- `type`
- `data` (jsonb)
- `is_read`
- `created_at`

## 7) Quan hệ chính
- Một `user` có thể tạo nhiều `trips`.
- Một `trip` liên kết với các `trip_members`, `trip_places`, `schedule_items`, `comments`, `reactions`, `albums`.
- Một `album` chứa nhiều `photos`.
- `comments` và `reactions` dùng cặp `target_type` + `target_id` để hỗ trợ nhiều kiểu đối tượng.
- `places` chuẩn hóa dữ liệu vị trí để tái sử dụng giữa các trip.

## 8) Định hướng monorepo
Đề xuất cấu trúc:
- `apps/web`: Frontend (Next.js)
- `apps/api`: Backend (NestJS/Fastify)
- `packages/db`: Prisma schema + migration
- `packages/types`: DTO, type chia sẻ
- `packages/ui`: Component UI dùng chung
- `packages/config`: eslint/tsconfig/prettier chung

## 9) Lộ trình tiếp theo
1. Hoàn thiện ERD và schema database MVP.
2. Scaffold monorepo (`pnpm` + `turborepo`).
3. Cài đặt auth, user, trip CRUD.
4. Tích hợp map, trip places, schedule.
5. Bổ sung comment/reaction và album upload.
6. Hoàn thiện phân quyền, thông báo và tối ưu UX.
