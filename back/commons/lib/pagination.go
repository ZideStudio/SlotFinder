package lib

import (
	"app/commons/constants"
	"errors"

	"github.com/gin-gonic/gin"
)

type PaginationQuery struct {
	Page   int `form:"page,default=1"`
	Limit  int `form:"limit,default=20"`
	Offset int `form:"-"`
}

type PaginationResult[T any] struct {
	Data  []T   `json:"data"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

func ParsePaginationQuery(c *gin.Context, p *PaginationQuery) error {
	if c == nil {
		return errors.New("context is nil")
	}

	if err := c.ShouldBindQuery(&p); err != nil {
		return constants.ERR_INVALID_PAGINATION_PARAMS.Err
	}

	if p.Page < 1 {
		return constants.ERR_INVALID_PAGINATION_PAGE.Err
	}
	if p.Limit < 1 || p.Limit > 100 {
		return constants.ERR_INVALID_PAGINATION_LIMIT.Err
	}
	p.Offset = (p.Page - 1) * p.Limit

	return nil
}

func GetPaginationResult[T any](data []T, query PaginationQuery, total int64) PaginationResult[T] {
	return PaginationResult[T]{
		Data:  data,
		Page:  query.Page,
		Limit: query.Limit,
		Total: total,
	}
}
