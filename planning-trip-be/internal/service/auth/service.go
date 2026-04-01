package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"planning-trip-be/internal/model"
	urepo "planning-trip-be/internal/repository/user"
)

var (
	ErrInvalidInput       = errors.New("invalid input")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailAlreadyUsed   = errors.New("email already in use")
)

const tokenTTL = 24 * time.Hour

type RegisterInput struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthUser struct {
	ID        string    `json:"id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	AvatarURL string    `json:"avatar_url,omitempty"`
	Bio       string    `json:"bio,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AuthResult struct {
	AccessToken string   `json:"access_token"`
	TokenType   string   `json:"token_type"`
	ExpiresIn   int64    `json:"expires_in"`
	User        AuthUser `json:"user"`
}

type Service interface {
	Register(ctx context.Context, input RegisterInput) (AuthResult, error)
	Login(ctx context.Context, input LoginInput) (AuthResult, error)
}

type service struct {
	userRepo urepo.Repository
	secret   string
}

func NewService(userRepo urepo.Repository, secret string) Service {
	return &service{
		userRepo: userRepo,
		secret:   strings.TrimSpace(secret),
	}
}

func (s *service) Register(ctx context.Context, input RegisterInput) (AuthResult, error) {
	fullName := strings.TrimSpace(input.FullName)
	email := strings.ToLower(strings.TrimSpace(input.Email))
	password := strings.TrimSpace(input.Password)

	if fullName == "" || email == "" || len(password) < 6 {
		return AuthResult{}, ErrInvalidInput
	}

	if _, err := s.userRepo.GetByEmail(ctx, email); err == nil {
		return AuthResult{}, ErrEmailAlreadyUsed
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return AuthResult{}, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return AuthResult{}, err
	}

	user := model.User{
		FullName:     fullName,
		Email:        email,
		PasswordHash: string(hashedPassword),
	}
	if err := s.userRepo.Create(ctx, &user); err != nil {
		return AuthResult{}, err
	}

	return s.buildAuthResult(user)
}

func (s *service) Login(ctx context.Context, input LoginInput) (AuthResult, error) {
	email := strings.ToLower(strings.TrimSpace(input.Email))
	password := strings.TrimSpace(input.Password)

	if email == "" || password == "" {
		return AuthResult{}, ErrInvalidInput
	}

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return AuthResult{}, ErrInvalidCredentials
		}
		return AuthResult{}, err
	}

	if compareErr := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); compareErr != nil {
		return AuthResult{}, ErrInvalidCredentials
	}

	return s.buildAuthResult(user)
}

func (s *service) buildAuthResult(user model.User) (AuthResult, error) {
	token, err := s.signAccessToken(user.ID, user.Email)
	if err != nil {
		return AuthResult{}, err
	}

	return AuthResult{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int64(tokenTTL.Seconds()),
		User: AuthUser{
			ID:        user.ID,
			FullName:  user.FullName,
			Email:     user.Email,
			AvatarURL: user.AvatarURL,
			Bio:       user.Bio,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}, nil
}

func (s *service) signAccessToken(userID, email string) (string, error) {
	if s.secret == "" {
		return "", fmt.Errorf("auth secret is required")
	}

	expiresAt := time.Now().UTC().Add(tokenTTL).Unix()
	payload := userID + "|" + email + "|" + strconv.FormatInt(expiresAt, 10)

	mac := hmac.New(sha256.New, []byte(s.secret))
	if _, err := mac.Write([]byte(payload)); err != nil {
		return "", err
	}

	signature := hex.EncodeToString(mac.Sum(nil))
	raw := payload + "|" + signature
	return base64.RawURLEncoding.EncodeToString([]byte(raw)), nil
}
