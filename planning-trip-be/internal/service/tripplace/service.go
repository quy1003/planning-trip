package tripplace

import repo "planning-trip-be/internal/repository/tripplace"

type Service interface{}

type service struct {
	repo repo.Repository
}

func NewService(repository repo.Repository) Service {
	return &service{repo: repository}
}
