package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {

	c := gin.Default()
	c.Use(func(context *gin.Context) {
		fmt.Println("中间件1")
	}, func(ctx *gin.Context) {
		fmt.Println("中间件2")
	})

	c.GET("/middle", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"code": "200"})
	})

	c.Run(":8000")

}
