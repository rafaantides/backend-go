package routes

import (
	"backend-go/internal/api/v1/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterDocsRoutes(router *gin.RouterGroup) {
	router.StaticFile("/swagger.json", "./docs/v1/swagger.json")

	router.GET("", func(c *gin.Context) {
		c.File("./docs/v1/docs.html")
	})
	router.GET("/redoc", func(c *gin.Context) {
		c.File("./docs/v1/redoc.html")
	})
	router.GET("/scalar", func(c *gin.Context) {
		c.File("./docs/v1/scalar.html")
	})
	router.GET("/swagger", func(c *gin.Context) {
		c.File("./docs/v1/swagger.html")
	})

}

func RegisterDebtRoutes(router *gin.RouterGroup, handler *handlers.DebtHandler) {
	router.POST("", handler.CreateDebtHandler)
	router.GET("", handler.ListDebtsHandler)
	router.GET("/:id", handler.GetDebtByIDHandler)
	router.PUT("/:id", handler.UpdateDebtHandler)
	router.DELETE("/:id", handler.DeleteDebtHandler)
}

func RegisterInvoiceRoutes(router *gin.RouterGroup, handler *handlers.InvoiceHandler) {
	router.POST("", handler.CreateInvoiceHandler)
	router.GET("", handler.ListInvoicesHandler)
	router.GET("/:id", handler.GetInvoiceByIDHandler)
	router.PUT("/:id", handler.UpdateInvoiceHandler)
	router.DELETE("/:id", handler.DeleteInvoiceHandler)
}

func RegisterCategoryRoutes(router *gin.RouterGroup, handler *handlers.CategoryHandler) {
	router.POST("", handler.CreateCategoryHandler)
	router.GET("", handler.ListCategorysHandler)
	router.GET("/:id", handler.GetCategoryByIDHandler)
	router.PUT("/:id", handler.UpdateCategoryHandler)
	router.DELETE("/:id", handler.DeleteCategoryHandler)
}

func RegisterPaymentStatusRoutes(router *gin.RouterGroup, handler *handlers.PaymentStatusHandler) {
	router.POST("", handler.CreatePaymentStatusHandler)
	router.GET("", handler.ListPaymentStatussHandler)
	router.GET("/:id", handler.GetPaymentStatusByIDHandler)
	router.PUT("/:id", handler.UpdatePaymentStatusHandler)
	router.DELETE("/:id", handler.DeletePaymentStatusHandler)
}
