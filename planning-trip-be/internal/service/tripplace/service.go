package tripplace

import (
	"context"
	"errors"
	"strings"
	"time"

	"planning-trip-be/internal/model"
	plrepo "planning-trip-be/internal/repository/place"
	trepo "planning-trip-be/internal/repository/trip"
	repo "planning-trip-be/internal/repository/tripplace"

	"gorm.io/gorm"
)

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrForbidden    = errors.New("forbidden")
	ErrNotFound     = gorm.ErrRecordNotFound
)

type PlaceDTO struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Address string  `json:"address,omitempty"`
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
}

type TripPlaceDTO struct {
	ID         string    `json:"id"`
	TripID     string    `json:"trip_id"`
	PlaceID    string    `json:"place_id"`
	Title      string    `json:"title,omitempty"`
	Note       string    `json:"note,omitempty"`
	DayIndex   int       `json:"day_index"`
	OrderIndex int       `json:"order_index"`
	CreatedBy  string    `json:"created_by"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Place      *PlaceDTO `json:"place,omitempty"`
}

type CreateInput struct {
	TripID     string `json:"trip_id"`
	PlaceID    string `json:"place_id"`
	Title      string `json:"title"`
	Note       string `json:"note"`
	DayIndex   int    `json:"day_index"`
	OrderIndex int    `json:"order_index"`
	CreatedBy  string `json:"created_by"`
}

type UpdateInput struct {
	Title      *string `json:"title"`
	Note       *string `json:"note"`
	DayIndex   *int    `json:"day_index"`
	OrderIndex *int    `json:"order_index"`
}

type Service interface {
	ListByTrip(ctx context.Context, tripID, userID string, dayIndex *int) ([]TripPlaceDTO, error)
	GetByID(ctx context.Context, id, userID string) (TripPlaceDTO, error)
	Create(ctx context.Context, input CreateInput) (TripPlaceDTO, error)
	Update(ctx context.Context, id, userID string, input UpdateInput) (TripPlaceDTO, error)
	Delete(ctx context.Context, id, userID string) error
}

type service struct {
	repo      repo.Repository
	tripRepo  trepo.Repository
	placeRepo plrepo.Repository
}

func NewService(repository repo.Repository, tripRepo trepo.Repository, placeRepo plrepo.Repository) Service {
	return &service{repo: repository, tripRepo: tripRepo, placeRepo: placeRepo}
}

func (s *service) ListByTrip(ctx context.Context, tripID, userID string, dayIndex *int) ([]TripPlaceDTO, error) {
	if strings.TrimSpace(tripID) == "" || strings.TrimSpace(userID) == "" {
		return nil, ErrInvalidInput
	}
	if err := s.ensureCanView(ctx, tripID, userID); err != nil {
		return nil, err
	}
	if dayIndex != nil && *dayIndex <= 0 {
		return nil, ErrInvalidInput
	}

	rows, err := s.repo.ListByTrip(ctx, tripID, dayIndex)
	if err != nil {
		return nil, err
	}
	result := make([]TripPlaceDTO, 0, len(rows))
	for _, row := range rows {
		result = append(result, s.toDTO(ctx, row))
	}
	return result, nil
}

func (s *service) GetByID(ctx context.Context, id, userID string) (TripPlaceDTO, error) {
	if strings.TrimSpace(id) == "" || strings.TrimSpace(userID) == "" {
		return TripPlaceDTO{}, ErrInvalidInput
	}
	row, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return TripPlaceDTO{}, err
	}
	if err := s.ensureCanView(ctx, row.TripID, userID); err != nil {
		return TripPlaceDTO{}, err
	}
	return s.toDTO(ctx, row), nil
}

func (s *service) Create(ctx context.Context, input CreateInput) (TripPlaceDTO, error) {
	tripID := strings.TrimSpace(input.TripID)
	placeID := strings.TrimSpace(input.PlaceID)
	createdBy := strings.TrimSpace(input.CreatedBy)
	if tripID == "" || placeID == "" || createdBy == "" || input.DayIndex <= 0 {
		return TripPlaceDTO{}, ErrInvalidInput
	}
	if err := s.ensureCanEdit(ctx, tripID, createdBy); err != nil {
		return TripPlaceDTO{}, err
	}
	if _, err := s.placeRepo.GetByID(ctx, placeID); err != nil {
		return TripPlaceDTO{}, err
	}

	row := model.TripPlace{
		TripID:     tripID,
		PlaceID:    placeID,
		Title:      strings.TrimSpace(input.Title),
		Note:       strings.TrimSpace(input.Note),
		DayIndex:   input.DayIndex,
		OrderIndex: input.OrderIndex,
		CreatedBy:  createdBy,
	}
	if err := s.repo.Create(ctx, &row); err != nil {
		return TripPlaceDTO{}, err
	}
	return s.toDTO(ctx, row), nil
}

func (s *service) Update(ctx context.Context, id, userID string, input UpdateInput) (TripPlaceDTO, error) {
	if strings.TrimSpace(id) == "" || strings.TrimSpace(userID) == "" {
		return TripPlaceDTO{}, ErrInvalidInput
	}
	row, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return TripPlaceDTO{}, err
	}
	if err := s.ensureCanEdit(ctx, row.TripID, userID); err != nil {
		return TripPlaceDTO{}, err
	}

	updates := map[string]interface{}{}
	if input.Title != nil {
		updates["title"] = strings.TrimSpace(*input.Title)
	}
	if input.Note != nil {
		updates["note"] = strings.TrimSpace(*input.Note)
	}
	if input.DayIndex != nil {
		if *input.DayIndex <= 0 {
			return TripPlaceDTO{}, ErrInvalidInput
		}
		updates["day_index"] = *input.DayIndex
	}
	if input.OrderIndex != nil {
		updates["order_index"] = *input.OrderIndex
	}
	if len(updates) == 0 {
		return TripPlaceDTO{}, ErrInvalidInput
	}

	updated, err := s.repo.Update(ctx, id, updates)
	if err != nil {
		return TripPlaceDTO{}, err
	}
	return s.toDTO(ctx, updated), nil
}

func (s *service) Delete(ctx context.Context, id, userID string) error {
	if strings.TrimSpace(id) == "" || strings.TrimSpace(userID) == "" {
		return ErrInvalidInput
	}
	row, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if err := s.ensureCanEdit(ctx, row.TripID, userID); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}

func (s *service) ensureCanView(ctx context.Context, tripID, userID string) error {
	ok, err := s.tripRepo.CanViewTrip(ctx, tripID, userID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrForbidden
	}
	return nil
}

func (s *service) ensureCanEdit(ctx context.Context, tripID, userID string) error {
	ok, err := s.tripRepo.CanEditTrip(ctx, tripID, userID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrForbidden
	}
	return nil
}

func (s *service) toDTO(ctx context.Context, row model.TripPlace) TripPlaceDTO {
	dto := TripPlaceDTO{
		ID:         row.ID,
		TripID:     row.TripID,
		PlaceID:    row.PlaceID,
		Title:      row.Title,
		Note:       row.Note,
		DayIndex:   row.DayIndex,
		OrderIndex: row.OrderIndex,
		CreatedBy:  row.CreatedBy,
		CreatedAt:  row.CreatedAt,
		UpdatedAt:  row.UpdatedAt,
	}
	place, err := s.placeRepo.GetByID(ctx, row.PlaceID)
	if err == nil {
		dto.Place = &PlaceDTO{
			ID:      place.ID,
			Name:    place.Name,
			Address: place.Address,
			Lat:     place.Lat,
			Lng:     place.Lng,
		}
	}
	return dto
}
