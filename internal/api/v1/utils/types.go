package utils

type UriId struct {
	ID int64 `uri:"id" binding:"required,gte=1"`
}
