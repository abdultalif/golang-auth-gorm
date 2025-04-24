package main

import (
	"net/http"
	"os"

	"github.com/abdultalif/golang-auth-gorm/config"
	"github.com/abdultalif/golang-auth-gorm/logger"
	"github.com/abdultalif/golang-auth-gorm/middlewares"
	"github.com/abdultalif/golang-auth-gorm/routes"
)

func main() {
	config.LoadENV()
	config.ConnectDB()

	router := routes.SetupRouter()

	port := os.Getenv("APP_PORT")
	logger.Log.Println("🚀 Server running on port:", port)
	logger.Log.Fatal(http.ListenAndServe(":"+port, middlewares.CORS(router)))
}