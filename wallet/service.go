package wallet

import (
	"context"
	"fmt"
	"github.com/georgysavva/generic_wallet/account"
	"github.com/georgysavva/generic_wallet/payment"
	"github.com/pkg/errors"
)

const (
	defaultPaginationLimit = 50
)

type Service interface {
	SendPayment(ctx context.Context, fromAccountId, toAccountId string, amount float64) error
	GetAllPayments(ctx context.Context, offset, limit *int) ([]*payment.Payment, int, error)
	GetAllAccounts(ctx context.Context, offset, limit *int) ([]*account.Account, int, error)
}

type service struct {
	payments payment.Repository
	accounts account.Repository
}

func NewService(payments payment.Repository, accounts account.Repository) Service {
	return &service{payments: payments, accounts: accounts}
}

func (s *service) SendPayment(ctx context.Context, fromAccountId, toAccountId string, amount float64) error {
	// Assumption: Account can't be deleted.
	// Assumption: Account currency can't be changed.
	if fromAccountId == toAccountId {
		return &IncorrectInputData{"source account and destination account are the same"}
	}
	if amount <= 0 {
		return &IncorrectInputData{"payment amount must be greater than 0"}
	}
	fromAccount, err := s.accounts.Get(ctx, fromAccountId)
	if err != nil {
		return err
	}
	if fromAccount == nil {
		return FromAccountNotFound
	}
	toAccount, err := s.accounts.Get(ctx, toAccountId)
	if err != nil {
		return err
	}
	if toAccount == nil {
		return ToAccountNotFound
	}
	if fromAccount.Currency != toAccount.Currency {
		return &DifferentCurrenciesError{fromAccount.Currency, toAccount.Currency}
	}
	err = s.payments.Save(ctx, fromAccountId, toAccountId, amount)
	return err
}

func (s *service) GetAllPayments(ctx context.Context, offset, limit *int) ([]*payment.Payment, int, error) {
	offset, limit, err := preparePagination(offset, limit)
	if err != nil {
		return nil, 0, err
	}
	paymentRecords, err := s.payments.GetAll(ctx, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	paymentsTotal, err := s.payments.CountAll(ctx)
	if err != nil {
		return nil, 0, err
	}
	return paymentRecords, paymentsTotal, nil
}

func (s *service) GetAllAccounts(ctx context.Context, offset, limit *int) ([]*account.Account, int, error) {
	offset, limit, err := preparePagination(offset, limit)
	if err != nil {
		return nil, 0, err
	}
	accountRecords, err := s.accounts.GetAll(ctx, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	accountsTotal, err := s.accounts.CountAll(ctx)
	if err != nil {
		return nil, 0, err
	}
	return accountRecords, accountsTotal, nil
}
func preparePagination(offset, limit *int) (*int, *int, error) {
	if offset != nil && *offset < 0 {
		return nil, nil, &IncorrectInputData{"'offset'pagination parameter must be >= 0"}
	}
	if limit != nil {
		if *limit < 0 {
			return nil, nil, &IncorrectInputData{"'limit'pagination parameter must be >= 0"}
		}
	} else {
		limit = new(int)
		*limit = defaultPaginationLimit
	}
	return offset, limit, nil
}

var FromAccountNotFound = errors.New("source account not found")
var ToAccountNotFound = errors.New("destination account not found")

type IncorrectInputData struct {
	Details string
}

func (e *IncorrectInputData) Error() string {
	return e.Details
}

type DifferentCurrenciesError struct {
	FromAccountCurrency string
	ToAccountCurrency   string
}

func (e *DifferentCurrenciesError) Error() string {
	return fmt.Sprintf(
		"source account has %s currency, but destination account has %s currency",
		e.FromAccountCurrency, e.ToAccountCurrency,
	)
}
