package repository

import (
	"context"
	"database/sql"
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
	Transfer    Transfer `json:"transfer"`
	EntryFrom   Entry    `json:"entryFrom"`
	EntryTo     Entry    `json:"entryTo"`
	AccountFrom Account  `json:"accountFrom"`
	AccountTo   Account  `json:"accountTo"`
}

func (s *SqlStore) TransferTX(ctx context.Context, params CreateTransferParams) (*TransferTxResult, error) {
	var result TransferTxResult
	var err error

	err = s.execTx(ctx, func(q *Queries) error {
		result.Transfer, err = q.CreateTransfer(ctx, params)
		if err != nil {
			return err
		}

		result.EntryFrom, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: params.AccountFrom,
			Amount:    -params.Amount,
		})
		if err != nil {
			return err
		}

		result.EntryTo, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: params.AccountTo,
			Amount:    params.Amount,
		})

		if params.AccountFrom < params.AccountTo {
			result.AccountFrom, result.AccountTo, err = updateBalances(
				ctx,
				q,
				params.AccountFrom,
				-params.Amount,
				params.AccountTo,
				params.Amount,
			)
		} else {
			result.AccountTo, result.AccountFrom, err = updateBalances(
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
