package service

import (
	"context"
	"wb-l0/internal/models"
	"wb-l0/internal/repository"
)

type OrderService struct {
	repo repository.Order
}

func NewOrderService(r repository.Order) *OrderService {
	return &OrderService{
		repo: r,
	}
}

func (o *OrderService) FindByUid(ctx context.Context, uid string) (models.Order, error) {
	return o.repo.FindByUid(ctx, uid)
}

func (o *OrderService) FindAll(ctx context.Context) ([]models.Order, error) {
	return o.repo.FindAll(ctx)
}
