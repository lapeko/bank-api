package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (s *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("err: %v, rallbackError: %v", err, rbErr)
		}
		return err
	}

	err = tx.Commit()

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("err: %v, rallbackError: %v", err, rbErr)
		}
		return err
	}

	return nil
}

type TransferTxResult struct {
	transfer    Transfer
	entryFrom   Entry
	entryTo     Entry
	accountFrom Account
	accountTo   Account
}

var transferTxCtxKey = struct{}{}

func (s *Store) TransferTX(ctx context.Context, params CreateTransferParams) (*TransferTxResult, error) {
	var result TransferTxResult
	var err error

	err = s.execTx(ctx, func(q *Queries) error {
		account, err := q.GetAccount(ctx, params.AccountFrom)

		ctxKey := ctx.Value(transferTxCtxKey)

		if err != nil {
			return err
		}
		if account.Balance < params.Amount {
			return errors.New("insufficient funds")
		}
		fmt.Printf("%s\t%s\n", ctxKey, "CreateTransfer")
		result.transfer, err = q.CreateTransfer(ctx, params)
		if err != nil {
			return err
		}
		fmt.Printf("%s\t%s\n", ctxKey, "CreateEntry 1")
		result.entryFrom, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: params.AccountFrom,
			Amount:    -params.Amount,
		})
		if err != nil {
			return err
		}
		fmt.Printf("%s\t%s\n", ctxKey, "CreateEntry 2")
		result.entryTo, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: params.AccountTo,
			Amount:    params.Amount,
		})
		// implement in a wrong way
		fmt.Printf("%s\t%s\n", ctxKey, "GetAccountForUpdate 1")
		accFrom, err := q.GetAccountForUpdate(ctx, params.AccountFrom)
		if err != nil {
			return err
		}
		fmt.Printf("%s\t%s\n", ctxKey, "UpdateAccount 1")
		result.accountFrom, err = q.UpdateAccount(ctx, UpdateAccountParams{ID: accFrom.ID, Balance: accFrom.Balance - params.Amount})
		if err != nil {
			return err
		}
		fmt.Printf("%s\t%s\n", ctxKey, "GetAccountForUpdate 2")
		accTo, err := q.GetAccountForUpdate(ctx, params.AccountTo)
		if err != nil {
			return err
		}
		fmt.Printf("%s\t%s\n", ctxKey, "UpdateAccount 2")
		result.accountTo, err = q.UpdateAccount(ctx, UpdateAccountParams{ID: accTo.ID, Balance: accTo.Balance + params.Amount})
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &result, err
}
