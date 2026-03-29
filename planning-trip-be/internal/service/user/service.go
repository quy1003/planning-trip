package user

import (
	"context"
	"errors"
	"strings"
	"time"

	"planning-trip-be/internal/model"
	urepo "planning-trip-be/internal/repository/user"

	"gorm.io/gorm"
)

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrNotFound     = gorm.ErrRecordNotFound
)

type UserDTO struct {
	ID        string    `json:"id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	AvatarURL string    `json:"avatar_url,omitempty"`
	Bio       string    `json:"bio,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateInput struct {
	FullName     string `json:"full_name"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	AvatarURL    string `json:"avatar_url"`
	Bio          string `json:"bio"`
}

type UpdateInput struct {
	FullName     *string `json:"full_name"`
	Email        *string `json:"email"`
	PasswordHash *string `json:"password_hash"`
	AvatarURL    *string `json:"avatar_url"`
	Bio          *string `json:"bio"`
}

type Service interface {
	List(ctx context.Context) ([]UserDTO, error)
	GetByID(ctx context.Context, id string) (UserDTO, error)
	Create(ctx context.Context, input CreateInput) (UserDTO, error)
	Update(ctx context.Context, id string, input UpdateInput) (UserDTO, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo urepo.Repository
}

func NewService(repo urepo.Repository) Service {
	return &service{repo: repo}
}

func (s *service) List(ctx context.Context) ([]UserDTO, error) {
	users, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]UserDTO, 0, len(users))
	for _, u := range users {
		result = append(result, toDTO(u))
	}
	return result, nil
}

func (s *service) GetByID(ctx context.Context, id string) (UserDTO, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return UserDTO{}, err
	}
	return toDTO(user), nil
}

func (s *service) Create(ctx context.Context, input CreateInput) (UserDTO, error) {
	if strings.TrimSpace(input.FullName) == "" || strings.TrimSpace(input.Email) == "" || strings.TrimSpace(input.PasswordHash) == "" {
		return UserDTO{}, ErrInvalidInput
	}

	user := model.User{
		FullName:     strings.TrimSpace(input.FullName),
		Email:        strings.ToLower(strings.TrimSpace(input.Email)),
		PasswordHash: strings.TrimSpace(input.PasswordHash),
		AvatarURL:    strings.TrimSpace(input.AvatarURL),
		Bio:          strings.TrimSpace(input.Bio),
	}

	if err := s.repo.Create(ctx, &user); err != nil {
		return UserDTO{}, err
	}
	return toDTO(user), nil
}

func (s *service) Update(ctx context.Context, id string, input UpdateInput) (UserDTO, error) {
	updates := map[string]interface{}{}

	if input.FullName != nil {
		trimmed := strings.TrimSpace(*input.FullName)
		if trimmed == "" {
			return UserDTO{}, ErrInvalidInput
		}
		updates["full_name"] = trimmed
	}
	if input.Email != nil {
		trimmed := strings.ToLower(strings.TrimSpace(*input.Email))
		if trimmed == "" {
			return UserDTO{}, ErrInvalidInput
		}
		updates["email"] = trimmed
	}
	if input.PasswordHash != nil {
		trimmed := strings.TrimSpace(*input.PasswordHash)
		if trimmed == "" {
			return UserDTO{}, ErrInvalidInput
		}
		updates["password_hash"] = trimmed
	}
	if input.AvatarURL != nil {
		updates["avatar_url"] = strings.TrimSpace(*input.AvatarURL)
	}
	if input.Bio != nil {
		updates["bio"] = strings.TrimSpace(*input.Bio)
	}

	if len(updates) == 0 {
		return UserDTO{}, ErrInvalidInput
	}

	updates["updated_at"] = time.Now().UTC()
	user, err := s.repo.Update(ctx, id, updates)
	if err != nil {
		return UserDTO{}, err
	}
	return toDTO(user), nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func toDTO(user model.User) UserDTO {
	return UserDTO{
		ID:        user.ID,
		FullName:  user.FullName,
		Email:     user.Email,
		AvatarURL: user.AvatarURL,
		Bio:       user.Bio,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
