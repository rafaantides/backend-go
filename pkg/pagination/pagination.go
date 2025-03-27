package pagination

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	DefaultPage     = "1"
	DefaultPageSize = "10"
	MaxPageSize     = 100
)

type Pagination struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	OrderBy  string `json:"order_by"`
	Search   string `json:"search"`
}

func NewPagination(c *gin.Context) (*Pagination, error) {
	page, err := strconv.Atoi(c.DefaultQuery("page", DefaultPage))
	if err != nil || page < 1 {
		return nil, fmt.Errorf("%s: %s", "page", "valor invalido")
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", DefaultPageSize))
	if err != nil || pageSize < 1 || pageSize > MaxPageSize {
		return nil, fmt.Errorf("%s: %s", "page_size", "valor invalido")
	}

	orderBy := c.Query("order_by")
	search := c.Query("search")

	return &Pagination{
		Page:     page,
		PageSize: pageSize,
		OrderBy:  orderBy,
		Search:   search,
	}, nil
}

func (p *Pagination) Offset() int {
	return (p.Page - 1) * p.PageSize
}

func (p *Pagination) ValidateOrderBy(defaultOrder string, validColumns map[string]bool) error {
	if p.OrderBy == "" {
		p.OrderBy = defaultOrder
		return nil
	}
	if p.OrderBy != "" && !validColumns[p.OrderBy] {
		return fmt.Errorf("%s: %s", "order_by", "valor invalido")
	}
	return nil
}

func (p *Pagination) SetPaginationHeaders(c *gin.Context, total int) {
	totalPages := (total + p.PageSize - 1) / p.PageSize

	c.Header("X-Page", strconv.Itoa(p.Page))
	c.Header("X-Page-Size", strconv.Itoa(p.PageSize))
	c.Header("X-Total-Count", strconv.Itoa(total))
	c.Header("X-Total-Pages", strconv.Itoa(totalPages))
}
