package health

import (
	"context"
	"time"
)

type Response struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

type Service interface {
	Check(ctx context.Context) (Response, error)
}

type service struct{}

func NewService() Service {
	return &service{}
}

func (s *service) Check(ctx context.Context) (Response, error) {
	return Response{
		Status: "dsadsa",
	}, nil
}
