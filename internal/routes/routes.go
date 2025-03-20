package routes

import (
	"api-go/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterDebtRoutes(router *gin.RouterGroup) {
	router.POST("/", handlers.CreateDebtHandler)
	// router.GET("/", handlers.GetAllDebtsHandler)
	// router.GET("/:id", handlers.GetDebtByIDHandler)
	// router.PUT("/:id", handlers.UpdateDebtHandler)
	// router.DELETE("/:id", handlers.DeleteDebtHandler)
}
