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
	// Filter   *string `json:"filter"`
	// TotalRows int64  `json:"total_rows"`
	// TotalPages int   `json:"total_pages"`
}

func NewPagination(c *gin.Context) (*Pagination, error) {
	page, err := strconv.Atoi(c.DefaultQuery("page", DefaultPage))
	if err != nil || page < 1 {
		return nil, errs.ErrInvalidPage
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", DefaultPageSize))
	if err != nil || pageSize < 1 || pageSize > MaxPageSize {
		return nil, errs.ErrInvalidPageSize
	}

	orderBy := c.DefaultQuery("order_by", "")
	search := c.Query("search")
	// filter := c.Query("filter")

	// totalPages := int(math.Ceil(float64(totalRows) / float64(pageSize)))

	return &Pagination{
		Page:     page,
		PageSize: pageSize,
		OrderBy:  orderBy,
		Search:   search,
		// Filter:   filter,
		// TotalRows: totalRows,
		// TotalPages: totalPages,
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
		return errs.ErrInvalidOrderBy(p.OrderBy)
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

// func ParseFilter(filter string, model interface{}) interface{} {
// 	if filter == "" {
// 		return reflect.New(reflect.TypeOf(model)).Elem().Interface()
// 	}

// 	pairs := strings.Split(filter, ",")
// 	for _, pair := range pairs {
// 		parts := strings.Split(pair, ":")
// 		if len(parts) == 2 {
// 			key := strings.TrimSpace(parts[0])
// 			value := strings.TrimSpace(parts[1])

// 			reflectValue := reflect.ValueOf(model).Elem()
// 			reflectType := reflectValue.Type()

// 			for i := 0; i < reflectType.NumField(); i++ {
// 				field := reflectType.Field(i)
// 				if strings.EqualFold(field.Name, key) {
// 					reflectFieldValue := reflectValue.Field(i)
// 					if reflectFieldValue.CanSet() {
// 						reflectFieldValue.SetString(value)
// 					}
// 					break
// 				}
// 			}
// 		}
// 	}
// 	return model
// }
