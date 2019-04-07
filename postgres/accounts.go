package postgres

import (
	"coins_wallet/account"
	"coins_wallet/config"
	"context"
	"github.com/go-pg/pg"
)

type AccountsRepository struct {
	db *pg.DB
}

func NewAccountsRepositry(settings *config.Postgres) (*AccountsRepository, error) {
	db, err := connect(settings)
	if err != nil {
		return nil, err
	}
	return &AccountsRepository{db: db}, nil
}

func (ar *AccountsRepository) GetAll(ctx context.Context, offset, limit int) ([]*account.Account, error) {
	var records []*account.Account
	_, err := ar.db.QueryContext(ctx,
		&records, "select id,balance,currency from accounts order by id offset ?0 limit ?1",
		offset, limit,
	)
	return records, err
}

func (ar *AccountsRepository) CountAll(ctx context.Context) (int, error) {
	var count int
	_, err := ar.db.QueryOneContext(ctx, pg.Scan(&count), "select count(*) from accounts")
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (ar *AccountsRepository) Get(ctx context.Context, accountId string) (*account.Account, error) {
	record := &account.Account{}
	_, err := ar.db.QueryOneContext(ctx,
		record, "select id,balance,currency from accounts where id=?0", accountId,
	)
	if err != nil {
		return nil, err
	}
	return record, nil
}
