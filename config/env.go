package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func DBURI() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("env file not found")
	}

	dburi := os.Getenv("DATABASE_URL")
	if dburi == "" {
		log.Fatal("DBURI not set")
	}
	return dburi
}
