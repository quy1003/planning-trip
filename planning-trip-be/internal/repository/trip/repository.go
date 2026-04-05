package trip

import (
	"context"
	"errors"
	"time"

	"planning-trip-be/internal/model"

	"gorm.io/gorm"
)

type ScheduleItemReorder struct {
	ID         string
	DayIndex   int
	OrderIndex int
}

type ScheduleItemAggregate struct {
	Item      model.TripScheduleItem
	TripPlace *model.TripPlace
	Place     *model.Place
}

type TripMemberAggregate struct {
	Member model.TripMember
	User   model.User
}

type Repository interface {
	Create(ctx context.Context, trip *model.Trip) error
	GetByID(ctx context.Context, id string) (model.Trip, error)
	ListByUser(ctx context.Context, userID string) ([]model.Trip, error)
	CanViewTrip(ctx context.Context, tripID, userID string) (bool, error)
	CanEditTrip(ctx context.Context, tripID, userID string) (bool, error)
	UpdateStatus(ctx context.Context, tripID, status string) (model.Trip, error)
	ListMembers(ctx context.Context, tripID string) ([]TripMemberAggregate, error)
	ListScheduleItems(ctx context.Context, tripID string, dayIndex *int) ([]ScheduleItemAggregate, error)
	CreateScheduleItem(ctx context.Context, item *model.TripScheduleItem) error
	UpdateScheduleItem(ctx context.Context, tripID, itemID string, updates map[string]interface{}) (model.TripScheduleItem, error)
	DeleteScheduleItem(ctx context.Context, tripID, itemID string) error
	ReorderScheduleItems(ctx context.Context, tripID string, items []ScheduleItemReorder) error
	ListAlbumPreview(ctx context.Context, tripID string, limit int) ([]model.Photo, error)
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

func (r *gormRepository) CanViewTrip(ctx context.Context, tripID, userID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.Trip{}).
		Joins("LEFT JOIN trip_members ON trip_members.trip_id = trips.id").
		Where("trips.id = ? AND (trips.creator_id = ? OR trip_members.user_id = ?)", tripID, userID, userID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *gormRepository) CanEditTrip(ctx context.Context, tripID, userID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("trips").
		Joins("LEFT JOIN trip_members ON trip_members.trip_id = trips.id").
		Where("trips.id = ?", tripID).
		Where("trips.creator_id = ? OR (trip_members.user_id = ? AND trip_members.role IN ?)", userID, userID, []string{"owner", "editor"}).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *gormRepository) UpdateStatus(ctx context.Context, tripID, status string) (model.Trip, error) {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now().UTC(),
	}

	result := r.db.WithContext(ctx).Model(&model.Trip{}).Where("id = ?", tripID).Updates(updates)
	if result.Error != nil {
		return model.Trip{}, result.Error
	}
	if result.RowsAffected == 0 {
		return model.Trip{}, gorm.ErrRecordNotFound
	}
	return r.GetByID(ctx, tripID)
}

func (r *gormRepository) ListMembers(ctx context.Context, tripID string) ([]TripMemberAggregate, error) {
	var members []model.TripMember
	if err := r.db.WithContext(ctx).
		Where("trip_id = ?", tripID).
		Order("joined_at ASC").
		Find(&members).Error; err != nil {
		return nil, err
	}

	if len(members) == 0 {
		return []TripMemberAggregate{}, nil
	}

	userIDs := make([]string, 0, len(members))
	for _, member := range members {
		userIDs = append(userIDs, member.UserID)
	}

	var users []model.User
	if err := r.db.WithContext(ctx).Where("id IN ?", userIDs).Find(&users).Error; err != nil {
		return nil, err
	}

	userByID := make(map[string]model.User, len(users))
	for _, user := range users {
		userByID[user.ID] = user
	}

	result := make([]TripMemberAggregate, 0, len(members))
	for _, member := range members {
		user, ok := userByID[member.UserID]
		if !ok {
			continue
		}
		result = append(result, TripMemberAggregate{
			Member: member,
			User:   user,
		})
	}
	return result, nil
}

