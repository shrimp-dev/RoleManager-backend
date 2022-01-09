package main

import (
	"drinkBack/database"
	"drinkBack/routes"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}

	client, err := database.NewClient()
	if err != nil {
		panic(err)
	}
	err = client.CheckConnection()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to database successfully")

	router := routes.Router{Client: client}
	port := ":" + os.Getenv("PORT")
	if port == ":" {
		port = ":3000"
	}
	fmt.Printf("Listening to %v\n", port)
	http.ListenAndServe(port, router.GenerateHandler())
}
