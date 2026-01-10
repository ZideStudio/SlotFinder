package lib

import (
	"app/commons/constants"
	"errors"

	"github.com/gin-gonic/gin"
)

type Pagination[T any] struct {
	Datas  []T   `json:"datas"`
	Page   int   `form:"page,default=1" json:"page"`
	Limit  int   `form:"limit,default=20" json:"limit"`
	Offset int   `form:"-" json:"-"`
	Total  int64 `json:"total"`
}

func (p *Pagination[T]) ParseQuery(c *gin.Context) error {
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
