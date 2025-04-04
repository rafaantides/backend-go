package cmd

import (
	"backend-go/internal/api/middlewares"
	"backend-go/internal/api/v1/handlers"
	queue "backend-go/internal/api/v1/queue/interfaces"
	"backend-go/internal/api/v1/queue/rabbitmq"
	repository "backend-go/internal/api/v1/repository/interfaces"
	"backend-go/internal/api/v1/repository/postgresql"
	"backend-go/internal/api/v1/routes"
	"backend-go/internal/api/v1/services"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var (
	apiPort     string
	enableCORS  bool
	enableDebug bool
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Start the API server",
	Run: func(cmd *cobra.Command, args []string) {
		startAPIServer()
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)
	apiCmd.Flags().StringVarP(&apiPort, "port", "p", "8080", "Port to run API server on")
	apiCmd.Flags().BoolVar(&enableCORS, "cors", false, "Enable CORS middleware")
	apiCmd.Flags().BoolVar(&enableDebug, "debug", false, "Enable debug mode")
}

func startAPIServer() {
	if enableDebug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	_ = godotenv.Load()

	db := connectDatabase()
	defer db.Close()

	mq := connectQueue()
	defer mq.Close()

	r := setupRouter(db, mq)

	for _, route := range r.Routes() {
		fmt.Printf("[%s] %s\n", route.Method, route.Path)
	}

	r.Run(":" + apiPort)
}

func connectDatabase() *postgresql.PostgreSQL {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := postgresql.NewPostgreSQL(dsn)
	if err != nil {
		log.Fatalf("Falha ao conectar ao banco: %v", err)
	}
	return db
}

func connectQueue() queue.MessageQueue {
	amqpURI := os.Getenv("AMQP_URI")
	queueName := os.Getenv("QUEUE_NAME")

	mq, err := rabbitmq.NewRabbitMQ(amqpURI, queueName)

	if err != nil {
		log.Fatalf("Falha ao conectar a fila: %v", err)
	}

	return mq

}

func setupRouter(db repository.Database, mq queue.MessageQueue) *gin.Engine {
	debtService := services.NewDebtService(db, mq)
	debtHandler := handlers.NewDebtHandler(debtService)

	invoiceService := services.NewInvoiceService(db)
	invoiceHandler := handlers.NewInvoiceHandler(invoiceService)

	categoryService := services.NewCategoryService(db)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	paymentStatusService := services.NewPaymentStatusService(db)
	paymentStatusHandler := handlers.NewPaymentStatusHandler(paymentStatusService)

	r := gin.Default()
	v1 := r.Group("/api/v1")

	if enableCORS {
		v1.Use(middlewares.CORSMiddleware())
	}
	v1.Use(middlewares.UUIDMiddleware())
	v1.Use(middlewares.ErrorMiddleware())

	r.Static("/static", "./static")
	routes.RegisterDocsRoutes(r.Group("/docs/v1"))
	routes.RegisterDebtRoutes(v1.Group("/debts"), debtHandler)
	routes.RegisterInvoiceRoutes(v1.Group("/invoices"), invoiceHandler)
	routes.RegisterCategoryRoutes(v1.Group("/categories"), categoryHandler)
	routes.RegisterPaymentStatusRoutes(v1.Group("/payment_status"), paymentStatusHandler)

	return r
}
