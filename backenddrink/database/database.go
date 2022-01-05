package database

import (
	"context"
	"drinkBack/models"
	"drinkBack/utils"
	"errors"
	"log"
	"fmt"

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
	CreateNewUser(usr models.User) (models.User, error)
	CreateNewDrink(drink models.Drink) (models.Drink, error)
	CreateNewDebt(debt models.Debt) (models.Debt, error)
	FindUserById(usrId primitive.ObjectID) (models.User, error)
	FindAllUsers() ([]models.User, error)
	FindDrinksOfUser(usrId primitive.ObjectID) ([]models.Drink, error)
	FindAllDrinks(usrId primitive.ObjectID) ([]models.Drink, error)
	FindAllDebts() ([]models.Debt, error)
	UpdateUserById(usrId primitive.ObjectID, usr models.User) (models.User, error)
	UpdateDrinksByIds(usrIds []primitive.ObjectID, done bool) ([]models.Drink, error)
	PayDebt(query bson.M) (models.Debt, error)
	FindDebtById(debtId primitive.ObjectID) (models.Debt, error)
	FindDebtsOfUser(usrId primitive.ObjectID) ([]models.Debt, error)
}

func NewClient() (DbClient, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://mongo:27017"))
	if err != nil {
		return nil, err
	}

	dbCl := &dbClient{Client: client}

	return dbCl, nil
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

// User
func (d *dbClient) CreateNewUser(usr models.User) (models.User, error) {
	fmt.Println("CreateNewUser")
	usr.Id = primitive.NewObjectID()

	userDb := d.getUserDatabase()
	_, err := userDb.InsertOne(context.TODO(), usr)
	return usr, err
}

func (d *dbClient) UpdateUserById(usrId primitive.ObjectID, usr models.User) (models.User, error) {
	userDb := d.getUserDatabase()

	update := make(bson.M)
	if usr.Name != "" {
		update["name"] = usr.Name
	}
	if usr.Path != "" {
		if err := utils.ValidateUserPath(usr.Path); err != nil {
			return models.User{}, err
		}
		update["path"] = usr.Path
	}
	if len(update) == 0 {
		return models.User{}, errors.New("no update in body")
	}

	query_options := options.FindOneAndUpdate()
	rd := options.After
	query_options.ReturnDocument = &rd
	result_fnu := userDb.FindOneAndUpdate(context.Background(), bson.M{"_id": usrId}, bson.M{"$set": update}, query_options)
	var doc_upd models.User
	if err := result_fnu.Decode(&doc_upd); err != nil {
		return models.User{}, err
	}
	return doc_upd, nil
}

func (d *dbClient) FindUserById(usrId primitive.ObjectID) (models.User, error) {
	cur, err := d.getUserDatabase().Find(context.TODO(), bson.M{"_id": usrId}, nil)
	if err != nil {
		return models.User{}, err
	}
	var user models.User
	ok := cur.Next(context.TODO())
	if !ok {
		return models.User{}, errors.New("user not found")
	}
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

// Drink
func (d *dbClient) CreateNewDrink(drink models.Drink) (models.Drink, error) {
	drink.Id = primitive.NewObjectID()

	drinkDb := d.getDrinkDatabase()
	_, err := drinkDb.InsertOne(context.TODO(), drink)
	return drink, err
}

func (d *dbClient) FindAllDrinks(usrId primitive.ObjectID) ([]models.Drink, error) {

	cur, err := d.getDrinkDatabase().Find(context.TODO(), bson.M{}, nil)
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

func (d *dbClient) FindDrinksOfUser(usrId primitive.ObjectID) ([]models.Drink, error) {

	cur, err := d.getDrinkDatabase().Find(context.TODO(), bson.M{"usrId": usrId}, nil)
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

func (d *dbClient) UpdateDrinksByIds(usrIds []primitive.ObjectID, done bool) ([]models.Drink, error) {
	drinkDb := d.getDrinkDatabase()
	filter := bson.M{
		"_id": bson.M{
			"$in": usrIds,
		},
	}
	query := bson.M{
		"$set": bson.M{
			"done": done,
		},
	}
	_, err := drinkDb.UpdateMany(context.Background(), filter, query)
	if err != nil {
		return nil, err
	}
	cursor, err := drinkDb.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	var drinks []models.Drink
	if err := cursor.All(context.Background(), &drinks); err != nil {
		return nil, err
	}
	return drinks, nil
}

// Debt
func (d *dbClient) CreateNewDebt(debt models.Debt) (models.Debt, error) {
	debt.Id = primitive.NewObjectID()

	debtDb := d.getDebtDatabase()
	_, err := debtDb.InsertOne(context.TODO(), debt)
	return debt, err
}

func (d *dbClient) FindDebtById(debtId primitive.ObjectID) (models.Debt, error) {
	debtDb := d.getDebtDatabase()
	var debt models.Debt
	if err := debtDb.FindOne(context.Background(), bson.M{"_id": debtId}).Decode(&debt); err != nil {
		return models.Debt{}, err
	}
	return debt, nil
}

func (d *dbClient) FindDebtsOfUser(usrId primitive.ObjectID) ([]models.Debt, error) {
	filter := bson.M{
		"debtors._id": usrId,
	}
	cur, err := d.getDebtDatabase().Find(context.TODO(), filter, nil)
	if err != nil {
		return nil, err
	}

	var debts []models.Debt

	for cur.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var elem models.Debt
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		debts = append(debts, elem)

	}

	return debts, nil
}

func (d *dbClient) FindAllDebts() ([]models.Debt, error) {
	cur, err := d.getDebtDatabase().Find(context.TODO(), bson.D{{}}, nil)
	if err != nil {
		return nil, err
	}

	var debts []models.Debt

	for cur.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var elem models.Debt
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		debts = append(debts, elem)

	}

	return debts, nil
}

func (d *dbClient) PayDebt(query bson.M) (models.Debt, error) {
	debtDb := d.getDebtDatabase()
	query_options := options.FindOneAndUpdate()
	rd := options.After
	query_options.ReturnDocument = &rd
	var debt models.Debt
	if err := debtDb.FindOneAndUpdate(context.Background(),
		bson.M{
			"_id":         query["_id"].(primitive.ObjectID),
			"debtors._id": query["usrId"].(primitive.ObjectID),
		},
		bson.M{
			"$set": bson.M{
				"debtors.$.paid": query["paid"].(bool),
			},
		},
		query_options,
	).Decode(&debt); err != nil {
		return models.Debt{}, err
	}
	return debt, nil
}
