package entry

import (
	"context"

	db "github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/db/sqlc"
	"golang.org/x/sync/errgroup"
)

type entryService struct {
	store db.Store
}

func (s *entryService) listEntries(ctx context.Context, args listEntriesRequest) (res listEntriesResponse, err error) {
	params := db.ListEntriesParams{Limit: args.Size, Offset: (args.Page - 1) * args.Size}

	var entryRows []db.Entry
	var totalCount int64

	g := errgroup.Group{}
	g.Go(func() (e error) {
		entryRows, e = s.store.ListEntries(ctx, params)
		return
	})
	g.Go(func() (e error) {
		totalCount, e = s.store.GetTotalEntriesCount(ctx)
		return
	})
	if err = g.Wait(); err != nil {
		return
	}

	return listEntriesResponse{
		Entries:    dbEntriesToEntryResponses(entryRows),
		TotalCount: totalCount,
	}, nil
}

func (s *entryService) getEntryById(ctx context.Context, id int64) (entryResponse, error) {
	dbRow, err := s.store.GetEntryById(ctx, id)
	return dbEntryToEntryResponse(dbRow), err
}

func (s *entryService) listEntriesByAccount(ctx context.Context, args listEntriesRequest, accountId int64) (res listEntriesResponse, err error) {
	params := db.ListEntriesByAccountParams{Limit: args.Size, Offset: (args.Page - 1) * args.Size, AccountID: accountId}

	var entryRows []db.Entry
	var totalCount int64

	g := errgroup.Group{}
	g.Go(func() (e error) {
		entryRows, e = s.store.ListEntriesByAccount(ctx, params)
		return
	})
	g.Go(func() (e error) {
		totalCount, e = s.store.GetTotalEntriesCountByAccount(ctx, accountId)
		return
	})
	if err = g.Wait(); err != nil {
		return
	}

	return listEntriesResponse{
		Entries:    dbEntriesToEntryResponses(entryRows),
		TotalCount: totalCount,
	}, nil
}
