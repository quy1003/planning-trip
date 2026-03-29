package phototag

import repo "planning-trip-be/internal/repository/phototag"

type Service interface{}

type service struct {
	repo repo.Repository
}

func NewService(repository repo.Repository) Service {
	return &service{repo: repository}
}
