package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type DBConn interface {
	DBTX
	Begin(ctx context.Context) (pgx.Tx, error)
}

type Store interface {
	Querier
	TransferMoney(context.Context, int64, int64, int64) error
	TransferExternalMoney(context.Context, int64, int64) (Account, error)
}

type store struct {
	db DBConn
	*Queries
}

type transferClientErrorMessage string

var (
	NotPositiveAmount transferClientErrorMessage = "amount should be positive"
	SameAccount       transferClientErrorMessage = "transfer accounts should not be same"
	NotSameCurrency   transferClientErrorMessage = "currency must be same for money transfer"
	InsufficientFunds transferClientErrorMessage = "insufficient funds"
)

type TransferClientError struct {
	message transferClientErrorMessage
}

func (e *TransferClientError) Error() string {
	return fmt.Sprintf("transfer error: %s", e.message)
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
		return &TransferClientError{NotPositiveAmount}
	}
	if accIdFrom == accIdTo {
		return &TransferClientError{SameAccount}
	}

	return s.execTX(ctx, func(q *Queries) error {
		accs, err := q.GetTwoAccountsByIdForUpdate(ctx, GetTwoAccountsByIdForUpdateParams{accIdFrom, accIdTo})
		if err != nil {
			return err
		}
		accFrom, accTo := accs[0], accs[1]
		if accs[0].ID == accIdTo {
			accFrom, accTo = accTo, accFrom
		}
		if accFrom.Currency != accTo.Currency {
			return &TransferClientError{NotSameCurrency}
		}
		if accFrom.Balance < amount {
			return &TransferClientError{InsufficientFunds}
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
		if _, err := q.OffsetAccountBalance(ctx, OffsetAccountBalanceParams{accIdFrom, -amount}); err != nil {
			return err
		}
		if _, err := q.OffsetAccountBalance(ctx, OffsetAccountBalanceParams{accIdTo, amount}); err != nil {
			return err
		}
		return nil
	})
}

func (s *store) TransferExternalMoney(ctx context.Context, accountId, amount int64) (updatedAcc Account, e error) {
	e = s.execTX(ctx, func(q *Queries) error {
		acc, err := q.GetAccountByIdForUpdate(ctx, accountId)
		if err != nil {
			return err
		}
		if amount < 0 && acc.Balance+amount < 0 {
			return &TransferClientError{InsufficientFunds}
		}
		if _, err = q.CreateEntry(ctx, CreateEntryParams{AccountID: accountId, Amount: amount}); err != nil {
			return err
		}
		updatedAcc, err = q.OffsetAccountBalance(ctx, OffsetAccountBalanceParams{ID: accountId, Delta: amount})
		return err
	})
	return
}
