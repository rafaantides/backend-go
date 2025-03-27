package routes

import (
	"backend-go/internal/api/v1/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterDocsRoutes(router *gin.RouterGroup) {
	router.StaticFile("/swagger.json", "./docs/v1/swagger.json")
	router.Static("/static", "./static")

	router.GET("/", func(c *gin.Context) {
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

func RegisterCategoryRoutes(router *gin.RouterGroup) {
	router.POST("/", handlers.CreateCategoryHandler)
	router.GET("/", handlers.ListCategorysHandler)
	router.GET("/:id", handlers.GetCategoryByIDHandler)
	router.PUT("/:id", handlers.UpdateCategoryHandler)
	router.DELETE("/:id", handlers.DeleteCategoryHandler)
}
