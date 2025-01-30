package main

import (
	"fmt"
	"go-mma/config"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	port = 8080
)

func main() {
	config, err := config.Load()
	if err != nil {
		log.Panic(err)
	}
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World",
		})
	})
	r.Run(fmt.Sprintf(":%d", config.HTTPPort))
}
