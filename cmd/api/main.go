package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	port = 8080
)

func main() {
	// basic gin web api
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World",
		})
	})
	r.Run(fmt.Sprintf(":%d", port))
}
