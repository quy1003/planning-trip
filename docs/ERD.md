# ERD - Planning Trip

Đây là sơ đồ ERD mô tả các bảng cốt lõi của dự án Planning Trip cùng các mối quan hệ chính.

## Diagram
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

## Ghi chú
- `target_type` trong `comments`/`reactions` hỗ trợ đa hình, cho phép bình luận/reaction trên `trip`, `trip_place`, `photo`.
- `trip_places` tham chiếu `places` để tái sử dụng vị trí và giữ dữ liệu bản đồ chuẩn.
- `photo_tags` cho phép gắn người tham gia vào mỗi bức ảnh.
- `notifications` là bảng mở rộng cho các thông báo người dùng.
