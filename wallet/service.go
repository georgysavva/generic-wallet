package wallet

import (
	"coins_wallet/account"
	"coins_wallet/payment"
	"context"
	"fmt"
	"github.com/pkg/errors"
)

const (
	defaultPaginationLimit = 50
)

type Service interface {
	SendPayment(ctx context.Context, fromAccountId, toAccountId string, amount float64) error
	GetAllPayments(ctx context.Context, offset, limit int) ([]*payment.Payment, int, error)
	GetAllAccounts(ctx context.Context, offset, limit int) ([]*account.Account, int, error)
}

type service struct {
	payments payment.Repository
	accounts account.Repository
}

func NewService(payments payment.Repository, accounts account.Repository) *service {
	return &service{payments: payments, accounts: accounts}
}

func (s *service) SendPayment(ctx context.Context, fromAccountId, toAccountId string, amount float64) error {
	// Assumption: Account can't be deleted.
	// Assumption: Account currency can't be changed.
	if fromAccountId == toAccountId {
		return PaymentToSameAccountErr
	}
	if amount <= 0 {
		return NotPositivePaymentAmount
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
	newPayment := &payment.PaymentRequest{fromAccountId, toAccountId, amount}
	err = s.payments.Save(ctx, newPayment)
	return err
}

func (s *service) GetAllPayments(ctx context.Context, offset, limit int) ([]*payment.Payment, int, error) {
	if limit == 0 {
		limit = defaultPaginationLimit
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

func (s *service) GetAllAccounts(ctx context.Context, offset, limit int) ([]*account.Account, int, error) {
	if limit == 0 {
		limit = defaultPaginationLimit
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

//var LowBalanceErr = errors.New("source account doesn't have enough money")
var FromAccountNotFound = errors.New("source account not found")
var ToAccountNotFound = errors.New("destination account not found")

var PaymentToSameAccountErr = errors.New("source account and destination account are the same")
var NotPositivePaymentAmount = errors.New("payment amount must be greater than 0")

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
