package place

import (
	"context"
	"errors"
	"strings"
	"time"

	"planning-trip-be/internal/model"
	repo "planning-trip-be/internal/repository/place"

	"gorm.io/gorm"
)

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrNotFound     = gorm.ErrRecordNotFound
)

type PlaceDTO struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Address       string    `json:"address,omitempty"`
	Lat           float64   `json:"lat"`
	Lng           float64   `json:"lng"`
	GooglePlaceID string    `json:"google_place_id,omitempty"`
	CreatedBy     *string   `json:"created_by,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type CreateInput struct {
	Name          string  `json:"name"`
	Address       string  `json:"address"`
	Lat           float64 `json:"lat"`
	Lng           float64 `json:"lng"`
	GooglePlaceID string  `json:"google_place_id"`
	CreatedBy     *string `json:"created_by"`
}

type UpdateInput struct {
	Name          *string  `json:"name"`
	Address       *string  `json:"address"`
	Lat           *float64 `json:"lat"`
	Lng           *float64 `json:"lng"`
	GooglePlaceID *string  `json:"google_place_id"`
}

type Service interface {
	List(ctx context.Context, keyword string, limit int) ([]PlaceDTO, error)
	GetByID(ctx context.Context, id string) (PlaceDTO, error)
	Create(ctx context.Context, input CreateInput) (PlaceDTO, error)
	Update(ctx context.Context, id string, input UpdateInput) (PlaceDTO, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo repo.Repository
}

func NewService(repository repo.Repository) Service {
	return &service{repo: repository}
}

func (s *service) List(ctx context.Context, keyword string, limit int) ([]PlaceDTO, error) {
	rows, err := s.repo.List(ctx, keyword, limit)
	if err != nil {
		return nil, err
	}
	result := make([]PlaceDTO, 0, len(rows))
	for _, row := range rows {
		result = append(result, toDTO(row))
	}
	return result, nil
}

func (s *service) GetByID(ctx context.Context, id string) (PlaceDTO, error) {
	if strings.TrimSpace(id) == "" {
		return PlaceDTO{}, ErrInvalidInput
	}
	row, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return PlaceDTO{}, err
	}
	return toDTO(row), nil
}

func (s *service) Create(ctx context.Context, input CreateInput) (PlaceDTO, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return PlaceDTO{}, ErrInvalidInput
	}
	if input.Lat < -90 || input.Lat > 90 || input.Lng < -180 || input.Lng > 180 {
		return PlaceDTO{}, ErrInvalidInput
	}

	row := model.Place{
		Name:          name,
		Address:       strings.TrimSpace(input.Address),
		Lat:           input.Lat,
		Lng:           input.Lng,
		GooglePlaceID: strings.TrimSpace(input.GooglePlaceID),
		CreatedBy:     input.CreatedBy,
	}
	if err := s.repo.Create(ctx, &row); err != nil {
		return PlaceDTO{}, err
	}
	return toDTO(row), nil
}

func (s *service) Update(ctx context.Context, id string, input UpdateInput) (PlaceDTO, error) {
	if strings.TrimSpace(id) == "" {
		return PlaceDTO{}, ErrInvalidInput
	}
	updates := map[string]interface{}{}

	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)
		if name == "" {
			return PlaceDTO{}, ErrInvalidInput
		}
		updates["name"] = name
	}
	if input.Address != nil {
		updates["address"] = strings.TrimSpace(*input.Address)
	}
	if input.Lat != nil {
		if *input.Lat < -90 || *input.Lat > 90 {
			return PlaceDTO{}, ErrInvalidInput
		}
		updates["lat"] = *input.Lat
	}
	if input.Lng != nil {
		if *input.Lng < -180 || *input.Lng > 180 {
			return PlaceDTO{}, ErrInvalidInput
		}
		updates["lng"] = *input.Lng
	}
	if input.GooglePlaceID != nil {
		updates["google_place_id"] = strings.TrimSpace(*input.GooglePlaceID)
	}
	if len(updates) == 0 {
		return PlaceDTO{}, ErrInvalidInput
	}

	row, err := s.repo.Update(ctx, id, updates)
	if err != nil {
		return PlaceDTO{}, err
	}
	return toDTO(row), nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return ErrInvalidInput
	}
	return s.repo.Delete(ctx, id)
}

func toDTO(row model.Place) PlaceDTO {
	return PlaceDTO{
		ID:            row.ID,
		Name:          row.Name,
		Address:       row.Address,
		Lat:           row.Lat,
		Lng:           row.Lng,
		GooglePlaceID: row.GooglePlaceID,
		CreatedBy:     row.CreatedBy,
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     row.UpdatedAt,
	}
}
