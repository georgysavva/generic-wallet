package wallet

import (
	"context"
	"generic_wallet/account"
	"generic_wallet/inmem_repository"
	"generic_wallet/payment"
	"github.com/stretchr/testify/assert"
	"testing"
)

func instantiateServiceForTests() *service {
	accounts := []*account.Account{
		{Id: "alice", Balance: 100.0, Currency: "USD"},
		{Id: "bob", Balance: 100.0, Currency: "USD"},
		{Id: "mark", Balance: 100.0, Currency: "USD"},
		{Id: "john", Balance: 100.0, Currency: "USD"},
		{Id: "kate_in_europe", Balance: 100.0, Currency: "EUR"},
	}
	accountsRepo, PaymentsRepo := inmem_repository.InstantiateRepositories(accounts, nil)
	return NewService(PaymentsRepo, accountsRepo)
}

func TestSendPayment_Single(t *testing.T) {
	s := instantiateServiceForTests()
	ctx := context.Background()
	var err error
	offset, limit := 0, 0

	err = s.SendPayment(ctx, "alice", "bob", 20.0)
	assert.Equal(t, err, nil)
	fromAccount, _ := s.accounts.Get(ctx, "alice")
	toAccount, _ := s.accounts.Get(ctx, "bob")
	assert.Equal(t, fromAccount.Balance, 80.0)
	assert.Equal(t, toAccount.Balance, 120.0)
	paymentsList, totalPayments, _ := s.GetAllPayments(ctx, offset, limit)
	expectedPayments := []*payment.Payment{
		{AccountId: "alice", ToAccountId: "bob", Amount: 20.0, Direction: payment.OutgoingDirection},
		{AccountId: "bob", FromAccountId: "alice", Amount: 20.0, Direction: payment.IncomingDirection},
	}
	assert.Equal(t, expectedPayments, paymentsList)
	assert.Equal(t, 2, totalPayments)
}

func TestSendPayment_Bunch(t *testing.T) {
	s := instantiateServiceForTests()
	ctx := context.Background()
	var err error
	offset, limit := 0, 0

	err = s.SendPayment(ctx, "alice", "bob", 20.0)
	assert.Equal(t, err, nil)
	err = s.SendPayment(ctx, "alice", "john", 30.0)
	assert.Equal(t, err, nil)
	err = s.SendPayment(ctx, "mark", "bob", 40.0)
	assert.Equal(t, err, nil)
	aliceAccount, _ := s.accounts.Get(ctx, "alice")
	bobAccount, _ := s.accounts.Get(ctx, "bob")
	johnAccount, _ := s.accounts.Get(ctx, "john")
	markAccount, _ := s.accounts.Get(ctx, "mark")
	assert.Equal(t, aliceAccount.Balance, 50.0)
	assert.Equal(t, bobAccount.Balance, 160.0)
	assert.Equal(t, johnAccount.Balance, 130.0)
	assert.Equal(t, markAccount.Balance, 60.0)
	paymentsList, totalPayments, _ := s.GetAllPayments(ctx, offset, limit)
	expectedPayments := []*payment.Payment{
		{AccountId: "alice", ToAccountId: "bob", Amount: 20.0, Direction: payment.OutgoingDirection},
		{AccountId: "bob", FromAccountId: "alice", Amount: 20.0, Direction: payment.IncomingDirection},
		{AccountId: "alice", ToAccountId: "john", Amount: 30.0, Direction: payment.OutgoingDirection},
		{AccountId: "john", FromAccountId: "alice", Amount: 30.0, Direction: payment.IncomingDirection},
		{AccountId: "mark", ToAccountId: "bob", Amount: 40.0, Direction: payment.OutgoingDirection},
		{AccountId: "bob", FromAccountId: "mark", Amount: 40.0, Direction: payment.IncomingDirection},
	}
	assert.Equal(t, expectedPayments, paymentsList)
	assert.Equal(t, 6, totalPayments)
}

func TestSendPayment_AccountDoesNotExist(t *testing.T) {
	s := instantiateServiceForTests()
	ctx := context.Background()
	var err error
	offset, limit := 0, 0

	err = s.SendPayment(ctx, "unknown_from_account", "bob", 20.0)
	assert.Equal(t, err, FromAccountNotFound)
	toAccount, _ := s.accounts.Get(ctx, "bob")
	assert.Equal(t, toAccount.Balance, 100.0)
	paymentsList, totalPayments, _ := s.GetAllPayments(ctx, offset, limit)
	assert.Equal(t, 0, len(paymentsList))
	assert.Equal(t, 0, totalPayments)

	err = s.SendPayment(ctx, "alice", "unknown_to_account", 20.0)
	assert.Equal(t, err, ToAccountNotFound)
	fromAccount, _ := s.accounts.Get(ctx, "alice")
	assert.Equal(t, fromAccount.Balance, 100.0)
	paymentsList, totalPayments, _ = s.GetAllPayments(ctx, offset, limit)
	assert.Equal(t, 0, len(paymentsList))
	assert.Equal(t, 0, totalPayments)
}

func TestSendPayment_DifferentCurrencies(t *testing.T) {
	s := instantiateServiceForTests()
	ctx := context.Background()
	var err error
	offset, limit := 0, 0

	err = s.SendPayment(ctx, "alice", "kate_in_europe", 20.0)
	_, ok := err.(*DifferentCurrenciesError)
	assert.Equal(t, true, ok, "DifferentCurrenciesError type assertion")
	fromAccount, _ := s.accounts.Get(ctx, "alice")
	toAccount, _ := s.accounts.Get(ctx, "kate_in_europe")
	assert.Equal(t, fromAccount.Balance, 100.0)
	assert.Equal(t, toAccount.Balance, 100.0)
	paymentsList, totalPayments, _ := s.GetAllPayments(ctx, offset, limit)
	assert.Equal(t, 0, len(paymentsList))
	assert.Equal(t, 0, totalPayments)
}

func TestSendPayment_LowBalance(t *testing.T) {
	s := instantiateServiceForTests()
	ctx := context.Background()
	var err error
	offset, limit := 0, 0

	err = s.SendPayment(ctx, "alice", "bob", 200.0)
	assert.Equal(t, err, payment.LowBalanceErr)
	fromAccount, _ := s.accounts.Get(ctx, "alice")
	toAccount, _ := s.accounts.Get(ctx, "bob")
	assert.Equal(t, fromAccount.Balance, 100.0)
	assert.Equal(t, toAccount.Balance, 100.0)
	paymentsList, totalPayments, _ := s.GetAllPayments(ctx, offset, limit)
	assert.Equal(t, 0, len(paymentsList))
	assert.Equal(t, 0, totalPayments)
}
