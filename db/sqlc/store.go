package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Store interface {
	execTX(context.Context, func(*Queries) error) error
	GetQueries() *Queries
	TransferMoney(ctx context.Context, accIdFrom, accIdTo, amount int64) error
}

type store struct {
	db      *pgx.Conn
	Queries *Queries
}

func NewStore(db *pgx.Conn) Store {
	return &store{db: db, Queries: New(db)}
}

func (s *store) execTX(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.Begin(ctx)

	if err != nil {
		return err
	}

	if err := fn(s.Queries.WithTx(tx)); err != nil {
		if e := tx.Rollback(ctx); e != nil {
			return fmt.Errorf("rallback error: %w. transaction error: %w", e, err)
		}
		return err
	}

	return tx.Commit(ctx)
}

func (s *store) GetQueries() *Queries {
	return s.Queries
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
		accFrom, accTo := accs[0], accs[1]
		if accs[0].ID == accIdTo {
			accFrom, accTo = accTo, accFrom
		}
		if err != nil {
			return err
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
