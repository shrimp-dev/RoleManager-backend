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

	_, err = client.CreateNewUser(models.User{Name: "Edson", Path: "teste"})
	if err != nil {
		panic(err)
	}

	v, err := client.FindAllUsers()
	if err != nil {
		panic(err)
	}

	client.CreateNewDrink(models.Drink{UsrId: v[0].Id, Name: "Tequila"})
	if err != nil {
		panic(err)
	}

	d, err := client.FindDrinksByUser(v[0].Id)
	if err != nil {
		panic(err)
	}

	fmt.Println(v)
	fmt.Println(d)

}
