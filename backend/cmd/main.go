package main

import (
	// "net/http"
	"log"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to the server!",
		})
	})
	err := router.Run()
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}