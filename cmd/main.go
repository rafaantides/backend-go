package main

import (
	"api-go/internal/repository"
	"api-go/internal/routes"

	"github.com/gin-gonic/gin"
)

// @title API GO
// @version 1.0
// @description Esta é a documentação da API
// @host localhost:8080
// @BasePath /
// @schemes http
func main() {
	repository.ConnectDB()

	r := gin.Default()

	r.StaticFile("/swagger.json", "./docs/swagger.json")
	r.Static("/static", "./static")

	r.GET("/", func(c *gin.Context) {
		c.File("./docs/docs.html")
	})

	api := r.Group("/api")

	routes.RegisterDebtRoutes(api.Group("/debts"))
	routes.RegisterInvoiceRoutes(api.Group("/invoices"))

	r.Run(":8080")
}
