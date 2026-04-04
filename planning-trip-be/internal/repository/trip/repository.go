package trip

import (
	"context"
	"time"

	"planning-trip-be/internal/model"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, trip *model.Trip) error
	GetByID(ctx context.Context, id string) (model.Trip, error)
	ListByUser(ctx context.Context, userID string) ([]model.Trip, error)
}

type gormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) Create(ctx context.Context, trip *model.Trip) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(trip).Error; err != nil {
			return err
		}

		member := model.TripMember{
			TripID:   trip.ID,
			UserID:   trip.CreatorID,
			Role:     "owner",
			JoinedAt: time.Now().UTC(),
		}
		return tx.Create(&member).Error
	})
}

func (r *gormRepository) GetByID(ctx context.Context, id string) (model.Trip, error) {
	var trip model.Trip
	err := r.db.WithContext(ctx).First(&trip, "id = ?", id).Error
	return trip, err
}

func (r *gormRepository) ListByUser(ctx context.Context, userID string) ([]model.Trip, error) {
	var trips []model.Trip
	err := r.db.WithContext(ctx).
		Model(&model.Trip{}).
		Distinct("trips.*").
		Joins("LEFT JOIN trip_members ON trip_members.trip_id = trips.id").
		Where("trips.creator_id = ? OR trip_members.user_id = ?", userID, userID).
		Order("trips.created_at DESC").
		Find(&trips).Error
	return trips, err
}
