package database

import (
	"context"
	"drinkBack/models"
	"drinkBack/utils"
	"fmt"
	"os"
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
	FindUserById(usrId primitive.ObjectID) (models.UserData, error)
	FindAllUsers() ([]models.UserData, error)
	FindDrinksOfUser(usrId primitive.ObjectID) ([]models.Drink, error)
	FindDebtsOfUser(usrId primitive.ObjectID) ([]models.Debt, error)
	UpdateUserById(usrId primitive.ObjectID, usr models.UserUpdate) (models.UserData, error)
	UpdateDrinksByIds(usrIds []primitive.ObjectID, done bool) ([]models.Drink, error)

	CreateNewDrink(drink models.Drink) (models.Drink, error)
	FindAllDrinks(usrId primitive.ObjectID) ([]models.Drink, error)

	CreateNewDebt(debt models.Debt) (models.Debt, error)
	FindAllDebts() ([]models.Debt, error)
	FindDebtById(debtId primitive.ObjectID) (models.Debt, error)
	PayDebt(query bson.M) (models.Debt, error)

	VerifyUserPassword(email string, password string, data *models.LoginResponse) (bool, error)
}

func NewClient() (DbClient, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URL")))
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

	// Check seed user
	var seed models.User
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	err = dbCl.getUserDatabase().FindOne(
		ctx,
		bson.M{
			"email": os.Getenv("SEED_MAIL"),
		},
	).Decode(&seed)

	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			if err := dbCl.createSeedUser(); err != nil {
				return nil, err
			}
		default:
			return nil, err
		}
	}

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

func (d *dbClient) createSeedUser() error {
	fmt.Println("(DATABASE) seed user not found, creating using informations in .env")
	if valid, err := utils.ValidatePassword([]byte(os.Getenv("SEED_PWD"))); err != nil || !valid {
		return fmt.Errorf("could not generate seed user. valid result: %v.\nError: %v", valid, err)
	}
	seedSalt, err := utils.GenerateRandomSalt(10)
	if err != nil {
		return fmt.Errorf("could not generate seed user. %v", err)
	}
	seedPwd, err := utils.HashPassword(os.Getenv("SEED_PWD"), seedSalt)
	if err != nil {
		return fmt.Errorf("could not generate seed user. %v", err)
	}
	_, err = d.getUserDatabase().InsertOne(context.Background(), models.User{
		UserData: models.UserData{
			Id: primitive.NewObjectID(),
			UserUpdate: models.UserUpdate{
				Name:  "Seed",
				Email: os.Getenv("SEED_MAIL"),
				Path:  "https://cdn.discordapp.com/attachments/580125063186087966/930244138325131264/59MvRFdT_400x400.png",
			},
		},
		Password: seedPwd,
		Salt:     seedSalt,
	})
	if err != nil {
		return fmt.Errorf("could not generate seed user. %v", err)
	}
	return nil
}
