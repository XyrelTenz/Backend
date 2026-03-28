package handlers

import (
	"backend/pkg/firebase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	firebaseApp *firebase.App
}

func NewUserHandler(firebaseApp *firebase.App) *UserHandler {
	return &UserHandler{
		firebaseApp: firebaseApp,
	}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"user_id": userID.(string),
	})
}
