package account

import db "github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/db/sqlc"

func listAccountRowsToAccountWithUserInfo(src []db.ListAccountsRow) []accountWithUserInfo {
	res := make([]accountWithUserInfo, len(src))
	for i, row := range src {
		res[i] = accountWithUserInfo{
			ID:        row.ID,
			UserID:    row.UserID,
			FullName:  row.FullName,
			Email:     row.Email,
			Currency:  row.Currency,
			Balance:   row.Balance,
			CreatedAt: row.CreatedAt,
		}
	}
	return res
}

func getAccountByIdRowToAccountWithUserInfo(src db.GetAccountByIdRow) *accountWithUserInfo {
	return &accountWithUserInfo{
		ID:        src.ID,
		UserID:    src.UserID,
		FullName:  src.FullName,
		Email:     src.Email,
		Currency:  src.Currency,
		Balance:   src.Balance,
		CreatedAt: src.CreatedAt,
	}
}
