package main

import (
	"log"
	"net/http"
	"os"

	"github.com/YoDobchev/Go-Online/src/database"
	"github.com/YoDobchev/Go-Online/src/routes"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	database.Connect()

	r := routes.New()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	http.ListenAndServe(":"+port, r)
}
