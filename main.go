package main

import (
	"akh/blog/config"
	"akh/blog/routes"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	config.ConnectDB()
	routes.RunRoute()
}
