package module

import "go.mongodb.org/mongo-driver/bson/primitive"

type Problem_t struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Number uint64              `bson:"Pnumber"	json:"Pnumber"`
	Dec    string             `bson:"dec"	json:"dec"`
	Title  string             `bson:"title"	json:"title"`
}
