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

func get_created_user(c *gin.Context) (module.User_t, error) {
	var info gin.H
	err := c.ShouldBindJSON(&info)
	var user module.User_t
	user.Account, _ = info["Account"].(string)
	user.Name, _ = info["NickName"].(string)
	user.Password, _ = info["Password"].(string)
	user.Level = 1
	return user, err
}
func get_login_info(c *gin.Context) (string, string) {
	var info gin.H
	err := c.ShouldBindJSON(&info)
	if err != nil {
		log.Fatal(err)
	}
	return info["Account"].(string), info["Password"].(string)
}
func (h *Handler) find_user(account string) (module.User_t, error) {
	coll := h.usercoll
	var res module.User_t
	filter := bson.M{"account": account}
	err := coll.FindOne(h.ctx, filter).Decode(&res)
	return res, err
}
func (h *Handler) find_and_do(account string, c *gin.Context, findfunc func(module.User_t) error, notfindfuc func() error) {
	user, err := h.find_user(account)
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
	coll := h.usercoll
	_, err := coll.InsertOne(h.ctx, user)
	if err == nil {
		return user, nil
	} else {
		return module.User_t{}, err
	}
}
func (h *Handler) http_add_user(c *gin.Context) {
	user, err := get_created_user(c)
	// 从请求中拿出一个JSON表示要创建的用户的信息
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid user json.")
		return
	}
	h.find_and_do(user.Account, c, func(u module.User_t) error {
		c.JSON(http.StatusOK, gin.H{
			"Result": "Fail",
			"Reason": "Account has existed",
		})
		return nil
	}, func() error {
		_, err := h.create_user(user)
		if err != nil && err != mongo.ErrNoDocuments {
			// 数据库出错
			c.String(http.StatusInternalServerError, "Data Base seems to raise an error.")
			return err
		}
		c.JSON(http.StatusOK, gin.H{
			"Result": "Success",
			"Reason": "null",
		})
		return nil
	})
}
func (h *Handler) http_login(c *gin.Context) {
	accout, password := get_login_info(c)
	h.find_and_do(accout, c, func(u module.User_t) error {
		if password == u.Password {
			c.JSON(http.StatusOK, gin.H{
				"Result": "Success",
				"Reason": "null",
				"Info": gin.H{
					"account":  u.Account,
					"NickName": u.Name,
					"figure":   u.Level,
				},
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"Result": "Fail",
				"Reason": "Incorrect password",
				"Info":   "null",
			})
		}
		return nil
	}, func() error {
		c.JSON(http.StatusOK, gin.H{
			"Result": "Fail",
			"Reason": "Account not found.",
			"Info":   "null",
		})
		return nil
	})
}
func (h *Handler) Init_user(r *gin.Engine) {
	coll := h.usercoll
	// 用户表项按用户账号升序排列，并且账号是唯一的
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "account", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err := coll.Indexes().CreateOne(h.ctx, indexModel)
	if err != nil {
		log.Fatal(err)
	}
	h.create_user(module.User_t{
		Account:  "0",
		Name:     "admin",
		Password: "123",
		Level:    module.Teacher,
	})
	group := r.Group("/user")
	// 下面可以进行一些路由注册
	group.POST("/register", h.http_add_user)
	group.POST("/login", h.http_login)
}
