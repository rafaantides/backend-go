// main.go
package main

import (
	"backend-go/internal/api/middlewares"
	"backend-go/internal/api/v1/repository"
	"backend-go/internal/api/v1/routes"

	"github.com/gin-gonic/gin"
)

// @title API GO
// @version 1.0
// @description Esta API é projetada para monitoramento de dívidas, ajudando a organizar financeiramente.
// @host localhost:8080
// @BasePath /api/v1
// @schemes http
func main() {
	repository.ConnectDB()

	r := gin.Default()
	v1 := r.Group("/api/v1")

	v1.Use(middlewares.UUIDValidatorMiddleware())

	routes.RegisterDocsRoutes(r.Group("/docs/v1"))
	routes.RegisterDebtRoutes(v1.Group("/debts"))
	routes.RegisterInvoiceRoutes(v1.Group("/invoices"))

	r.Run(":8080")
}
