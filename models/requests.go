package models

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
