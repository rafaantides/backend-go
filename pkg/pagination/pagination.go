package pagination

import (
	"backend-go/internal/api/errs"
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
		// TODO: voltar e rever isso daqui
		return nil, errs.ErrBadRequest
		// return nil, errs.ErrInvalidPage
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", DefaultPageSize))
	if err != nil || pageSize < 1 || pageSize > MaxPageSize {
		// TODO: voltar e rever isso daqui
		// return nil, errs.ErrInvalidPageSize
		return nil, errs.ErrBadRequest
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
		return errs.InvalidOrderBy(p.OrderBy)
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
