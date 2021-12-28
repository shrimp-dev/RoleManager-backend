package main

import (
	"drinkBack/database"
	"drinkBack/models"
	"fmt"
)

func main() {
	client, err := database.NewClient()
	if err != nil {
		panic(err)
	}

	err = client.CreateNewUser(models.User{Name: "Edson", Path: "teste"})
	if err != nil {
		panic(err)
	}

	v, err := client.FindAllUsers()
	if err != nil {
		panic(err)
	}

	fmt.Println(v)
}
