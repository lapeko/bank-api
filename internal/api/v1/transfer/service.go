package transfer

import (
	"context"

	db "github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/db/sqlc"
	"golang.org/x/sync/errgroup"
)

type transferService struct {
	store db.Store
}

func (s *transferService) transferMoney(ctx context.Context, args transferRequest) error {
	return s.store.TransferMoney(ctx, args.AccountFrom, args.AccountTo, args.Amount)
}

func (s *transferService) transferIn(ctx context.Context, args externalTransferRequest) error {
	_, err := s.store.TransferExternalMoney(ctx, args.Account, args.Amount)
	return err
}

func (s *transferService) transferOut(ctx context.Context, args externalTransferRequest) (accountResponse, error) {
	acc, err := s.store.TransferExternalMoney(ctx, args.Account, args.Amount)
	return dbAccoutToAccountResponse(acc), err
}

func (s *transferService) listTransfers(ctx context.Context, args listTransfersRequest) (res listTransfersResponse, err error) {
	var transfers []db.Transfer
	var totalCount int64

	params := db.ListTransfersParams{Limit: args.Size, Offset: (args.Page - 1) * args.Size}

	g := errgroup.Group{}
	g.Go(func() (e error) {
		transfers, e = s.store.ListTransfers(ctx, params)
		return
	})
	g.Go(func() (e error) {
		totalCount, e = s.store.GetTotalTransfersCount(ctx)
		return
	})
	if err = g.Wait(); err != nil {
		return
	}

	res.Transfers = dbTransfersToTransferResponses(transfers)
	res.TotalCount = totalCount
	return res, nil
}

func (s *transferService) listTransfersByAccount(ctx context.Context, args listTransfersRequest, accountId int64) (res listTransfersResponse, err error) {
	var transfers []db.Transfer
	var totalCount int64

	params := db.ListTransfersByAccountParams{AccoutID: accountId, Limit: args.Size, Offset: (args.Page - 1) * args.Size}

	g := errgroup.Group{}
	g.Go(func() (e error) {
		transfers, e = s.store.ListTransfersByAccount(ctx, params)
		return
	})
	g.Go(func() (e error) {
		totalCount, e = s.store.GetTotalTransfersCountByAccount(ctx, accountId)
		return
	})
	if err = g.Wait(); err != nil {
		return
	}

	res.Transfers = dbTransfersToTransferResponses(transfers)
	res.TotalCount = totalCount
	return res, nil
}
