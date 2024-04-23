package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func not_implement(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"return-type": "not implement yet.",
	})
}
