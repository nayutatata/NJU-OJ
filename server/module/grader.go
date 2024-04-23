package module

import "go.mongodb.org/mongo-driver/bson/primitive"

type Grader_t struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"	json:"id"`
	Pnumber uint64             `bson:"Pnumber"	json:"Pnumber"`
	Inputs  []string           `bson:"inputs"	json:"inputs"`
	Outputs []string           `bson:"outputs"	json:"outputs"`
}
