package trip

import (
	"context"
	"errors"
	"strings"
	"time"

	"planning-trip-be/internal/model"
	repo "planning-trip-be/internal/repository/trip"

	"gorm.io/gorm"
)

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrNotFound     = gorm.ErrRecordNotFound
	ErrForbidden    = errors.New("forbidden")
)

const dateLayout = "2006-01-02"

type TripDTO struct {
	ID            string     `json:"id"`
	CreatorID     string     `json:"creator_id"`
	Title         string     `json:"title"`
	Description   string     `json:"description,omitempty"`
	StartDate     *time.Time `json:"start_date,omitempty"`
	EndDate       *time.Time `json:"end_date,omitempty"`
	Visibility    string     `json:"visibility"`
	Status        string     `json:"status"`
	CoverImageURL string     `json:"cover_image_url,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type TripMemberDTO struct {
	UserID    string    `json:"user_id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	AvatarURL string    `json:"avatar_url,omitempty"`
	Role      string    `json:"role"`
	JoinedAt  time.Time `json:"joined_at"`
}

type PlaceDTO struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Address string  `json:"address,omitempty"`
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
}

type TripPlaceDTO struct {
	ID         string   `json:"id"`
	Title      string   `json:"title,omitempty"`
	Note       string   `json:"note,omitempty"`
	DayIndex   int      `json:"day_index"`
	OrderIndex int      `json:"order_index"`
	Place      PlaceDTO `json:"place"`
}

type ScheduleItemDTO struct {
	ID          string        `json:"id"`
	TripID      string        `json:"trip_id"`
	TripPlaceID *string       `json:"trip_place_id,omitempty"`
	Title       string        `json:"title"`
	Description string        `json:"description,omitempty"`
	StartTime   *time.Time    `json:"start_time,omitempty"`
	EndTime     *time.Time    `json:"end_time,omitempty"`
	DayIndex    int           `json:"day_index"`
	OrderIndex  int           `json:"order_index"`
	CreatedBy   string        `json:"created_by"`
	TripPlace   *TripPlaceDTO `json:"trip_place,omitempty"`
}

type AlbumPhotoDTO struct {
	ID           string     `json:"id"`
	URL          string     `json:"url"`
	ThumbnailURL string     `json:"thumbnail_url,omitempty"`
	Caption      string     `json:"caption,omitempty"`
	TakenAt      *time.Time `json:"taken_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

type TripDetailDTO struct {
	TripDTO
	Members      []TripMemberDTO   `json:"members"`
	Schedule     []ScheduleItemDTO `json:"schedule"`
	AlbumPreview []AlbumPhotoDTO   `json:"album_preview"`
}

type CreateInput struct {
	CreatorID     string  `json:"creator_id"`
	Title         string  `json:"title"`
	Description   string  `json:"description"`
	StartDate     *string `json:"start_date"`
	EndDate       *string `json:"end_date"`
	Visibility    string  `json:"visibility"`
	Status        string  `json:"status"`
	CoverImageURL string  `json:"cover_image_url"`
}

type UpdateStatusInput struct {
	Status string `json:"status"`
}

type CreateScheduleItemInput struct {
	TripPlaceID *string    `json:"trip_place_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	StartTime   *time.Time `json:"start_time"`
	EndTime     *time.Time `json:"end_time"`
	DayIndex    int        `json:"day_index"`
	OrderIndex  int        `json:"order_index"`
}

type UpdateScheduleItemInput struct {
	TripPlaceID *string     `json:"trip_place_id"`
	Title       *string     `json:"title"`
	Description *string     `json:"description"`
	StartTime   **time.Time `json:"start_time"`
	EndTime     **time.Time `json:"end_time"`
	DayIndex    *int        `json:"day_index"`
	OrderIndex  *int        `json:"order_index"`
}

type ReorderScheduleItem struct {
	ID         string `json:"id"`
	DayIndex   int    `json:"day_index"`
	OrderIndex int    `json:"order_index"`
}

