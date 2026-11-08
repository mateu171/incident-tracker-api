package routes

import (
	"incident-tracker-api/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")

	incidents := api.Group("/incidents")
	{
		incidents.POST("", handlers.CreateIncident)
		incidents.GET("", handlers.ListIncidents)
		incidents.GET("/:id", handlers.GetIncident)
		incidents.PATCH("/:id/status", handlers.UpdateIncidentStatus)
		incidents.DELETE("/:id", handlers.DeleteIncident)
	}

	auth := api.Group("/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
	}
}
