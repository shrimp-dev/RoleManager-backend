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

func (d *dbClient) CreateNewUser(usr models.User) (models.GetUserResponse, error) {
	usr.Id = primitive.NewObjectID()

	userDb := d.getUserDatabase()
	_, err := userDb.InsertOne(context.TODO(), usr)

	filtered := models.GetUserResponse{
		Id:             usr.Id,
		Name:           usr.Name,
		Email:          usr.Email,
		Path:           usr.Path,
		WalletAccounts: usr.WalletAccounts,
		PixAccounts:    usr.PixAccounts,
		CreatedBy:      usr.CreatedBy,
	}
	return filtered, err
}

func (d *dbClient) UpdateUserById(usrId primitive.ObjectID, update models.UpdateUserRequest) (models.GetUserResponse, error) {
	userDb := d.getUserDatabase()

	if update.Email != "" {
		if _, err := mail.ParseAddress(update.Email); err != nil {
			return models.GetUserResponse{}, err
		}
	}

	if update.Path != "" {
		if err := utils.ValidateUserPath(update.Path); err != nil {
			return models.GetUserResponse{}, err
		}
	}

	query_options := options.FindOneAndUpdate()
	rd := options.After
	query_options.ReturnDocument = &rd

	result_fnu := userDb.FindOneAndUpdate(context.Background(), bson.M{"_id": usrId}, bson.M{"$set": update}, query_options)

	var doc_upd models.GetUserResponse
	if err := result_fnu.Decode(&doc_upd); err != nil {
		fmt.Println("ERR")
		fmt.Println(err)
		fmt.Println(doc_upd)
		return models.GetUserResponse{}, err
	}
	return doc_upd, nil
}

func (d *dbClient) FindUserById(usrId primitive.ObjectID) (models.GetUserResponse, error) {
	cur, err := d.getUserDatabase().Find(context.TODO(), bson.M{"_id": usrId}, nil)
	if err != nil {
		return models.GetUserResponse{}, err
	}

	ok := cur.Next(context.TODO())
	if !ok {
		return models.GetUserResponse{}, errors.New("user not found")
	}

	var user models.GetUserResponse
	err = cur.Decode(&user)
	if err != nil {
		log.Fatal(err)
	}

	return user, nil
}

func (d *dbClient) FindAllUsers() ([]models.GetUserResponse, error) {

	cur, err := d.getUserDatabase().Find(context.TODO(), bson.D{{}}, nil)
	if err != nil {
		return nil, err
	}

	var users []models.GetUserResponse

	for cur.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var elem models.GetUserResponse
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
		GetUserResponse: models.GetUserResponse{
			Id:    user.Id,
			Name:  user.Name,
			Email: user.Email,
			Path:  user.Path,
		},
		Token: token,
	}
	return ok, nil
}
