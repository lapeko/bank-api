package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type DBConn interface {
	DBTX
	Begin(ctx context.Context) (pgx.Tx, error)
}

type Store interface {
	Querier
	TransferMoney(ctx context.Context, accIdFrom, accIdTo, amount int64) error
}

type store struct {
	db DBConn
	*Queries
}

func NewStore(db DBConn) Store {
	return &store{db: db, Queries: New(db)}
}

func (s *store) execTX(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.Begin(ctx)

	if err != nil {
		return err
	}

	if err := fn(s.Queries.WithTx(tx)); err != nil {
		if e := tx.Rollback(ctx); e != nil {
			return errors.Join(err, e)
		}
		return err
	}

	return tx.Commit(ctx)
}

func (s *store) TransferMoney(ctx context.Context, accIdFrom, accIdTo, amount int64) error {
	if amount <= 0 {
		return errors.New("amount should be positive")
	}
	if accIdFrom == accIdTo {
		return errors.New("transfer accounts should not be same")
	}

	return s.execTX(ctx, func(q *Queries) error {
		accs, err := q.GetAccountsByIdForUpdate(ctx, GetAccountsByIdForUpdateParams{accIdFrom, accIdTo})
		if err != nil {
			return err
		}
		accFrom, accTo := accs[0], accs[1]
		if accs[0].ID == accIdTo {
			accFrom, accTo = accTo, accFrom
		}
		if accFrom.Currency != accTo.Currency {
			return errors.New("currency must be same for money transfer")
		}
		if accFrom.Balance < amount {
			return errors.New("insufficient funds")
		}
		if _, err := q.CreateEntry(ctx, CreateEntryParams{accIdFrom, -amount}); err != nil {
			return err
		}
		if _, err := q.CreateEntry(ctx, CreateEntryParams{accIdTo, amount}); err != nil {
			return err
		}
		if _, err := q.CreateTransfer(ctx, CreateTransferParams{accIdFrom, accIdTo, amount}); err != nil {
			return err
		}
		if _, err := q.OffsetBalance(ctx, OffsetBalanceParams{accIdFrom, -amount}); err != nil {
			return err
		}
		if _, err := q.OffsetBalance(ctx, OffsetBalanceParams{accIdTo, amount}); err != nil {
			return err
		}
		return nil
	})
}
