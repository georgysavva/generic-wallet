package wallet

import (
	"context"
	"github.com/georgysavva/generic_wallet/account"
	"github.com/georgysavva/generic_wallet/payment"
	"github.com/go-kit/kit/log"
)

type loggingService struct {
	logger log.Logger
	Service
}

// NewLoggingService returns a new instance of a logging Service.
func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s *loggingService) SendPayment(ctx context.Context, fromAccountId, toAccountId string, amount float64) error {
	s.logger.Log(
		"method", "send_payment",
		"from_account", fromAccountId,
		"to_account", toAccountId,
		"amount", amount,
	)
	return s.Service.SendPayment(ctx, fromAccountId, toAccountId, amount)
}

func (s *loggingService) GetAllPayments(ctx context.Context, offset, limit *int) ([]*payment.Payment, int, error) {
	s.logger.Log(
		"method", "get_all_payments",
		"offset", offset,
		"limit", limit,
	)
	return s.Service.GetAllPayments(ctx, offset, limit)
}

func (s *loggingService) GetAllAccounts(ctx context.Context, offset, limit *int) ([]*account.Account, int, error) {
	s.logger.Log(
		"method", "get_all_accounts",
		"offset", offset,
		"limit", limit,
	)
	return s.Service.GetAllAccounts(ctx, offset, limit)
}
