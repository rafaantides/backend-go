package pagination

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	DefaultPage     = 1
	DefaultPageSize = 10
	MaxPageSize     = 100
)

func GetPage(page *int) int {
	if page == nil || *page < 1 {
		return DefaultPage
	}
	return *page
}

func GetPageSize(pageSize *int) int {
	if pageSize == nil || *pageSize < 1 || *pageSize > MaxPageSize {
		return DefaultPageSize
	}
	return *pageSize
}

func GetOrderBy(orderBy *string, defaultOrder string, validColumns map[string]bool) (string, error) {
	if orderBy == nil || *orderBy == "" {
		return defaultOrder, nil
	}

	if !validColumns[*orderBy] {
		return "", errors.New("")
	}

	return *orderBy, nil
}

func SetPaginationHeaders(c *gin.Context, page, pageSize, total int) {
	totalPages := (total + pageSize - 1) / pageSize

	c.Header("X-Page", strconv.Itoa(page))
	c.Header("X-Page-Size", strconv.Itoa(pageSize))
	c.Header("X-Total-Count", strconv.Itoa(total))
	c.Header("X-Total-Pages", strconv.Itoa(totalPages))
}
