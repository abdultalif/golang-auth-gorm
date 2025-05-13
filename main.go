package main

import (
	"net/http"
	"os"

	"github.com/abdultalif/golang-auth-gorm/config"
	"github.com/abdultalif/golang-auth-gorm/logger"
	"github.com/abdultalif/golang-auth-gorm/middlewares"
	"github.com/abdultalif/golang-auth-gorm/routes"
	"github.com/abdultalif/golang-auth-gorm/workers"
)

func main() {
	config.LoadENV()
	config.ConnectDB()
	config.InitRabbitMQ()
	workers.ConsumeVerificationQueue()

	router := routes.SetupRouter()

	port := os.Getenv("APP_PORT")
	logger.Log.Println("🚀 Server running on port:", port)
	logger.Log.Fatal(http.ListenAndServe(":"+port, middlewares.CORS(router)))
}