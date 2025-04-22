package main

import (
	"log"
	"net/http"
	"os"

	"github.com/abdultalif/golang-auth-gorm/config"
	"github.com/abdultalif/golang-auth-gorm/models"
	"github.com/abdultalif/golang-auth-gorm/routes"
)

func main() {
	config.LoadENV()
	config.ConnectDB()
	config.DB.AutoMigrate(&models.User{}, &models.VerificationCode{}, &models.PasswordReset{})

	router := routes.SetupRouter()

	port := os.Getenv("APP_PORT")
	log.Println("🚀 Server running on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}