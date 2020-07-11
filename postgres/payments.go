package postgres

import (
	"context"
	"github.com/georgysavva/generic-wallet/config"
	"github.com/georgysavva/generic-wallet/payment"
	"github.com/go-pg/pg"
)

type PaymentsRepository struct {
	db *pg.DB
}

func NewPaymentsRepository(settings *config.Postgres) (*PaymentsRepository, error) {
	db, err := connect(settings)
	if err != nil {
		return nil, err
	}
	return &PaymentsRepository{db: db}, nil
}

func (pr *PaymentsRepository) GetAll(ctx context.Context, offset, limit *int) ([]*payment.Payment, error) {
	var records []*payment.Payment
	_, err := pr.db.QueryContext(ctx,
		&records,
		"select account_id,to_account_id,from_account_id,amount,direction "+
			"from payments order by id offset ?0 limit ?1",
		offset, limit,
	)
	return records, err
}

func (pr *PaymentsRepository) CountAll(ctx context.Context) (int, error) {
	var count int
	_, err := pr.db.QueryOneContext(ctx, pg.Scan(&count), "select count(*) from payments")
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (pr *PaymentsRepository) Save(ctx context.Context, fromAccountId, toAccountId string, amount float64) error {
	err := pr.db.RunInTransaction(func(tx *pg.Tx) error {
		var fromAccountBalance float64
		// We need to lock source account row
		// to prevent race condition on the balance field.
		_, err := tx.QueryOneContext(ctx,
			pg.Scan(&fromAccountBalance),
			"select balance from accounts where id=?0 for update",
			fromAccountId,
		)
		if err != nil {
			return err
		}
		if fromAccountBalance-amount < 0 {
			return payment.LowBalanceErr
		}

		// Create an outgoing payment.
		_, err = tx.ExecOneContext(ctx,
			"insert into payments (account_id,to_account_id,amount,direction) values (?0,?1,?2,?3)",
			fromAccountId, toAccountId, amount, payment.OutgoingDirection,
		)
		if err != nil {
			return err
		}

		// Create an incoming payment.
		_, err = tx.ExecOneContext(ctx,
			"insert into payments (account_id,from_account_id,amount,direction) values (?0,?1,?2,?3)",
			toAccountId, fromAccountId, amount, payment.IncomingDirection,
		)
		if err != nil {
			return err
		}

		// Decrease source account.
		_, err = tx.ExecOneContext(ctx,
			"update accounts set balance = balance - ?0 where id=?1",
			amount, fromAccountId,
		)
		if err != nil {
			return err
		}

		// Increase destination account.
		_, err = tx.ExecOneContext(ctx,
			"update accounts set balance = balance + ?0 where id=?1",
			amount, toAccountId,
		)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}
