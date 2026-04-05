package place

import (
	"context"
	"strings"
	"time"

	"planning-trip-be/internal/model"

	"gorm.io/gorm"
)

type Repository interface {
	List(ctx context.Context, keyword string, limit int) ([]model.Place, error)
	GetByID(ctx context.Context, id string) (model.Place, error)
	Create(ctx context.Context, place *model.Place) error
	Update(ctx context.Context, id string, updates map[string]interface{}) (model.Place, error)
	Delete(ctx context.Context, id string) error
}

type gormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) List(ctx context.Context, keyword string, limit int) ([]model.Place, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	query := r.db.WithContext(ctx).Model(&model.Place{})
	trimmed := strings.TrimSpace(keyword)
	if trimmed != "" {
		query = query.Where("name ILIKE ? OR address ILIKE ?", "%"+trimmed+"%", "%"+trimmed+"%")
	}

	var places []model.Place
	err := query.Order("created_at DESC").Limit(limit).Find(&places).Error
	return places, err
}

func (r *gormRepository) GetByID(ctx context.Context, id string) (model.Place, error) {
	var place model.Place
	err := r.db.WithContext(ctx).First(&place, "id = ?", id).Error
	return place, err
}

func (r *gormRepository) Create(ctx context.Context, place *model.Place) error {
	return r.db.WithContext(ctx).Create(place).Error
}

func (r *gormRepository) Update(ctx context.Context, id string, updates map[string]interface{}) (model.Place, error) {
	updates["updated_at"] = time.Now().UTC()
	result := r.db.WithContext(ctx).Model(&model.Place{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return model.Place{}, result.Error
	}
	if result.RowsAffected == 0 {
		return model.Place{}, gorm.ErrRecordNotFound
	}
	return r.GetByID(ctx, id)
}

func (r *gormRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&model.Place{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
