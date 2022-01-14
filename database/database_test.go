package database_test

import (
	"drinkBack/database"
	"drinkBack/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateClient(t *testing.T) {
	client, err := database.NewClient()

	assert.NotNil(t, client)
	assert.Nil(t, err)
}

func TestCreateUser(t *testing.T) {
	client, err := database.NewClient()
	assert.Nil(t, err)

	_, err = client.CreateNewUser(models.User{UserData: models.UserData{Name: "foo", Path: "bar"}})
	assert.Nil(t, err)

}
