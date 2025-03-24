// main.go
package main

import (
	"api-go/internal/repository"
	"api-go/internal/routes"

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
	api := r.Group("/api/v1")

	routes.RegisterDocsRoutes(r.Group("/docs/v1"))
	routes.RegisterDebtRoutes(api.Group("/debts"))
	routes.RegisterInvoiceRoutes(api.Group("/invoices"))

	r.Run(":8080")
}
