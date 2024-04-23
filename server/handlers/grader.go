package handlers

import (
	"log"
	"net/http"
	"server/judger"
	"server/module"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func run_grader(grader module.Grader_t, submission module.Submission_t) string {
	return judger.Judge_samples(submission.Answer, grader.Inputs, grader.Outputs)
}
func (h *Handler) grade_submission(submission module.Submission_t) string {
	grader, err := h.find_grader(submission.Pnumber)
	if err != nil && err == mongo.ErrNoDocuments {
		return "The Grader does not exist."
	} else if err != nil {
		log.Fatal(err)
	}
	return run_grader(grader, submission)
}
func (h *Handler) find_grader(pnumber uint64) (module.Grader_t, error) {
	coll := h.gracoll
	var grader module.Grader_t
	filter := bson.M{
		"Pnumber": pnumber,
	}
	err := coll.FindOne(h.ctx, filter).Decode(&grader)
	return grader, err
}
func (h *Handler) add_grader(grader module.Grader_t) error {
	coll := h.gracoll
	_, err := coll.InsertOne(h.ctx, grader)
	return err
}
func (h *Handler) update_grader(grader module.Grader_t) error {
	coll := h.gracoll
	_, err := h.find_grader(grader.Pnumber)
	if err != nil && err == mongo.ErrNoDocuments {
		return h.add_grader(grader)
	}
	filter := bson.M{"Pnumber": grader.Pnumber}
	_, err = coll.UpdateOne(h.ctx, filter, grader)
	return err
}
func (h *Handler)http_add_grader(c *gin.Context) {
	get_basic := func(c *gin.Context) (uint64,error) {
		a := c.Param("Pnumber")
		pnumber, err := strconv.ParseUint(a, 10, 64)
		return pnumber,err
	}
	Pnumber, err := get_basic(c)
	if err != nil {
		c.String(http.StatusBadRequest,"Invalid Pnumber.")
		return
	}
	var grader module.Grader_t
	err = c.ShouldBindJSON(&grader)
	if err != nil {
		c.String(http.StatusBadRequest,"Invalid information of a grader.")
		return
	}
	grader.Pnumber = Pnumber
	err = h.add_grader(grader)
	if err != nil {
		c.String(http.StatusInternalServerError,"Data Base seems to raise an error.")
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"result":"success",
	})
}
func (h *Handler) Init_grader(r *gin.Engine) {
	coll := h.gracoll
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{
				Key:   "Pnumber",
				Value: 1,
			},
		},
		Options: options.Index().SetUnique(true),
	}
	_, err := coll.Indexes().CreateOne(h.ctx, indexModel)
	if err != nil {
		log.Fatal(err)
	}
	h.add_grader(
		module.Grader_t{
			Pnumber: 1,
			Inputs:  []string{"100", "200", "300"},
			Outputs: []string{"101", "201", "300"},
		},
	)
	group := r.Group("grader")
	group.POST("/:Pnumber", h.http_add_grader)
}
