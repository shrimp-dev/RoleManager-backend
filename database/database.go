package database

import (
	"context"
	"drinkBack/models"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type dbClient struct {
	//TODO: encapsulate client
	Client *mongo.Client
}

type DbClient interface {
	CreateNewUser(usr models.User) error
	CreateNewDrink(drink models.Drink) error
	FindUserById(usrId primitive.ObjectID) (models.User, error)
	FindAllUsers() ([]models.User, error)
	FindDrinksByUser(usrId primitive.ObjectID) ([]models.Drink, error)
}

func NewClient() (DbClient, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, err
	}

	dbCl := &dbClient{Client: client}

	return dbCl, nil
}

func (d *dbClient) getUserDatabase() *mongo.Collection {
	usersDb := d.Client.Database("testv2").Collection("users")
	return usersDb
}

func (d *dbClient) getDrinkDatabase() *mongo.Collection {
	usersDb := d.Client.Database("testv2").Collection("drinks")
	return usersDb
}

func (d *dbClient) CreateNewUser(usr models.User) (string, error) {
	usr.Id = primitive.NewObjectID()

	userDb := d.getUserDatabase()
	id, err := userDb.InsertOne(context.TODO(), usr)
	return id.InsertedID.(primitive.ObjectID).Hex(), err
}

func (d *dbClient) CreateNewDrink(drink models.Drink) (string, error) {
	drink.Id = primitive.NewObjectID()

	drinkDb := d.getDrinkDatabase()
	id, err := drinkDb.InsertOne(context.TODO(), drink)
	return id.InsertedID.(primitive.ObjectID).Hex(), err
}

func (d *dbClient) FindUserById(usrId primitive.ObjectID) (models.User, error) {

	cur, err := d.getUserDatabase().Find(context.TODO(), bson.D{{"_id", usrId}}, nil)
	if err != nil {
		return models.User{}, err
	}

	var user models.User

	cur.Next(context.TODO())
	//Create a value into which the single document can be decoded
	var elem models.User
	err = cur.Decode(&elem)
	if err != nil {
		log.Fatal(err)
	}

	user = elem

	return user, nil
}

func (d *dbClient) FindAllUsers() ([]models.User, error) {

	cur, err := d.getUserDatabase().Find(context.TODO(), bson.D{{}}, nil)
	if err != nil {
		return nil, err
	}

	var users []models.User

	for cur.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var elem models.User
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		users = append(users, elem)

	}

	return users, nil
}

func (d *dbClient) FindDrinksByUser(usrId primitive.ObjectID) ([]models.Drink, error) {

	cur, err := d.getDrinkDatabase().Find(context.TODO(), bson.D{{"usrId", usrId}}, nil)
	if err != nil {
		return nil, err
	}

	var drinks []models.Drink

	for cur.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var elem models.Drink
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		drinks = append(drinks, elem)

	}

	return drinks, nil
}
