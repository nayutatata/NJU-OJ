package module

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Submission_t struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"	json:"id"`
	Account string             `bson:"account"	json:"account"`
	Pnumber uint64             `bson:"Pnumber"	json:"Pnumber"`
	State   string             `bson:"state"	json:"state"`
	SubTime time.Time          `bson:"time"	json:"time"`
	Answer  string             `bson:"answer"	json:"answer"`
	Graded  bool               `bson:"graded"	json:"graded"`
	Queuing bool               `bson:"queuing"	json:"queuing"`
}
