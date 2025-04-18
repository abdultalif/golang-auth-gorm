package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)


func LoadENV() {
	env := os.Getenv("APP_ENV")
	if env == "production" {
        err := godotenv.Load(".env.production", "/app/.env.production", ".env", "/app/.env")
        if err != nil {
            log.Printf("⚠️ Warning: Failed to load .env files: %v", err)
            log.Println("ℹ️ Using environment variables from container")
        }
        log.Println("✅ Loaded environment variables")
	} else {
		err := godotenv.Load(".env.local")
		if err != nil {
			log.Fatal("❌ Failed to load .env.local")
		}
		log.Println("✅ Loaded .env.local")
	}
}
