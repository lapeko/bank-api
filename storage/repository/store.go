package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/golang/mock/mockgen/model"
)

type Store interface {
	Querier
	TransferTX(ctx context.Context, params CreateTransferParams) (*TransferTxResult, error)
}

type SqlStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *SqlStore {
	return &SqlStore{
		db:      db,
		Queries: New(db),
	}
}

func (s *SqlStore) execTx(ctx context.Context, fn func(*Queries) error) error {
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

func (s *SqlStore) TransferTX(ctx context.Context, params CreateTransferParams) (*TransferTxResult, error) {
	var result TransferTxResult
	var err error

	err = s.execTx(ctx, func(q *Queries) error {
		// TODO add FOR UPDATE
		account, err := q.GetAccount(ctx, params.AccountFrom)
		if err != nil {
			return err
		}
		if account.Balance < params.Amount {
			return errors.New("insufficient funds")
		}

		result.transfer, err = q.CreateTransfer(ctx, params)
		if err != nil {
			return err
		}

		result.entryFrom, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: params.AccountFrom,
			Amount:    -params.Amount,
		})
		if err != nil {
			return err
		}

		result.entryTo, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: params.AccountTo,
			Amount:    params.Amount,
		})

		if params.AccountFrom < params.AccountTo {
			result.accountFrom, result.accountTo, err = updateBalances(
				ctx,
				q,
				params.AccountFrom,
				-params.Amount,
				params.AccountTo,
				params.Amount,
			)
		} else {
			result.accountTo, result.accountFrom, err = updateBalances(
				ctx,
				q,
				params.AccountTo,
				params.Amount,
				params.AccountFrom,
				-params.Amount,
			)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &result, err
}

func updateBalances(ctx context.Context, q *Queries, acc1Id int64, amount1 int64, acc2Id int64, amount2 int64) (acc1 Account, acc2 Account, err error) {
	acc1, err = q.UpdateAccountBalanceBy(ctx, UpdateAccountBalanceByParams{ID: acc1Id, Amount: amount1})
	if err != nil {
		return
	}
	acc2, err = q.UpdateAccountBalanceBy(ctx, UpdateAccountBalanceByParams{ID: acc2Id, Amount: amount2})
	return
}
