package transfer

import (
	"context"

	db "github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/db/sqlc"
)

type transferService struct {
	store db.Store
}

func (s *transferService) transferMoney(ctx context.Context, args transferRequest) error {
	return s.store.TransferMoney(ctx, args.AccountFrom, args.AccountTo, args.Amount)
}
