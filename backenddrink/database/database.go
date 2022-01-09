package database

import (
	"context"
	"drinkBack/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbName = "testv4"
)

type dbClient struct {
	//TODO: encapsulate client
	Client *mongo.Client
}

type DbClient interface {
	CheckConnection() error
	CreateNewUser(usr models.User) (models.UserData, error)
	CreateNewDrink(drink models.Drink) (models.Drink, error)
	CreateNewDebt(debt models.Debt) (models.Debt, error)
	FindUserById(usrId primitive.ObjectID) (models.UserData, error)
	FindAllUsers() ([]models.User, error)
	VerifyUserPassword(email string, password string, data *models.LoginResponse) (bool, error)
	FindDrinksOfUser(usrId primitive.ObjectID) ([]models.Drink, error)
	FindAllDrinks(usrId primitive.ObjectID) ([]models.Drink, error)
	FindAllDebts() ([]models.Debt, error)
	UpdateUserById(usrId primitive.ObjectID, usr models.User) (models.UserData, error)
	UpdateDrinksByIds(usrIds []primitive.ObjectID, done bool) ([]models.Drink, error)
	PayDebt(query bson.M) (models.Debt, error)
	FindDebtById(debtId primitive.ObjectID) (models.Debt, error)
	FindDebtsOfUser(usrId primitive.ObjectID) ([]models.Debt, error)
}

func NewClient() (DbClient, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, err
	}

	if _, err := client.Database(dbName).Collection("users").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"email": int32(1)},
		Options: options.Index().SetUnique(true),
	}); err != nil {
		return nil, err
	}

	dbCl := &dbClient{Client: client}

	return dbCl, nil
}

func (d *dbClient) CheckConnection() error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := d.Client.Ping(ctx, nil)
	if err != nil {
		return err
	}
	return nil
}

func (d *dbClient) getUserDatabase() *mongo.Collection {
	usersDb := d.Client.Database(dbName).Collection("users")
	return usersDb
}

func (d *dbClient) getDrinkDatabase() *mongo.Collection {
	usersDb := d.Client.Database(dbName).Collection("drinks")
	return usersDb
}

func (d *dbClient) getDebtDatabase() *mongo.Collection {
	debtDb := d.Client.Database(dbName).Collection("debt")
	return debtDb
}
