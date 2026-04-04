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

type Service interface {
	Create(ctx context.Context, input CreateInput) (TripDTO, error)
	GetByID(ctx context.Context, id string) (TripDTO, error)
	ListByUser(ctx context.Context, userID string) ([]TripDTO, error)
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
