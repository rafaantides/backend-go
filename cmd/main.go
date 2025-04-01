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

	// Configuração do CORS
	// r.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"*"}, // Frontend
	// 	// AllowOrigins:     []string{"http://localhost:3000"}, // Frontend
	// 	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	// 	AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
	// 	AllowCredentials: true, // Permitir cookies
	// }))

	// v1.Use(middlewares.CORSMiddleware())
	v1.Use(middlewares.UUIDMiddleware())
	v1.Use(middlewares.ErrorMiddleware())

	routes.RegisterDocsRoutes(r.Group("/docs/v1"))
	routes.RegisterDebtRoutes(v1.Group("/debts"))
	routes.RegisterInvoiceRoutes(v1.Group("/invoices"))
	routes.RegisterCategoryRoutes(v1.Group("/categories"))

	r.Run(":8080")
}
