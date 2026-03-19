package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"showcase_project/data/request/utils"
)

func (h *Handler) CheckLoginUnique(c *gin.Context) {
	var req utils.CheckUniqueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isUnique, appErr := h.service.Utils.CheckLoginUnique(req)
	if appErr != nil {
		c.JSON(appErr.Code(), gin.H{"error": appErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"is_unique": isUnique,
	})
}

func (h *Handler) CheckPhoneUnique(c *gin.Context) {
	var req utils.CheckUniqueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isUnique, appErr := h.service.Utils.CheckPhoneUnique(req)
	if appErr != nil {
		c.JSON(appErr.Code(), gin.H{"error": appErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"is_unique": isUnique,
	})
}
