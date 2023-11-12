package repository

import (
	"context"
	"wb-l0/internal/cache"
	"wb-l0/internal/models"
	"wb-l0/pkg/database"

	"github.com/nats-io/stan.go"
)

type Order interface {
	FindByUid(ctx context.Context, uid string) (models.Order, error)
	FindAll(ctx context.Context) ([]models.Order, error)
	CreateOrder(msg *stan.Msg)
}

type Repository struct {
	Order
}

func NewRepository(db database.PGXQuerier, cache *cache.Cache) *Repository {
	return &Repository{
		Order: NewOrderRepository(db, cache),
	}
}
