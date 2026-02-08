package main

import (
	"log"
	"net/http"
	"os"

	"github.com/YoDobchev/Go-Online/src/database"
	gogame "github.com/YoDobchev/Go-Online/src/game/go"
	"github.com/YoDobchev/Go-Online/src/katago"
	"github.com/YoDobchev/Go-Online/src/routes"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	database.Connect()

	gogame.LoadGamesFromDB()

	katago.InitEng()

	r := routes.New()

	port := os.Getenv("PORT")

	log.Printf("Listening on :%v...", port)
	http.ListenAndServe(":"+port, r)
}
