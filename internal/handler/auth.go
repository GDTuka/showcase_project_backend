package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"showcase_project/data/request/auth"
)

func (h *Handler) Register(c *gin.Context) {
	var req auth.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, appErr := h.service.Auth.Register(req)
	if appErr != nil {
		c.JSON(appErr.Code(), gin.H{"error": appErr.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"id":      id,
	})
}

func (h *Handler) Login(c *gin.Context) {
	var req auth.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, appErr := h.service.Auth.Login(req)
	if appErr != nil {
		c.JSON(appErr.Code(), gin.H{"error": appErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user":    user,
	})
}
