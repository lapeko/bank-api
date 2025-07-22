package transfer

type transferRequest struct {
	AccountFrom int64 `json:"accountFrom" binding:"required,gte=1"`
	AccountTo   int64 `json:"accountTo" binding:"required,gte=1"`
	Amount      int64 `json:"amount" binding:"required,gte=1"`
}

type externalTransferRequest struct {
	Account int64 `json:"accountFrom" binding:"required,gte=1"`
	Amount  int64 `json:"amount" binding:"required"`
}
