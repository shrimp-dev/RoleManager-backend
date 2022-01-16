package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id             primitive.ObjectID `bson:"_id" json:"_id"`
	Name           string             `bson:"name,omitempty" json:"name,omitempty"`
	Email          string             `bson:"email,omitempty" json:"email,omitempty"`
	Path           string             `bson:"path,omitempty" json:"path,omitempty"`
	PixAccounts    PixAccounts        `bson:"pixAcc,omitempty" json:"pixAcc,omitempty"`
	WalletAccounts WalletAccounts     `bson:"WaletAcc,omitempty" json:"WalletAcc,omitempty"`
	CreatedBy      primitive.ObjectID `bson:"createdby" json:"createdby"`
	Password       string             `bson:"password" json:"password"`
	Salt           []byte             `bson:"salt" json:"salt"`
}

type PixAccounts struct {
	PixPreferred uint     `bson:"pixPref,omitempty" json:"pixPref,omitempty"`
	PixAccounts  []string `bson:"pixAccs,omitempty" json:"pixAccs,omitempty"`
}

type WalletAccounts struct {
	WalletPreferred uint     `bson:"walletPref,omitempty" json:"walletPref,omitempty"`
	WalletAccounts  []string `bson:"walletAccs,omitempty" json:"WalletAccs,omitempty"`
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

type AccessTokenClaims struct {
	Id  string    `bson:"_id" json:"_id"`
	Exp time.Time `bson:"exp" json:"exp"`
}
