package helpers

import (
	"math"
	"strconv"

	"gorm.io/gorm"
)

type MetaPagination struct {
	ItemsPerPage int    `json:"itemsPerPage"`
	TotalItems   int64  `json:"totalItems"`
	CurrentPage  int    `json:"currentPage"`
	TotalPages   int    `json:"totalPages"`
	SortBy       string `json:"sortBy"`
}

type Pagination struct {
	Data any            `json:"data"`
	Meta MetaPagination `json:"meta"`
	Link string         `json:"link"`
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination) GetLimit() int {
	if p.Meta.ItemsPerPage == 0 {
		p.Meta.ItemsPerPage = 20
	}
	return p.Meta.ItemsPerPage
}

func (p *Pagination) GetPage() int {
	if p.Meta.CurrentPage == 0 {
		p.Meta.CurrentPage = 1
	}
	return p.Meta.CurrentPage
}

func (p *Pagination) GetSort() string {
	if p.Meta.SortBy == "" {
		p.Meta.SortBy = "Id desc"
	}
	return p.Meta.SortBy
}

func Paginate(value any, pagination *Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var totalRows int64
	db.Model(value).Count(&totalRows)

	pagination.Meta.TotalItems = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.Meta.ItemsPerPage)))
	pagination.Meta.TotalPages = totalPages

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(pagination.GetSort())
	}
}

func SetConds(queries map[string][]string, pagination *Pagination) map[string]any {

	queryStrings := make(map[string]any)

	if queries != nil {

		limit := queries["limit"]

		if len(limit) != 0 {

			pagination.Meta.ItemsPerPage, _ = strconv.Atoi(limit[0])
		}

		delete(queries, "limit")

		page := queries["page"]

		if len(page) != 0 {

			pagination.Meta.CurrentPage, _ = strconv.Atoi(page[0])
		}

		delete(queries, "page")

		for key, value := range queries { // Order not specified
			if len(value) != 0 {
				v := value[0]
				queryStrings[key] = v
			}
		}

	}

	return queryStrings
}
