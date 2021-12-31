package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id   primitive.ObjectID `bson:"_id" json:"_id"`
	Name string             `bson:"name" json:"name"`
	Path string             `bson:"path" json:"path"`
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
	Amount      float32            `bson:"amount" json:"amount"`
	Open        bool               `bson:"open" json:"open"`
}

type Request struct {
	Id          string   `bson:"_id" json:"_id"`
	UsrId       string   `bson:"usrId" json:"usrId"`
	Name        string   `bson:"name" json:"name"`
	Ids         []string `bson:"ids" json:"ids"`
	Done        bool     `bson:"done" json:"done"`
	Description string   `bson:"description" json:"description"`
	Creditor    string   `bson:"creditor" json:"creditor"`
	Amount      float32  `bson:"amount" json:"amount"`

	Debtors []struct {
		Id     string  `bson:"_id" json:"_id"`
		Amount float32 `bson:"amount" json:"amount"`
		Paid   bool    `bson:"paid" json:"paid"`
	} `bson:"debtors" json:"debtors"`
}
