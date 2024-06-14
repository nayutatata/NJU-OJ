package module

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	Teacher int = 2
	Student int = 1
)

type User_t struct {
	Uid      primitive.ObjectID `bson:"_id,omitempty"	json:"_id"`
	Name     string             `bson:"username"		json:"username"`
	Level    int                `bson:"level"			json:"level"`
	Account  string             `bson:"account"			json:"account"`
	Password string             `bson:"password"		json:"password"`
}
