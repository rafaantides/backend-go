package main

import (
	"api-go/internal/repository"
	"api-go/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	repository.ConnectDB()

	r := gin.Default()

	api := r.Group("/api")
	routes.RegisterDebtRoutes(api.Group("/debts"))

	r.Run(":8080")
}
