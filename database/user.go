package database

import (
	"context"
	"drinkBack/models"
	"drinkBack/utils"
	"errors"
	"fmt"
	"log"
	"net/mail"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d *dbClient) CreateNewUser(usr models.User) (models.UserData, error) {
	usr.Id = primitive.NewObjectID()

	userDb := d.getUserDatabase()
	_, err := userDb.InsertOne(context.TODO(), usr)

	filtered := models.UserData{
		Id: usr.Id,
		UserUpdate: models.UserUpdate{
			Name:  usr.Name,
			Email: usr.Email,
			Path:  usr.Path,
		},
	}
	return filtered, err
}

func (d *dbClient) UpdateUserById(usrId primitive.ObjectID, update models.UserUpdate) (models.UserData, error) {
	userDb := d.getUserDatabase()

	if update.Email != "" {
		if _, err := mail.ParseAddress(update.Email); err != nil {
			return models.UserData{}, err
		}
	}

	if update.Path != "" {
		if err := utils.ValidateUserPath(update.Path); err != nil {
			return models.UserData{}, err
		}
	}

	query_options := options.FindOneAndUpdate()
	rd := options.After
	query_options.ReturnDocument = &rd

	result_fnu := userDb.FindOneAndUpdate(context.Background(), bson.M{"_id": usrId}, bson.M{"$set": update}, query_options)
	var doc_upd models.UserData
	if err := result_fnu.Decode(&doc_upd); err != nil {
		fmt.Println("ERR")
		fmt.Println(err)
		fmt.Println(doc_upd)
		return models.UserData{}, err
	}
	return doc_upd, nil
}

func (d *dbClient) FindUserById(usrId primitive.ObjectID) (models.UserData, error) {
	cur, err := d.getUserDatabase().Find(context.TODO(), bson.M{"_id": usrId}, nil)
	if err != nil {
		return models.UserData{}, err
	}

	ok := cur.Next(context.TODO())
	if !ok {
		return models.UserData{}, errors.New("user not found")
	}

	var user models.UserData
	err = cur.Decode(&user)
	if err != nil {
		log.Fatal(err)
	}

	return user, nil
}

func (d *dbClient) FindAllUsers() ([]models.UserData, error) {

	cur, err := d.getUserDatabase().Find(context.TODO(), bson.D{{}}, nil)
	if err != nil {
		return nil, err
	}

	var users []models.UserData

	for cur.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var elem models.UserData
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		users = append(users, elem)

	}

	return users, nil
}

func (d *dbClient) VerifyUserPassword(email string, password string, data *models.LoginResponse) (bool, error) {
	cur, err := d.getUserDatabase().Find(context.TODO(), bson.M{"email": email}, nil)
	if err != nil {
		return false, err
	}

	ok := cur.Next(context.TODO())
	if !ok {
		return false, errors.New("user not found")
	}

	var user models.User
	err = cur.Decode(&user)
	if err != nil {
		return false, err
	}

	ok, err = utils.MatchPassword(user.Password, password, []byte(user.Salt))
	if err != nil || !ok {
		return false, err
	}

	token, err := utils.GenerateAuthenticationToken(user.Id.Hex(), utils.AUTH)
	if err != nil {
		return ok, err
	}

	*data = models.LoginResponse{
		UserData: models.UserData{
			Id: user.Id,
			UserUpdate: models.UserUpdate{
				Name:  user.Name,
				Email: user.Email,
				Path:  user.Path,
			},
		},
		Token: token,
	}
	return ok, nil
}
