package handler

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	// Auth routes
	auth := router.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.POST("/refresh", h.RefreshToken)
	}

	// Utils routes
	utils := router.Group("/utils")
	{
		utils.POST("/login-unique", h.CheckLoginUnique)
		utils.POST("/phone-unique", h.CheckPhoneUnique)
	}

	// Simple ping route
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	return router
}
