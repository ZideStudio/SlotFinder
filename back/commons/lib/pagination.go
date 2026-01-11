package lib

import (
	"app/commons/constants"
	"errors"

	"github.com/gin-gonic/gin"
)

type Pagination[T any] struct {
	Data   []T   `json:"data"`
	Page   int   `form:"page,default=1" json:"page" binding:"min=1,max=100"`
	Limit  int   `form:"limit,default=20" json:"limit" binding:"min=1,max=50"`
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
	p.Offset = (p.Page - 1) * p.Limit

	return nil
}