func (r *gormRepository) ListScheduleItems(ctx context.Context, tripID string, dayIndex *int) ([]ScheduleItemAggregate, error) {
	query := r.db.WithContext(ctx).
		Where("trip_id = ?", tripID).
		Order("day_index ASC, order_index ASC, start_time ASC")

	if dayIndex != nil {
		query = query.Where("day_index = ?", *dayIndex)
	}

	var items []model.TripScheduleItem
	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return []ScheduleItemAggregate{}, nil
	}

	tripPlaceIDs := make([]string, 0, len(items))
	for _, item := range items {
		if item.TripPlaceID != nil {
			tripPlaceIDs = append(tripPlaceIDs, *item.TripPlaceID)
		}
	}

	tripPlaceByID := map[string]model.TripPlace{}
	placeByID := map[string]model.Place{}

	if len(tripPlaceIDs) > 0 {
		var tripPlaces []model.TripPlace
		if err := r.db.WithContext(ctx).Where("id IN ?", tripPlaceIDs).Find(&tripPlaces).Error; err != nil {
			return nil, err
		}
		placeIDs := make([]string, 0, len(tripPlaces))
		for _, tripPlace := range tripPlaces {
			tripPlaceByID[tripPlace.ID] = tripPlace
			placeIDs = append(placeIDs, tripPlace.PlaceID)
		}

		if len(placeIDs) > 0 {
			var places []model.Place
			if err := r.db.WithContext(ctx).Where("id IN ?", placeIDs).Find(&places).Error; err != nil {
				return nil, err
			}
			for _, place := range places {
				placeByID[place.ID] = place
			}
		}
	}

	result := make([]ScheduleItemAggregate, 0, len(items))
	for _, item := range items {
		agg := ScheduleItemAggregate{Item: item}
		if item.TripPlaceID != nil {
			if tripPlace, ok := tripPlaceByID[*item.TripPlaceID]; ok {
				t := tripPlace
				agg.TripPlace = &t
				if place, ok := placeByID[tripPlace.PlaceID]; ok {
					p := place
					agg.Place = &p
				}
			}
		}
		result = append(result, agg)
	}

	return result, nil
}

func (r *gormRepository) CreateScheduleItem(ctx context.Context, item *model.TripScheduleItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *gormRepository) UpdateScheduleItem(ctx context.Context, tripID, itemID string, updates map[string]interface{}) (model.TripScheduleItem, error) {
	updates["updated_at"] = time.Now().UTC()
	result := r.db.WithContext(ctx).Model(&model.TripScheduleItem{}).Where("id = ? AND trip_id = ?", itemID, tripID).Updates(updates)
	if result.Error != nil {
		return model.TripScheduleItem{}, result.Error
	}
	if result.RowsAffected == 0 {
		return model.TripScheduleItem{}, gorm.ErrRecordNotFound
	}

	var item model.TripScheduleItem
	err := r.db.WithContext(ctx).First(&item, "id = ?", itemID).Error
	return item, err
}

func (r *gormRepository) DeleteScheduleItem(ctx context.Context, tripID, itemID string) error {
	result := r.db.WithContext(ctx).Delete(&model.TripScheduleItem{}, "id = ? AND trip_id = ?", itemID, tripID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *gormRepository) ReorderScheduleItems(ctx context.Context, tripID string, items []ScheduleItemReorder) error {
	if len(items) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			updates := map[string]interface{}{
				"day_index":   item.DayIndex,
				"order_index": item.OrderIndex,
				"updated_at":  time.Now().UTC(),
			}
			result := tx.Model(&model.TripScheduleItem{}).Where("id = ? AND trip_id = ?", item.ID, tripID).Updates(updates)
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return gorm.ErrRecordNotFound
			}
		}
		return nil
	})
}

func (r *gormRepository) ListAlbumPreview(ctx context.Context, tripID string, limit int) ([]model.Photo, error) {
	if limit <= 0 {
		limit = 12
	}

	var photos []model.Photo
	err := r.db.WithContext(ctx).
		Where("trip_id = ?", tripID).
		Order("created_at DESC").
		Limit(limit).
		Find(&photos).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []model.Photo{}, nil
		}
		return nil, err
	}
	return photos, nil
}