type ReorderScheduleInput struct {
	Items []ReorderScheduleItem `json:"items"`
}

type Service interface {
	Create(ctx context.Context, input CreateInput) (TripDTO, error)
	GetByID(ctx context.Context, id string) (TripDTO, error)
	GetDetail(ctx context.Context, tripID, userID string) (TripDetailDTO, error)
	ListByUser(ctx context.Context, userID string) ([]TripDTO, error)
	UpdateStatus(ctx context.Context, tripID, userID string, input UpdateStatusInput) (TripDTO, error)
	ListMembers(ctx context.Context, tripID, userID string) ([]TripMemberDTO, error)
	ListScheduleItems(ctx context.Context, tripID, userID string, dayIndex *int) ([]ScheduleItemDTO, error)
	CreateScheduleItem(ctx context.Context, tripID, userID string, input CreateScheduleItemInput) (ScheduleItemDTO, error)
	UpdateScheduleItem(ctx context.Context, tripID, itemID, userID string, input UpdateScheduleItemInput) (ScheduleItemDTO, error)
	DeleteScheduleItem(ctx context.Context, tripID, itemID, userID string) error
	ReorderSchedule(ctx context.Context, tripID, userID string, input ReorderScheduleInput) error
	ListAlbumPreview(ctx context.Context, tripID, userID string, limit int) ([]AlbumPhotoDTO, error)
}

type service struct {
	repo repo.Repository
}

func NewService(repository repo.Repository) Service {
	return &service{repo: repository}
}

func (s *service) Create(ctx context.Context, input CreateInput) (TripDTO, error) {
	creatorID := strings.TrimSpace(input.CreatorID)
	title := strings.TrimSpace(input.Title)
	if creatorID == "" || title == "" {
		return TripDTO{}, ErrInvalidInput
	}

	startDate, err := parseDatePtr(input.StartDate)
	if err != nil {
		return TripDTO{}, err
	}
	endDate, err := parseDatePtr(input.EndDate)
	if err != nil {
		return TripDTO{}, err
	}
	if startDate != nil && endDate != nil && startDate.After(*endDate) {
		return TripDTO{}, ErrInvalidInput
	}

	visibility := strings.ToLower(strings.TrimSpace(input.Visibility))
	if visibility == "" {
		visibility = "private"
	}

	status := strings.ToLower(strings.TrimSpace(input.Status))
	if status == "" {
		status = "draft"
	}

	trip := model.Trip{
		CreatorID:     creatorID,
		Title:         title,
		Description:   strings.TrimSpace(input.Description),
		StartDate:     startDate,
		EndDate:       endDate,
		Visibility:    visibility,
		Status:        status,
		CoverImageURL: strings.TrimSpace(input.CoverImageURL),
	}
	if err := s.repo.Create(ctx, &trip); err != nil {
		return TripDTO{}, err
	}
	return toDTO(trip), nil
}

func (s *service) GetByID(ctx context.Context, id string) (TripDTO, error) {
	tripID := strings.TrimSpace(id)
	if tripID == "" {
		return TripDTO{}, ErrInvalidInput
	}

	trip, err := s.repo.GetByID(ctx, tripID)
	if err != nil {
		return TripDTO{}, err
	}
	return toDTO(trip), nil
}

