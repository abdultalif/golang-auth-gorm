package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)


func LoadENV() {
	env := os.Getenv("APP_ENV")
	if env == "production" {
		err := godotenv.Load(".env.docker")
		if err != nil {
			log.Fatal("❌ Failed to load .env.docker")
		}
		log.Println("✅ Loaded .env.docker")
	} else {
		err := godotenv.Load(".env.local")
		if err != nil {
			log.Fatal("❌ Failed to load .env.local")
		}
		log.Println("✅ Loaded .env.local")
	}
}
