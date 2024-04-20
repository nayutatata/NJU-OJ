package handlers

import (
	"log"
	"net/http"
	"server/module"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var coll *mongo.Collection

func (h *Handler) find_user(Uid string) (module.User_t, error) {
	var res module.User_t
	filter := bson.M{"username": Uid}
	err := coll.FindOne(h.ctx, filter).Decode(&res)
	return res, err
}
func (h *Handler) find_and_do(Uid string, c *gin.Context, findfunc func(module.User_t) error, notfindfuc func() error) {
	user, err := h.find_user(Uid)
	if err == nil {
		findfunc(user)
	} else if err == mongo.ErrNoDocuments {
		err = notfindfuc()
	}
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
	}
}
func (h *Handler) create_user(user module.User_t) (module.User_t, error) {
	_, err := coll.InsertOne(h.ctx, user)
	if err == nil {
		return h.find_user(user.Name)
	} else {
		return module.User_t{}, nil
	}
}
func (h *Handler) post_user(c *gin.Context) {
	var user module.User_t
	// 从请求中拿出一个JSON表示要创建的用户的信息
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}
	h.find_and_do(user.Name, c, func(u module.User_t) error {
		c.JSON(http.StatusOK, u)
		return nil
	}, func() error {
		u, err := h.create_user(user)
		if err != nil {
			return err
		}
		c.JSON(http.StatusCreated, u)
		return nil
	})
}
func (h *Handler) get_user(c *gin.Context) {

}
func (h *Handler) Init_user(r *gin.Engine) {
	coll = h.usercoll
	// 用户表项按用户名升序排列，并且用户名是唯一的
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err := coll.Indexes().CreateOne(h.ctx, indexModel)
	if err != nil {
		log.Fatal(err)
	}
	group := r.Group("/user")
	// 下面可以进行一些路由注册
	group.POST("", func(c *gin.Context) {
		c.JSON(http.StatusAccepted, gin.H{
			"name": "not implemented yet",
		})
	})
}
