package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id       primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name     string             `bson:"name,omitempty" json:"name,omitempty"`
	Email    string             `bson:"email" json:"email"`
	Path     string             `bson:"path,omitempty" json:"path,omitempty"`
	Password string             `bson:"password" json:"password"`
	Salt     []byte             `bson:"salt" json:"salt"`
}

type UserData struct {
	Id    primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name  string             `bson:"name,omitempty" json:"name,omitempty"`
	Email string             `bson:"email" json:"email"`
	Path  string             `bson:"path,omitempty" json:"path,omitempty"`
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
	Id          primitive.ObjectID   `bson:"_id,omitempty" json:"_id,omitempty"`
	UsrId       primitive.ObjectID   `bson:"usrId,omitempty" json:"usrId,omitempty"`
	Name        string               `bson:"name,omitempty" json:"name,omitempty"`
	Email       string               `bson:"email" json:"email"`
	Password    string               `bson:"password" json:"password"`
	Ids         []primitive.ObjectID `bson:"ids,omitempty" json:"ids,omitempty"`
	Done        bool                 `bson:"done,omitempty" json:"done,omitempty"`
	Description string               `bson:"description,omitempty" json:"description,omitempty"`
	Creditor    primitive.ObjectID   `bson:"creditor,omitempty" json:"creditor,omitempty"`
	Amount      float32              `bson:"amount,omitempty" json:"amount,omitempty"`
	Paid        bool                 `bson:"paid,omitempty" json:"paid,omitempty"`

	Debtors []struct {
		Id     primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
		Amount float32            `bson:"amount,omitempty" json:"amount,omitempty"`
		Paid   bool               `bson:"paid,omitempty" json:"paid,omitempty"`
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