func (s *service) GetDetail(ctx context.Context, tripID, userID string) (TripDetailDTO, error) {
	if strings.TrimSpace(tripID) == "" || strings.TrimSpace(userID) == "" {
		return TripDetailDTO{}, ErrInvalidInput
	}

	if err := s.ensureCanView(ctx, tripID, userID); err != nil {
		return TripDetailDTO{}, err
	}

	trip, err := s.repo.GetByID(ctx, tripID)
	if err != nil {
		return TripDetailDTO{}, err
	}

	members, err := s.repo.ListMembers(ctx, tripID)
	if err != nil {
		return TripDetailDTO{}, err
	}
	schedule, err := s.repo.ListScheduleItems(ctx, tripID, nil)
	if err != nil {
		return TripDetailDTO{}, err
	}
	photos, err := s.repo.ListAlbumPreview(ctx, tripID, 12)
	if err != nil {
		return TripDetailDTO{}, err
	}

	return TripDetailDTO{
		TripDTO:      toDTO(trip),
		Members:      toMemberDTOs(members),
		Schedule:     toScheduleItemDTOs(schedule),
		AlbumPreview: toAlbumPhotoDTOs(photos),
	}, nil
}

func (s *service) ListByUser(ctx context.Context, userID string) ([]TripDTO, error) {
	trimmedUserID := strings.TrimSpace(userID)
	if trimmedUserID == "" {
		return nil, ErrInvalidInput
	}

	trips, err := s.repo.ListByUser(ctx, trimmedUserID)
	if err != nil {
		return nil, err
	}

	result := make([]TripDTO, 0, len(trips))
	for _, trip := range trips {
		result = append(result, toDTO(trip))
	}
	return result, nil
}

func (s *service) UpdateStatus(ctx context.Context, tripID, userID string, input UpdateStatusInput) (TripDTO, error) {
	if err := s.ensureCanEdit(ctx, tripID, userID); err != nil {
		return TripDTO{}, err
	}

	status := strings.ToLower(strings.TrimSpace(input.Status))
	if status != "draft" && status != "published" {
		return TripDTO{}, ErrInvalidInput
	}

	trip, err := s.repo.UpdateStatus(ctx, tripID, status)
	if err != nil {
		return TripDTO{}, err
	}
	return toDTO(trip), nil
}

func (s *service) ListMembers(ctx context.Context, tripID, userID string) ([]TripMemberDTO, error) {
	if err := s.ensureCanView(ctx, tripID, userID); err != nil {
		return nil, err
	}
	members, err := s.repo.ListMembers(ctx, tripID)
	if err != nil {
		return nil, err
	}
	return toMemberDTOs(members), nil
}

func (s *service) ListScheduleItems(ctx context.Context, tripID, userID string, dayIndex *int) ([]ScheduleItemDTO, error) {
	if err := s.ensureCanView(ctx, tripID, userID); err != nil {
		return nil, err
	}
	if dayIndex != nil && *dayIndex <= 0 {
		return nil, ErrInvalidInput
	}

	items, err := s.repo.ListScheduleItems(ctx, tripID, dayIndex)
	if err != nil {
		return nil, err
	}
	return toScheduleItemDTOs(items), nil
}

func (s *service) CreateScheduleItem(ctx context.Context, tripID, userID string, input CreateScheduleItemInput) (ScheduleItemDTO, error) {
	if err := s.ensureCanEdit(ctx, tripID, userID); err != nil {
		return ScheduleItemDTO{}, err
	}

	title := strings.TrimSpace(input.Title)
	if title == "" || input.DayIndex <= 0 {
		return ScheduleItemDTO{}, ErrInvalidInput
	}

	item := model.TripScheduleItem{
		TripID:      tripID,
		TripPlaceID: cleanStringPtr(input.TripPlaceID),
		Title:       title,
		Description: strings.TrimSpace(input.Description),
		StartTime:   input.StartTime,
		EndTime:     input.EndTime,
		DayIndex:    input.DayIndex,
		OrderIndex:  input.OrderIndex,
		CreatedBy:   userID,
	}

	if err := s.repo.CreateScheduleItem(ctx, &item); err != nil {
		return ScheduleItemDTO{}, err
	}

	items, err := s.repo.ListScheduleItems(ctx, tripID, nil)
	if err != nil {
		return ScheduleItemDTO{}, err
	}
	for _, row := range items {
		if row.Item.ID == item.ID {
			mapped := toScheduleItemDTO(row)
			return mapped, nil
		}
	}
	return ScheduleItemDTO{}, gorm.ErrRecordNotFound
}

