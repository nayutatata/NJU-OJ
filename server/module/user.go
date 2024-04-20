package module

import "go.mongodb.org/mongo-driver/bson/primitive"
const(
	teacher int = 1
	student int = 2
)
type User_t struct {
	Uid primitive.ObjectID  `bson:"_id,omitempty"`
	Name string				`bson:"username"`
	Level int				`bson:"level"`
}

