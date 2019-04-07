package payment

import (
	"context"
	"github.com/pkg/errors"
)

const (
	OutgoingDirection = "outgoing"
	IncomingDirection = "incoming"
)

type Payment struct {
	AccountId     string  `json:"account"`
	ToAccountId   string  `json:"to_account,omitempty"`
	FromAccountId string  `json:"from_account,omitempty"`
	Amount        float64 `json:"amount"`
	Direction     string  `json:"direction"`
}

type PaymentRequest struct {
	FromAccountId string
	ToAccountId   string
	Amount        float64
}

var LowBalanceErr = errors.New("account doesn't have enough money to send the payment")

type Repository interface {
	GetAll(ctx context.Context, offset, limit int) ([]*Payment, error)
	CountAll(ctx context.Context) (int, error)
	Save(ctx context.Context, payment *PaymentRequest) error
}