func (s *service) UpdateScheduleItem(ctx context.Context, tripID, itemID, userID string, input UpdateScheduleItemInput) (ScheduleItemDTO, error) {
	if err := s.ensureCanEdit(ctx, tripID, userID); err != nil {
		return ScheduleItemDTO{}, err
	}
	if strings.TrimSpace(itemID) == "" {
		return ScheduleItemDTO{}, ErrInvalidInput
	}

	updates := map[string]interface{}{}
	if input.TripPlaceID != nil {
		updates["trip_place_id"] = cleanStringPtr(input.TripPlaceID)
	}
	if input.Title != nil {
		title := strings.TrimSpace(*input.Title)
		if title == "" {
			return ScheduleItemDTO{}, ErrInvalidInput
		}
		updates["title"] = title
	}
	if input.Description != nil {
		updates["description"] = strings.TrimSpace(*input.Description)
	}
	if input.StartTime != nil {
		updates["start_time"] = *input.StartTime
	}
	if input.EndTime != nil {
		updates["end_time"] = *input.EndTime
	}
	if input.DayIndex != nil {
		if *input.DayIndex <= 0 {
			return ScheduleItemDTO{}, ErrInvalidInput
		}
		updates["day_index"] = *input.DayIndex
	}
	if input.OrderIndex != nil {
		updates["order_index"] = *input.OrderIndex
	}
	if len(updates) == 0 {
		return ScheduleItemDTO{}, ErrInvalidInput
	}

	item, err := s.repo.UpdateScheduleItem(ctx, tripID, itemID, updates)
	if err != nil {
		return ScheduleItemDTO{}, err
	}

	items, err := s.repo.ListScheduleItems(ctx, tripID, nil)
	if err != nil {
		return ScheduleItemDTO{}, err
	}
	for _, row := range items {
		if row.Item.ID == item.ID {
			mapped := toScheduleItemDTO(row)
			return mapped, nil
		}
	}
	return ScheduleItemDTO{}, gorm.ErrRecordNotFound
}

func (s *service) DeleteScheduleItem(ctx context.Context, tripID, itemID, userID string) error {
	if err := s.ensureCanEdit(ctx, tripID, userID); err != nil {
		return err
	}
	if strings.TrimSpace(itemID) == "" {
		return ErrInvalidInput
	}
	return s.repo.DeleteScheduleItem(ctx, tripID, itemID)
}

func (s *service) ReorderSchedule(ctx context.Context, tripID, userID string, input ReorderScheduleInput) error {
	if err := s.ensureCanEdit(ctx, tripID, userID); err != nil {
		return err
	}
	if len(input.Items) == 0 {
		return ErrInvalidInput
	}

	items := make([]repo.ScheduleItemReorder, 0, len(input.Items))
	for _, row := range input.Items {
		if strings.TrimSpace(row.ID) == "" || row.DayIndex <= 0 {
			return ErrInvalidInput
		}
		items = append(items, repo.ScheduleItemReorder{
			ID:         row.ID,
			DayIndex:   row.DayIndex,
			OrderIndex: row.OrderIndex,
		})
	}

	return s.repo.ReorderScheduleItems(ctx, tripID, items)
}

func (s *service) ListAlbumPreview(ctx context.Context, tripID, userID string, limit int) ([]AlbumPhotoDTO, error) {
	if err := s.ensureCanView(ctx, tripID, userID); err != nil {
		return nil, err
	}
	photos, err := s.repo.ListAlbumPreview(ctx, tripID, limit)
	if err != nil {
		return nil, err
	}
	return toAlbumPhotoDTOs(photos), nil
}

