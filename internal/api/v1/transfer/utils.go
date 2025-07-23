package transfer

import db "github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/db/sqlc"

func dbAccoutToAccountResponse(row db.Account) accountResponse {
	return accountResponse{
		ID:        row.ID,
		UserID:    row.UserID,
		Currency:  row.Currency,
		Balance:   row.Balance,
		CreatedAt: row.CreatedAt,
	}
}

func dbTransferToTransferResoinse(row db.Transfer) transferResponse {
	return transferResponse{
		ID:          row.ID,
		AccountFrom: row.AccountFrom,
		AccountTo:   row.AccountTo,
		Amount:      row.Amount,
		CreatedAt:   row.CreatedAt,
	}
}

func dbTransfersToTransferResponses(src []db.Transfer) []transferResponse {
	transfers := make([]transferResponse, len(src))
	for i, row := range src {
		transfers[i] = dbTransferToTransferResoinse(row)
	}
	return transfers
}
