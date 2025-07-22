package entry

import db "github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/db/sqlc"

func dbEntryToEntryResponse(row db.Entry) entryResponse {
	return entryResponse{
		ID:        row.ID,
		AccountID: row.AccountID,
		Amount:    row.Amount,
		CreatedAt: row.CreatedAt,
	}
}

func dbEntriesToEntryResponses(src []db.Entry) []entryResponse {
	entries := make([]entryResponse, len(src))
	for i, row := range src {
		entries[i] = dbEntryToEntryResponse(row)
	}
	return entries
}
