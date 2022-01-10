package database

import (
	"context"
	"drinkBack/models"
	"drinkBack/utils"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d *dbClient) CreateNewUser(usr models.User) (models.UserData, error) {
	usr.Id = primitive.NewObjectID()

	userDb := d.getUserDatabase()
	_, err := userDb.InsertOne(context.TODO(), usr)

	filtered := models.UserData{
		Id:    usr.Id,
		Name:  usr.Name,
		Email: usr.Email,
		Path:  usr.Path,
	}
	return filtered, err
}

func (d *dbClient) UpdateUserById(usrId primitive.ObjectID, usr models.User) (models.UserData, error) {
	userDb := d.getUserDatabase()

	update := make(bson.M)
	if usr.Name != "" {
		update["name"] = usr.Name
	}
	if usr.Path != "" {
		if err := utils.ValidateUserPath(usr.Path); err != nil {
			return models.UserData{}, err
		}
		update["path"] = usr.Path
	}
	if len(update) == 0 {
		return models.UserData{}, errors.New("no update in body")
	}

	query_options := options.FindOneAndUpdate()
	rd := options.After
	query_options.ReturnDocument = &rd
	result_fnu := userDb.FindOneAndUpdate(context.Background(), bson.M{"_id": usrId}, bson.M{"$set": update}, query_options)
	var doc_upd models.UserData
	if err := result_fnu.Decode(&doc_upd); err != nil {
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
	token, err := utils.GenerateAuthenticationToken(user.Id.Hex())
	if err != nil {
		return ok, err
	}
	*data = models.LoginResponse{
		UserData: models.UserData{
			Id:    user.Id,
			Name:  user.Name,
			Email: user.Email,
			Path:  user.Path,
		},
		Token: token,
	}
	return ok, nil
}
