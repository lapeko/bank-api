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
	transfer  *Transfer
	entryFrom *Entry
	entryTo   *Entry
}

func (s *Store) TransferTX(ctx context.Context, params CreateTransferParams) (*TransferTxResult, error) {
	var transfer Transfer
	var entryFrom Entry
	var entryTo Entry
	var err error

	err = s.execTx(ctx, func(q *Queries) error {
		account, err := q.GetAccount(ctx, params.AccountFrom)
		if err != nil {
			return err
		}
		if account.Balance < params.Amount {
			return errors.New("insufficient funds")
		}
		transfer, err = q.CreateTransfer(ctx, params)
		if err != nil {
			return err
		}
		entryFrom, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: params.AccountFrom,
			Amount:    -params.Amount,
		})
		if err != nil {
			return err
		}
		entryTo, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: params.AccountTo,
			Amount:    params.Amount,
		})
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &TransferTxResult{
		transfer:  &transfer,
		entryFrom: &entryFrom,
		entryTo:   &entryTo,
	}, err
}
