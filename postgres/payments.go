package postgres

import (
	"coins_wallet/config"
	"coins_wallet/payment"
	"context"
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

func (pr *PaymentsRepository) GetAll(ctx context.Context, offset, limit int) ([]*payment.Payment, error) {
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

func (pr *PaymentsRepository) Save(ctx context.Context, paymentReq *payment.PaymentRequest) error {
	err := pr.db.RunInTransaction(func(tx *pg.Tx) error {
		var fromAccountBalance float64
		_, err := tx.QueryOneContext(ctx,
			pg.Scan(&fromAccountBalance),
			"select balance from accounts where id=?0 for update",
			paymentReq.FromAccountId,
		)
		if err != nil {
			return err
		}
		if fromAccountBalance-paymentReq.Amount < 0 {
			return payment.LowBalanceErr
		}

		// Create an outgoing payment.
		_, err = tx.ExecOneContext(ctx,
			"insert into payments (account_id,to_account_id,amount,direction) values (?0,?1,?2,?3)",
			paymentReq.FromAccountId, paymentReq.ToAccountId, paymentReq.Amount, payment.OutgoingDirection,
		)
		if err != nil {
			return err
		}

		// Create an incoming payment.
		_, err = tx.ExecOneContext(ctx,
			"insert into payments (account_id,from_account_id,amount,direction) values (?0,?1,?2,?3)",
			paymentReq.ToAccountId, paymentReq.FromAccountId, paymentReq.Amount, payment.IncomingDirection,
		)
		if err != nil {
			return err
		}

		// Decrease source account.
		_, err = tx.ExecOneContext(ctx,
			"update accounts set balance = balance - ?0 where id=?1",
			paymentReq.Amount, paymentReq.FromAccountId,
		)
		if err != nil {
			return err
		}

		// Increase destination account.
		_, err = tx.ExecOneContext(ctx,
			"update accounts set balance = balance + ?0 where id=?1",
			paymentReq.Amount, paymentReq.ToAccountId,
		)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}
