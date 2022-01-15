package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type GetUserResponse struct {
	Id             primitive.ObjectID `bson:"_id" json:"_id"`
	Name           string             `bson:"name,omitempty" json:"name,omitempty"`
	Email          string             `bson:"email,omitempty" json:"email,omitempty"`
	Path           string             `bson:"path,omitempty" json:"path,omitempty"`
	PixAccounts    PixAccounts        `bson:"pixAcc,omitempty" json:"pixAcc,omitempty"`
	WalletAccounts WalletAccounts     `bson:"WaletAcc,omitempty" json:"WalletAcc,omitempty"`
	CreatedBy      primitive.ObjectID `bson:"createdby" json:"createdby"`
}

type LoginResponse struct {
	GetUserResponse
	Token string `bson:"token" json:"token"`
}