func (s *service) ensureCanView(ctx context.Context, tripID, userID string) error {
	if strings.TrimSpace(tripID) == "" || strings.TrimSpace(userID) == "" {
		return ErrInvalidInput
	}
	ok, err := s.repo.CanViewTrip(ctx, tripID, userID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrForbidden
	}
	return nil
}

func (s *service) ensureCanEdit(ctx context.Context, tripID, userID string) error {
	if strings.TrimSpace(tripID) == "" || strings.TrimSpace(userID) == "" {
		return ErrInvalidInput
	}
	ok, err := s.repo.CanEditTrip(ctx, tripID, userID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrForbidden
	}
	return nil
}

func parseDatePtr(value *string) (*time.Time, error) {
	if value == nil {
		return nil, nil
	}

	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil, nil
	}

	date, err := time.Parse(dateLayout, trimmed)
	if err != nil {
		return nil, ErrInvalidInput
	}
	return &date, nil
}

func cleanStringPtr(input *string) *string {
	if input == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*input)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func toDTO(trip model.Trip) TripDTO {
	return TripDTO{
		ID:            trip.ID,
		CreatorID:     trip.CreatorID,
		Title:         trip.Title,
		Description:   trip.Description,
		StartDate:     trip.StartDate,
		EndDate:       trip.EndDate,
		Visibility:    trip.Visibility,
		Status:        trip.Status,
		CoverImageURL: trip.CoverImageURL,
		CreatedAt:     trip.CreatedAt,
		UpdatedAt:     trip.UpdatedAt,
	}
}

func toMemberDTOs(rows []repo.TripMemberAggregate) []TripMemberDTO {
	result := make([]TripMemberDTO, 0, len(rows))
	for _, row := range rows {
		result = append(result, TripMemberDTO{
			UserID:    row.User.ID,
			FullName:  row.User.FullName,
			Email:     row.User.Email,
			AvatarURL: row.User.AvatarURL,
			Role:      row.Member.Role,
			JoinedAt:  row.Member.JoinedAt,
		})
	}
	return result
}

func toScheduleItemDTOs(rows []repo.ScheduleItemAggregate) []ScheduleItemDTO {
	result := make([]ScheduleItemDTO, 0, len(rows))
	for _, row := range rows {
		result = append(result, toScheduleItemDTO(row))
	}
	return result
}

func toScheduleItemDTO(row repo.ScheduleItemAggregate) ScheduleItemDTO {
	item := row.Item
	dto := ScheduleItemDTO{
		ID:          item.ID,
		TripID:      item.TripID,
		TripPlaceID: item.TripPlaceID,
		Title:       item.Title,
		Description: item.Description,
		StartTime:   item.StartTime,
		EndTime:     item.EndTime,
		DayIndex:    item.DayIndex,
		OrderIndex:  item.OrderIndex,
		CreatedBy:   item.CreatedBy,
	}

	if row.TripPlace != nil && row.Place != nil {
		dto.TripPlace = &TripPlaceDTO{
			ID:         row.TripPlace.ID,
			Title:      row.TripPlace.Title,
			Note:       row.TripPlace.Note,
			DayIndex:   row.TripPlace.DayIndex,
			OrderIndex: row.TripPlace.OrderIndex,
			Place: PlaceDTO{
				ID:      row.Place.ID,
				Name:    row.Place.Name,
				Address: row.Place.Address,
				Lat:     row.Place.Lat,
				Lng:     row.Place.Lng,
			},
		}
	}

	return dto
}

func toAlbumPhotoDTOs(photos []model.Photo) []AlbumPhotoDTO {
	result := make([]AlbumPhotoDTO, 0, len(photos))
	for _, photo := range photos {
		result = append(result, AlbumPhotoDTO{
			ID:           photo.ID,
			URL:          photo.URL,
			ThumbnailURL: photo.ThumbnailURL,
			Caption:      photo.Caption,
			TakenAt:      photo.TakenAt,
			CreatedAt:    photo.CreatedAt,
		})
	}
	return result
}
