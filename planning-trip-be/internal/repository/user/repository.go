package user

import (
	"context"

	"planning-trip-be/internal/model"

	"gorm.io/gorm"
)

type Repository interface {
	List(ctx context.Context) ([]model.User, error)
	GetByID(ctx context.Context, id string) (model.User, error)
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, id string, updates map[string]interface{}) (model.User, error)
	Delete(ctx context.Context, id string) error
}

type gormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) List(ctx context.Context) ([]model.User, error) {
	var users []model.User
	err := r.db.WithContext(ctx).Order("created_at desc").Find(&users).Error
	return users, err
}

func (r *gormRepository) GetByID(ctx context.Context, id string) (model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error
	return user, err
}

func (r *gormRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *gormRepository) Update(ctx context.Context, id string, updates map[string]interface{}) (model.User, error) {
	if err := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return model.User{}, err
	}
	return r.GetByID(ctx, id)
}

func (r *gormRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&model.User{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
