package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id   primitive.ObjectID `bson:"_id" json:"_id"`
	Name string             `bson:"name" json:"name"`
	Path string             `bson:"path" json:"path"`
}

type Drink struct {
	Id    primitive.ObjectID `bson:"_id" json:"_id"`
	Name  string             `bson:"name" json:"name"`
	Bebeu bool               `bson:"bebeu" json:"bebeu"`
}
