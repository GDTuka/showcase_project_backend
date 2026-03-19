package handler

import (
	"net/http"

	"showcase_project/data/request/auth"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Register(c *gin.Context) {
	var req auth.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, at, rt, appErr := h.service.Auth.Register(req)
	if appErr != nil {
		c.JSON(appErr.Code(), gin.H{"error": appErr.Error()})
		return
	}

	c.Header("X-Access-Token", at.Token)
	c.Header("X-Refresh-Token", rt.Token)

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

	user, at, rt, appErr := h.service.Auth.Login(req)
	if appErr != nil {
		c.JSON(appErr.Code(), gin.H{"error": appErr.Error()})
		return
	}

	c.Header("X-Access-Token", at.Token)
	c.Header("X-Refresh-Token", rt.Token)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user":    user,
	})
}

func (h *Handler) RefreshToken(c *gin.Context) {
	var req auth.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	at, rt, appErr := h.service.Auth.RefreshToken(req.RefreshToken)
	if appErr != nil {
		c.JSON(appErr.Code(), gin.H{"error": appErr.Error()})
		return
	}

	c.Header("X-Access-Token", at.Token)
	c.Header("X-Refresh-Token", rt.Token)

	c.JSON(http.StatusOK, gin.H{
		"message": "Tokens refreshed successfully",
	})
}


