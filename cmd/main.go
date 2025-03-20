package main

import (
	"api-go/internal/handlers"
	"api-go/internal/repository"

	"github.com/gin-gonic/gin"
)

func main() {
	repository.ConnectDB()

	r := gin.Default()

	// Criando rota POST
	r.POST("/debts", handlers.CreateDebtHandler)

	r.Run(":8080")
}