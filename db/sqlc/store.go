package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Store interface {
	execTX(context.Context, func(*Queries) error) error
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
