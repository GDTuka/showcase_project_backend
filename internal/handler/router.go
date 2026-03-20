package handler

import (
	"showcase_project/internal/middleware"

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

	// SMS route (Public)
	router.POST("/user/sms/send", h.SendSmsCode)

	// User routes (Protected)
	userRoutes := router.Group("/user", middleware.AuthMiddleware(h.service.JWT))
	{
		userRoutes.GET("/me", h.GetCurrentUser)
		userRoutes.GET("/search", h.SearchUsers)
		userRoutes.GET("/:id", h.GetUserById)
	}

	// Simple ping route
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	return router
}
