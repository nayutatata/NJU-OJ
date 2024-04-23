package handlers

import (
	"server/module"

	"github.com/gin-gonic/gin"
)
// 从http中找到当前用户的account
func (h *Handler) get_cur_user(c *gin.Context) module.User_t {
	user := module.User_t{
		Name: "admin",
		Account: "0",
		Level: 2,
		Password: "0",
	}
	return user
}
