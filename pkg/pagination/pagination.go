package pagination

import (
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

func SetPaginationHeaders(c *gin.Context, page, pageSize, total int) {
	totalPages := (total + pageSize - 1) / pageSize

	c.Header("X-Page", strconv.Itoa(page))
	c.Header("X-Page-Size", strconv.Itoa(pageSize))
	c.Header("X-Total-Count", strconv.Itoa(total))
	c.Header("X-Total-Pages", strconv.Itoa(totalPages))
}
