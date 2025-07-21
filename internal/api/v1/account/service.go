package account

import (
	"context"

	db "github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/db/sqlc"
	"golang.org/x/sync/errgroup"
)

type accountService struct {
	store db.Store
}

func (a *accountService) createAccount(ctx context.Context, params *createAccountRequest) (db.Account, error) {
	reqParams := db.CreateAccountParams{UserID: params.UserID, Currency: params.Currency}
	return a.store.CreateAccount(ctx, reqParams)
}

func (a *accountService) listAccounts(ctx context.Context, params listAccountsRequest) (listAccountsResponse, error) {
	reqParams := db.ListAccountsParams{
		Limit:  params.Size,
		Offset: (params.Page - 1) * params.Size,
	}

	var rows []db.ListAccountsRow
	var totalCount int64

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		var err error
		rows, err = a.store.ListAccounts(ctx, reqParams)
		return err
	})
	g.Go(func() error {
		var err error
		totalCount, err = a.store.GetTotalAccountsCount(ctx)
		return err
	})

	if err := g.Wait(); err != nil {
		return listAccountsResponse{}, err
	}

	return listAccountsResponse{
		Accounts:   listAccountRowsToAccountWithUserInfo(rows),
		TotalCount: totalCount,
	}, nil
}

func (a *accountService) getAccountById(ctx context.Context, id int64) (*accountWithUserInfo, error) {
	row, err := a.store.GetAccountById(ctx, id)
	if err != nil {
		return nil, err
	}
	return getAccountByIdRowToAccountWithUserInfo(row), nil
}

func (a *accountService) deleteAccountById(ctx context.Context, id int64) error {
	_, err := a.store.DeleteAccount(ctx, id)
	return err
}
