package database

import (
	"context"
	"drinkBack/models"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
