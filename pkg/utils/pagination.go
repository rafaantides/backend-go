package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func SetPaginationHeaders(c *gin.Context, page, pageSize, total int) {
	totalPages := (total + pageSize - 1) / pageSize

	c.Header("X-Page", strconv.Itoa(page))
	c.Header("X-Page-Size", strconv.Itoa(pageSize))
	c.Header("X-Total-Count", strconv.Itoa(total))
	c.Header("X-Total-Pages", strconv.Itoa(totalPages))
}
