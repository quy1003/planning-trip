package tripplace

import (
	"context"
	"time"

	"planning-trip-be/internal/model"

	"gorm.io/gorm"
)

type Repository interface {
	ListByTrip(ctx context.Context, tripID string, dayIndex *int) ([]model.TripPlace, error)
	GetByID(ctx context.Context, id string) (model.TripPlace, error)
	Create(ctx context.Context, item *model.TripPlace) error
	Update(ctx context.Context, id string, updates map[string]interface{}) (model.TripPlace, error)
	Delete(ctx context.Context, id string) error
}

type gormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) ListByTrip(ctx context.Context, tripID string, dayIndex *int) ([]model.TripPlace, error) {
	query := r.db.WithContext(ctx).Where("trip_id = ?", tripID).Order("day_index ASC, order_index ASC, created_at ASC")
	if dayIndex != nil {
		query = query.Where("day_index = ?", *dayIndex)
	}
	var rows []model.TripPlace
	err := query.Find(&rows).Error
	return rows, err
}

func (r *gormRepository) GetByID(ctx context.Context, id string) (model.TripPlace, error) {
	var row model.TripPlace
	err := r.db.WithContext(ctx).First(&row, "id = ?", id).Error
	return row, err
}

func (r *gormRepository) Create(ctx context.Context, item *model.TripPlace) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *gormRepository) Update(ctx context.Context, id string, updates map[string]interface{}) (model.TripPlace, error) {
	updates["updated_at"] = time.Now().UTC()
	result := r.db.WithContext(ctx).Model(&model.TripPlace{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return model.TripPlace{}, result.Error
	}
	if result.RowsAffected == 0 {
		return model.TripPlace{}, gorm.ErrRecordNotFound
	}
	return r.GetByID(ctx, id)
}

func (r *gormRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&model.TripPlace{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
