package main

import (
	"drinkBack/database"
	"drinkBack/routes"
	"net/http"
)

func main() {
	client, err := database.NewClient()
	if err != nil {
		panic(err)
	}
	router := routes.Router{Client: client}
	http.ListenAndServe(":3000", router.GenerateHandler())
}
