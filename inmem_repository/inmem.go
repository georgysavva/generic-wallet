package inmem_repository

import (
	"context"
	"github.com/georgysavva/generic_wallet/account"
	"github.com/georgysavva/generic_wallet/payment"
	"github.com/pkg/errors"
	"sort"
)

func InstantiateRepositories(accounts []*account.Account, payments []*payment.Payment) (*AccountsRepository, *PaymentsRepository) {
	accountsRepo := &AccountsRepository{accounts: map[string]*account.Account{}}
	for _, a := range accounts {
		accountsRepo.accounts[a.Id] = a
	}
	paymentsRepo := &PaymentsRepository{accountsRepo: accountsRepo}
	for _, p := range payments {
		paymentsRepo.payments = append(paymentsRepo.payments, p)
	}
	return accountsRepo, paymentsRepo
}

type AccountsRepository struct {
	accounts map[string]*account.Account
}

func (ar *AccountsRepository) GetAll(ctx context.Context, offset, limit *int) ([]*account.Account, error) {
	var accountsList []*account.Account
	for _, accountRecord := range ar.accounts {
		accountsList = append(accountsList, accountRecord)
	}
	sort.Slice(accountsList, func(i, j int) bool {
		return accountsList[i].Id < accountsList[j].Id
	})
	return accountsList, nil
}

func (ar *AccountsRepository) CountAll(ctx context.Context) (int, error) {
	return len(ar.accounts), nil
}

func (ar *AccountsRepository) Get(ctx context.Context, accountId string) (*account.Account, error) {
	accountRecords := ar.accounts[accountId]
	return accountRecords, nil
}

type PaymentsRepository struct {
	payments     []*payment.Payment
	accountsRepo *AccountsRepository
}

func (pr *PaymentsRepository) GetAll(ctx context.Context, offset, limit *int) ([]*payment.Payment, error) {
	return pr.payments, nil
}

func (pr *PaymentsRepository) CountAll(ctx context.Context) (int, error) {
	return len(pr.payments), nil
}

func (pr *PaymentsRepository) Save(ctx context.Context, fromAccountId, toAccountId string, amount float64) error {
	fromAccount := pr.accountsRepo.accounts[fromAccountId]
	if fromAccount == nil {
		return errors.New("source account not found")
	}
	toAccount := pr.accountsRepo.accounts[toAccountId]
	if toAccount == nil {
		return errors.New("destination account not found")
	}
	if fromAccount.Balance-amount < 0 {
		return payment.LowBalanceErr
	}
	pr.payments = append(pr.payments, &payment.Payment{
		AccountId:   fromAccount.Id,
		ToAccountId: toAccount.Id,
		Amount:      amount,
		Direction:   payment.OutgoingDirection,
	})
	pr.payments = append(pr.payments, &payment.Payment{
		AccountId:     toAccount.Id,
		FromAccountId: fromAccount.Id,
		Amount:        amount,
		Direction:     payment.IncomingDirection,
	})
	fromAccount.Balance -= amount
	toAccount.Balance += amount
	return nil
}
