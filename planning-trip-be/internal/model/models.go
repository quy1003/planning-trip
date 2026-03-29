package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type User struct {
	ID           string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	FullName     string    `gorm:"size:255;not null"`
	Email        string    `gorm:"size:255;not null;uniqueIndex"`
	PasswordHash string    `gorm:"size:255;not null"`
	AvatarURL    string    `gorm:"size:500"`
	Bio          string    `gorm:"type:text"`
	CreatedAt    time.Time `gorm:"not null;default:now()"`
	UpdatedAt    time.Time `gorm:"not null;default:now()"`
}

type Trip struct {
	ID            string     `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CreatorID     string     `gorm:"type:uuid;not null;index"`
	Title         string     `gorm:"size:255;not null"`
	Description   string     `gorm:"type:text"`
	StartDate     *time.Time `gorm:"type:date"`
	EndDate       *time.Time `gorm:"type:date"`
	Visibility    string     `gorm:"size:20;not null;default:private"`
	Status        string     `gorm:"size:20;not null;default:draft"`
	CoverImageURL string     `gorm:"size:500"`
	CreatedAt     time.Time  `gorm:"not null;default:now()"`
	UpdatedAt     time.Time  `gorm:"not null;default:now()"`
}

type TripMember struct {
	ID       string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TripID   string    `gorm:"type:uuid;not null;uniqueIndex:ux_trip_members_trip_user"`
	UserID   string    `gorm:"type:uuid;not null;uniqueIndex:ux_trip_members_trip_user"`
	Role     string    `gorm:"size:20;not null;default:viewer"`
	JoinedAt time.Time `gorm:"not null;default:now()"`
}

type Place struct {
	ID            string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name          string    `gorm:"size:255;not null"`
	Address       string    `gorm:"size:500"`
	Lat           float64   `gorm:"type:numeric(10,7);not null"`
	Lng           float64   `gorm:"type:numeric(10,7);not null"`
	GooglePlaceID string    `gorm:"size:255"`
	CreatedBy     *string   `gorm:"type:uuid;index"`
	CreatedAt     time.Time `gorm:"not null;default:now()"`
	UpdatedAt     time.Time `gorm:"not null;default:now()"`
}

type TripPlace struct {
	ID            string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TripID        string    `gorm:"type:uuid;not null;index"`
	PlaceID       string    `gorm:"type:uuid;not null;index"`
	Title         string    `gorm:"size:255"`
	Note          string    `gorm:"type:text"`
	DayIndex      int       `gorm:"not null;default:1"`
	OrderIndex    int       `gorm:"not null;default:0"`
	EstimatedCost float64   `gorm:"type:numeric(12,2)"`
	CreatedBy     string    `gorm:"type:uuid;not null;index"`
	CreatedAt     time.Time `gorm:"not null;default:now()"`
	UpdatedAt     time.Time `gorm:"not null;default:now()"`
}

type TripScheduleItem struct {
	ID          string     `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TripID      string     `gorm:"type:uuid;not null;index"`
	TripPlaceID *string    `gorm:"type:uuid;index"`
	Title       string     `gorm:"size:255;not null"`
	Description string     `gorm:"type:text"`
	StartTime   *time.Time `gorm:"type:timestamptz"`
	EndTime     *time.Time `gorm:"type:timestamptz"`
	DayIndex    int        `gorm:"not null;default:1"`
	OrderIndex  int        `gorm:"not null;default:0"`
	CreatedBy   string     `gorm:"type:uuid;not null;index"`
	CreatedAt   time.Time  `gorm:"not null;default:now()"`
	UpdatedAt   time.Time  `gorm:"not null;default:now()"`
}

type Album struct {
	ID          string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TripID      string    `gorm:"type:uuid;not null;index"`
	Title       string    `gorm:"size:255;not null"`
	Description string    `gorm:"type:text"`
	CreatedBy   string    `gorm:"type:uuid;not null;index"`
	CreatedAt   time.Time `gorm:"not null;default:now()"`
	UpdatedAt   time.Time `gorm:"not null;default:now()"`
}

type Photo struct {
	ID           string     `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	AlbumID      string     `gorm:"type:uuid;not null;index"`
	TripID       string     `gorm:"type:uuid;not null;index"`
	UploaderID   string     `gorm:"type:uuid;not null;index"`
	URL          string     `gorm:"size:1000;not null"`
	ThumbnailURL string     `gorm:"size:1000"`
	Caption      string     `gorm:"type:text"`
	TakenAt      *time.Time `gorm:"type:timestamptz"`
	TripPlaceID  *string    `gorm:"type:uuid;index"`
	CreatedAt    time.Time  `gorm:"not null;default:now()"`
	UpdatedAt    time.Time  `gorm:"not null;default:now()"`
}

type PhotoTag struct {
	ID        string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	PhotoID   string    `gorm:"type:uuid;not null;uniqueIndex:ux_photo_tags_photo_user"`
	UserID    string    `gorm:"type:uuid;not null;uniqueIndex:ux_photo_tags_photo_user"`
	TaggedBy  string    `gorm:"type:uuid;not null;index"`
	CreatedAt time.Time `gorm:"not null;default:now()"`
}

type Comment struct {
	ID         string         `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	AuthorID   string         `gorm:"type:uuid;not null;index"`
	TargetType string         `gorm:"size:50;not null;index"`
	TargetID   string         `gorm:"type:uuid;not null;index"`
	ParentID   *string        `gorm:"type:uuid;index"`
	Content    string         `gorm:"type:text;not null"`
	CreatedAt  time.Time      `gorm:"not null;default:now()"`
	UpdatedAt  time.Time      `gorm:"not null;default:now()"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

type Reaction struct {
	ID           string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID       string    `gorm:"type:uuid;not null;uniqueIndex:ux_reactions_user_target_reaction"`
	TargetType   string    `gorm:"size:50;not null;uniqueIndex:ux_reactions_user_target_reaction"`
	TargetID     string    `gorm:"type:uuid;not null;uniqueIndex:ux_reactions_user_target_reaction"`
	ReactionType string    `gorm:"size:20;not null;uniqueIndex:ux_reactions_user_target_reaction"`
	CreatedAt    time.Time `gorm:"not null;default:now()"`
}

type Notification struct {
	ID        string         `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    string         `gorm:"type:uuid;not null;index"`
	Type      string         `gorm:"size:100;not null"`
	Data      datatypes.JSON `gorm:"type:jsonb;not null;default:'{}'"`
	IsRead    bool           `gorm:"not null;default:false"`
	CreatedAt time.Time      `gorm:"not null;default:now()"`
}
