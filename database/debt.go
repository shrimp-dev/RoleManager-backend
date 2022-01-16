package database

import (
	"context"
	"drinkBack/models"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

func (d *dbClient) PayDebt(query models.PayDebtRequest) (models.Debt, error) {

	var debtors []primitive.ObjectID
	for i := 0; i < len(query.Debtors); i++ {
		if debtor, err := primitive.ObjectIDFromHex(query.Debtors[i]); err != nil {
			return models.Debt{}, fmt.Errorf("invalid debtor")
		} else {
			debtors = append(debtors, debtor)
		}
	}

	id, err := primitive.ObjectIDFromHex(query.Id)
	if err != nil {
		return models.Debt{}, fmt.Errorf("invalid debtId")
	}

	creditor, err := primitive.ObjectIDFromHex(query.Creditor)
	if err != nil {
		return models.Debt{}, fmt.Errorf("invalid Creditor")
	}

	debtDb := d.getDebtDatabase()
	query_options := options.FindOneAndUpdate()
	rd := options.After
	af := options.ArrayFilters{
		Filters: []interface{}{
			bson.M{
				"elem._id": bson.M{
					"$in": debtors,
				},
			},
		}}

	query_options.ReturnDocument = &rd
	query_options.ArrayFilters = &af
	var debt models.Debt
	if err := debtDb.FindOneAndUpdate(context.Background(),
		bson.M{
			"_id":      id,
			"creditor": creditor,
		},
		bson.M{
			"$set": bson.M{
				"debtors.$[elem].paid": query.Paid,
			},
		},
		query_options,
	).Decode(&debt); err != nil {
		return models.Debt{}, err
	}
	return debt, nil
}
