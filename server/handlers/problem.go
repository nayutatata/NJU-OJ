package handlers

import (
	"log"
	"net/http"
	"server/module"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var global_pnumber uint64
var mutex sync.Mutex

func (h *Handler) init_global_number() {
	coll := h.procoll
	findoptions := options.FindOne()
	findoptions.SetSort(bson.D{{Key: "Pnumber", Value: -1}})
	var res module.Problem_t
	err := coll.FindOne(h.ctx, bson.D{}, findoptions).Decode(&res)
	if err == mongo.ErrNoDocuments {
		global_pnumber = 1
	} else {
		global_pnumber = res.Number + 1
	}
}
func (h *Handler) add_problem(problem module.Problem_t) {
	coll := h.procoll
	coll.InsertOne(h.ctx, problem)
}
func (h *Handler) get_problem_by_number(pnumber uint64) (module.Problem_t, error) {
	coll := h.procoll
	var res module.Problem_t
	filter := bson.M{"Pnumber": pnumber}
	err := coll.FindOne(h.ctx, filter).Decode(&res)
	return res, err
}
func (h *Handler) get_all_problems() []module.Problem_t {
	coll := h.procoll
	var res []module.Problem_t = make([]module.Problem_t, 0)
	cursor, _ := coll.Find(h.ctx, bson.M{})
	defer cursor.Close(h.ctx)
	cursor.All(h.ctx, &res)
	return res
}
func (h *Handler) http_get_all_problems(c *gin.Context) {
	problems := h.get_all_problems()
	type prodec_t struct {
		Title   string `json:"title"`
		Pnumber uint64 `json:"Pnumber"`
	}
	pdecs := make([]prodec_t, 0)
	for _, problem := range problems {
		pdecs = append(pdecs, prodec_t{
			Title:   problem.Title,
			Pnumber: problem.Number,
		})
	}
	c.JSON(http.StatusOK, pdecs)
}
func (h *Handler) http_get_problem_by_number(c *gin.Context) {
	get_info := func(c *gin.Context) (uint64, string, error) {
		pnumber := c.Param("Pnumber")
		account := c.Query("account")
		pres, err := strconv.ParseUint(pnumber, 10, 64)
		return pres, account, err
	}

	pnumber, account, err := get_info(c)
	if err != nil {
		c.String(http.StatusBadRequest, "Wrong URL.")
		return
	}
	problem, err := h.get_problem_by_number(pnumber)
	if err != nil {
		c.String(http.StatusBadRequest, "Pnumber does not exist.")
		return
	}
	state := h.get_finish_state(account, pnumber)
	c.JSON(http.StatusOK, gin.H{
		"dec":   problem.Dec,
		"state": state,
	})
}
func (h *Handler) http_add_problem(c *gin.Context) {
	var problem module.Problem_t
	err := c.ShouldBindJSON(&problem)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid problem json.")
	}
	mutex.Lock()
	problem.Number = global_pnumber
	global_pnumber++
	mutex.Unlock()
	h.add_problem(problem)
	c.JSON(http.StatusOK, gin.H{
		"Result":  "Success",
		"Problem": problem,
	})
}
func (h *Handler) http_submit(c *gin.Context) {
	get_info := func(c *gin.Context) (uint64, string, error) {
		pnumber := c.Param("Pnumber")
		account := c.Query("account")
		pres, err := strconv.ParseUint(pnumber, 10, 64)
		return pres, account, err
	}
	pnumber, account, err := get_info(c)
	if err != nil {
		c.String(http.StatusBadRequest, "Wrong URL.")
		return
	}
	var j gin.H
	err = c.ShouldBindJSON(&j)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid json.")
		return
	}
	answer := j["answer"].(string)
	submission := module.Submission_t{
		Account: account,
		Pnumber: pnumber,
		Answer:  answer,
		State:   "NOT",
		SubTime: time.Now(),
		Graded:  false,
		Queuing: false,
	}
	state := h.grade_submission(submission)
	submission.Graded = true
	submission.State = state
	h.add_submission(submission)
	c.JSON(http.StatusOK, gin.H{
		"result": state,
	})
}
func (h *Handler) Init_problems(r *gin.Engine) {
	coll := h.procoll
	indexmodel := mongo.IndexModel{
		Keys: bson.D{
			{
				Key:   "Pnumber",
				Value: 1,
			},
		},
		Options: options.Index().SetUnique(true),
	}
	_, err := coll.Indexes().CreateOne(h.ctx, indexmodel)
	if err != nil {
		log.Fatal(err)
	}
	h.init_global_number()
	group := r.Group("/problems")
	group.GET("", h.http_get_all_problems)
	group.POST("", h.http_add_problem)
	group.GET("/:Pnumber", h.http_get_problem_by_number)
	group.POST("/:Pnumber", h.http_submit)
}
