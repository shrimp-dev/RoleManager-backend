package models

type UpdateUserRequest struct {
	Name           string         `bson:"name,omitempty" json:"name,omitempty"`
	Email          string         `bson:"email,omitempty" json:"email,omitempty"`
	Path           string         `bson:"path,omitempty" json:"path,omitempty"`
	PixAccounts    PixAccounts    `bson:"pixAcc,omitempty" json:"pixAcc,omitempty"`
	WalletAccounts WalletAccounts `bson:"WaletAcc,omitempty" json:"WalletAcc,omitempty"`
}
