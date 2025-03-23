package routes

import (
	"api-go/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterDebtRoutes(router *gin.RouterGroup) {
	router.POST("/", handlers.CreateDebtHandler)
	router.GET("/", handlers.ListDebtsHandler)
	router.GET("/:id", handlers.GetDebtByIDHandler)
	router.PUT("/:id", handlers.UpdateDebtHandler)
	router.DELETE("/:id", handlers.DeleteDebtHandler)
}

func RegisterInvoiceRoutes(router *gin.RouterGroup) {
	router.POST("/", handlers.CreateInvoiceHandler)
	router.GET("/", handlers.ListInvoicesHandler)
	router.GET("/:id", handlers.GetInvoiceByIDHandler)
	router.PUT("/:id", handlers.UpdateInvoiceHandler)
	router.DELETE("/:id", handlers.DeleteInvoiceHandler)
}
