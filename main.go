package main

import (
	"incident-tracker-api/config"
	"incident-tracker-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDatabase()

	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	routes.SetupRoutes(router)

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
