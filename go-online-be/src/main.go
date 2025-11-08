package main

import (
	"net/http"

	"github.com/YoDobchev/Go-Online/src/database"
	"github.com/YoDobchev/Go-Online/src/routes"
)

func main() {
	database.Connect()

	r := routes.New()

	http.ListenAndServe(":3000", r)
}
