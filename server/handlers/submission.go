package handlers

import (
	"log"
	"server/module"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h *Handler) get_finish_state(account string, pnumber uint64) string {
	submissions := h.get_submission(account,pnumber)
	for _, submission := range submissions {
		if submission.State == "AC" {
			return "AC"
		}
	}
	return "NOT"
}
func (h* Handler) get_submission(account string, pnumber uint64) []module.Submission_t {
	coll := h.subcoll
	var res []module.Submission_t = make([]module.Submission_t, 0)
	filter := bson.M{"account":account,"Pnumber":pnumber}
	cursor, _ := coll.Find(h.ctx,filter)
	defer cursor.Close(h.ctx)
	cursor.All(h.ctx,&res)
	return res
}
func (h* Handler) add_submission(submission module.Submission_t) error {
	coll := h.subcoll
	_, err := coll.InsertOne(h.ctx,submission)
	return err
}
func (h *Handler) Init_submission(r *gin.Engine){
	coll := h.subcoll
	indexModel := mongo.IndexModel{
		Keys:bson.D{
			{Key:"time",Value:1},
			{Key:"account",Value: 1},
			{Key: "Pnumber",Value: 1},
		},
	}
	_, err := coll.Indexes().CreateOne(h.ctx,indexModel)
	if err != nil {
		log.Fatal(err)
	}
	//group := r.Group("/submissions")
}