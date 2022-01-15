package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type UpdateUserRequest struct {
	Name           string         `bson:"name,omitempty" json:"name,omitempty"`
	Email          string         `bson:"email,omitempty" json:"email,omitempty"`
	Path           string         `bson:"path,omitempty" json:"path,omitempty"`
	PixAccounts    PixAccounts    `bson:"pixAcc,omitempty" json:"pixAcc,omitempty"`
	WalletAccounts WalletAccounts `bson:"WaletAcc,omitempty" json:"WalletAcc,omitempty"`
}

// Drinks
type CreateDrinkRequest struct {
	UsrId string `bson:"usrId" json:"usrId"`
	Name  string `bson:"name" json:"name"`
}

type UpdateDrinkDoneRequest struct {
	Ids  []string `bson:"_id" json:"_id"`
	Done bool     `bson:"status" json:"status"`
}

// Auth
type AuthenticateUserRequest struct {
	Email    string `bson:"email" json:"email"`
	Password string `bson:"password" json:"password"`
}

// Debts
type CreateDebtRequest struct {
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	Creditor    primitive.ObjectID `bson:"creditor,omitempty" json:"creditor,omitempty"`
	Debtors     []struct {
		Id     primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
		Amount float32            `bson:"amount,omitempty" json:"amount,omitempty"`
	} `bson:"debtors,omitempty" json:"debtors,omitempty"`
	Amount float32 `bson:"amount,omitempty" json:"amount,omitempty"`
}

type PayDebtRequest struct {
	Id       string
	Paid     bool     `json:"paid"`
	Debtors  []string `bson:"debtors,omitempty" json:"debtors,omitempty"`
	Creditor string   `bson:"creditor,omitempty" json:"creditor,omitempty"`
}
