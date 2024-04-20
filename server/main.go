package main

import (
	"context"
	"server/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	h := handlers.GetHandler(context.Background())
	r := gin.Default()
	h.Init_user(r)
	r.Run()
}
