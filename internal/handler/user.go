package handler

import (
	"net/http"
	"strconv"

	"showcase_project/data/request/user"

	"github.com/gin-gonic/gin"
)

func (h *Handler) SendSmsCode(c *gin.Context) {
	var req user.SendSmsCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.User.SendSmsCode(req.Phone); err != nil {
		c.JSON(err.Code(), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "SMS code sent successfully"})
}

func (h *Handler) GetCurrentUser(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userWithProfile, err := h.service.User.GetCurrentUser(userId.(int))
	if err != nil {
		c.JSON(err.Code(), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userWithProfile)
}

func (h *Handler) GetUserById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	userWithProfile, appErr := h.service.User.GetUserById(id)
	if appErr != nil {
		c.JSON(appErr.Code(), gin.H{"error": appErr.Error()})
		return
	}

	c.JSON(http.StatusOK, userWithProfile)
}

func (h *Handler) SearchUsers(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req user.SearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set defaults if not provided
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	users, appErr := h.service.User.SearchUsers(req, userId.(int))
	if appErr != nil {
		c.JSON(appErr.Code(), gin.H{"error": appErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"limit": req.Limit,
		"offset": req.Offset,
	})
}
