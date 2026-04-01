# Planning Trip

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

## 6) Diagram

```mermaid
erDiagram
    USERS {
        uuid id PK
        varchar full_name
        varchar email
        varchar password_hash
        varchar avatar_url
        varchar bio
        timestamp created_at
        timestamp updated_at
    }

    TRIPS {
        uuid id PK
        uuid creator_id FK
        varchar title
        text description
        date start_date
        date end_date
        varchar visibility
        varchar status
        varchar cover_image_url
        timestamp created_at
        timestamp updated_at
    }

    TRIP_MEMBERS {
        uuid id PK
        uuid trip_id FK
        uuid user_id FK
        varchar role
        timestamp joined_at
    }

    PLACES {
        uuid id PK
        varchar name
        varchar address
        decimal lat
        decimal lng
        varchar google_place_id
        uuid created_by FK
        timestamp created_at
        timestamp updated_at
    }

    TRIP_PLACES {
        uuid id PK
        uuid trip_id FK
        uuid place_id FK
        varchar title
        text note
        int day_index
        int order_index
        decimal estimated_cost
        uuid created_by FK
        timestamp created_at
        timestamp updated_at
    }

    TRIP_SCHEDULE_ITEMS {
        uuid id PK
        uuid trip_id FK
        uuid trip_place_id FK
        varchar title
        text description
        timestamptz start_time
        timestamptz end_time
        int day_index
        int order_index
        uuid created_by FK
        timestamp created_at
        timestamp updated_at
    }

    ALBUMS {
        uuid id PK
        uuid trip_id FK
        varchar title
        text description
        uuid created_by FK
        timestamp created_at
        timestamp updated_at
    }

    PHOTOS {
        uuid id PK
        uuid album_id FK
        uuid trip_id FK
        uuid uploader_id FK
        varchar url
        varchar thumbnail_url
        text caption
        timestamptz taken_at
        uuid trip_place_id FK
        timestamp created_at
        timestamp updated_at
    }

    PHOTO_TAGS {
        uuid id PK
        uuid photo_id FK
        uuid user_id FK
        uuid tagged_by FK
        timestamp created_at
    }

    COMMENTS {
        uuid id PK
        uuid author_id FK
        varchar target_type
        uuid target_id
        uuid parent_id FK
        text content
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at
    }

    REACTIONS {
        uuid id PK
        uuid user_id FK
        varchar target_type
        uuid target_id
        varchar reaction_type
        timestamp created_at
    }

    NOTIFICATIONS {
        uuid id PK
        uuid user_id FK
        varchar type
        json data
        boolean is_read
        timestamp created_at
    }

    USERS ||--o{ TRIPS : creates
    USERS ||--o{ TRIP_MEMBERS : participates
    USERS ||--o{ COMMENTS : authors
    USERS ||--o{ REACTIONS : authors
    USERS ||--o{ ALBUMS : owns
    USERS ||--o{ PHOTOS : uploads
    TRIPS ||--o{ TRIP_MEMBERS : contains
    TRIPS ||--o{ TRIP_PLACES : contains
    TRIPS ||--o{ TRIP_SCHEDULE_ITEMS : schedules
    TRIPS ||--o{ ALBUMS : hosts
    TRIPS ||--o{ COMMENTS : receives
    TRIPS ||--o{ REACTIONS : receives
    TRIPS ||--o{ PHOTOS : aggregates
    TRIP_PLACES ||--o{ TRIP_SCHEDULE_ITEMS : supports
    TRIP_PLACES ||--o{ PHOTOS : anchors
    PLACES ||--o{ TRIP_PLACES : referenced_by
    ALBUMS ||--o{ PHOTOS : contains
    PHOTOS ||--o{ PHOTO_TAGS : tagged_users
    COMMENTS ||--o{ REACTIONS : receives
    USERS ||--o{ NOTIFICATIONS : receives
    TRIPS ||--o{ NOTIFICATIONS : triggers
```
## 7) Quan hệ chính

- Một `user` có thể tạo nhiều `trips`.
- Một `trip` liên kết với các `trip_members`, `trip_places`, `schedule_items`, `comments`, `reactions`, `albums`.
- Một `album` chứa nhiều `photos`.
- `comments` và `reactions` dùng cặp `target_type` + `target_id` để hỗ trợ nhiều kiểu đối tượng.
- `places` chuẩn hóa dữ liệu vị trí để tái sử dụng giữa các trip.

---

Author: Selten03 aka Nguyễn Thi Quý  
Contact: quy021003@gmail.com  
Linkedin: https://www.linkedin.com/in/nguyenquythi/
