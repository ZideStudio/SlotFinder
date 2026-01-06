package lib

type PaginationQuery struct {
	Page  int `form:"page,default=1"`
	Limit int `form:"limit,default=20"`
}

type PaginationResult[T any] struct {
	Data  []T   `json:"data"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

func GetPaginationResult[T any](data []T, query PaginationQuery, total int64) PaginationResult[T] {
	return PaginationResult[T]{
		Data:  data,
		Page:  query.Page,
		Limit: query.Limit,
		Total: total,
	}
}
