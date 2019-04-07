package account

import (
	"context"
)

type Account struct {
	Id       string  `json:"id"`
	Balance  float64 `json:"balance"`
	Currency string  `json:"currency"`
}

type Repository interface {
	GetAll(ctx context.Context, offset, limit int) ([]*Account, error)
	CountAll(ctx context.Context) (int, error)
	Get(ctx context.Context, accountId string) (*Account, error)
}
