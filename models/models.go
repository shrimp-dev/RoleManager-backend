package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	UserData
	Password  string             `bson:"password" json:"password"`
	Salt      []byte             `bson:"salt" json:"salt"`
	CreatedBy primitive.ObjectID `bson:"createdby" json:"createdby"`
}

type UserData struct {
	UserUpdate
	Id        primitive.ObjectID `bson:"_id" json:"_id"`
	CreatedBy primitive.ObjectID `bson:"createdby" json:"createdby"`
}

type UserUpdate struct {
	Name  string `bson:"name,omitempty" json:"name,omitempty"`
	Email string `bson:"email,omitempty" json:"email,omitempty"`
	Path  string `bson:"path,omitempty" json:"path,omitempty"`
}

type Drink struct {
	Id    primitive.ObjectID `bson:"_id" json:"_id"`
	UsrId primitive.ObjectID `bson:"usrId" json:"usrId"`
	Name  string             `bson:"name" json:"name"`
	Done  bool               `bson:"done" json:"done"`
}

type Debtor struct {
	Id     primitive.ObjectID `bson:"_id" json:"_id"`
	Amount float32            `bson:"amount" json:"amount"`
	Paid   bool               `bson:"paid" json:"paid"`
}

type Debt struct {
	Id          primitive.ObjectID `bson:"_id" json:"_id"`
	Description string             `bson:"description" json:"description"`
	Creditor    primitive.ObjectID `bson:"creditor" json:"creditor"`
	Debtors     []Debtor           `bson:"debtors" json:"debtors"`
}

type Request struct {
	Id          primitive.ObjectID   `bson:"_id" json:"_id"`
	UsrId       primitive.ObjectID   `bson:"usrId" json:"usrId"`
	Name        string               `bson:"name" json:"name"`
	Email       string               `bson:"email" json:"email"`
	Password    string               `bson:"password" json:"password"`
	Ids         []primitive.ObjectID `bson:"ids" json:"ids"`
	Done        bool                 `bson:"done" json:"done"`
	Description string               `bson:"description" json:"description"`
	Creditor    primitive.ObjectID   `bson:"creditor" json:"creditor"`
	Amount      float32              `bson:"amount" json:"amount"`
	Paid        bool                 `bson:"paid" json:"paid"`

	Debtors []struct {
		Id     primitive.ObjectID `bson:"_id" json:"_id"`
		Amount float32            `bson:"amount" json:"amount"`
		Paid   bool               `bson:"paid" json:"paid"`
	} `bson:"debtors" json:"debtors"`
}

type LoginResponse struct {
	UserData
	Token string `bson:"token" json:"token"`
}

type AccessTokenClaims struct {
	Id  string    `bson:"_id" json:"_id"`
	Exp time.Time `bson:"exp" json:"exp"`
}
