package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id   primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name string             `bson:"name,omitempty" json:"name,omitempty"`
	Path string             `bson:"path,omitempty" json:"path,omitempty"`
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
