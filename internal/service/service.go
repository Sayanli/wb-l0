package service

import (
	"context"
	"wb-l0/internal/models"
	"wb-l0/internal/repository"
)

type Order interface {
	FindByUid(ctx context.Context, uid string) (models.Order, error)
	FindAll(ctx context.Context) ([]models.Order, error)
}

type Service struct {
	Order
}

func NewService(r *repository.Repository) *Service {
	return &Service{
		Order: NewOrderService(r.Order),
	}
}
