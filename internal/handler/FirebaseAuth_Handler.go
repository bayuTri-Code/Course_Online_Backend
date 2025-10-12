package handlers

import (
	"course_online_backend/internal/services"
	"course_online_backend/internal/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	Service *services.AuthService
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{Service: service}
}

func (h *AuthHandler) FirebaseLoginHandler(c *gin.Context) {
	var req struct {
		Token string `json:"token"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	user, err := h.Service.LoginWithFirebaseToken(req.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	userResponse := gin.H{
		"id":         user.ID,
		"firebaseId": user.FirebaseUID,
		"email":      user.EmailAddress,
		"username":   user.Username,
		"roles":      user.Roles,
		"isActive":   user.IsActive,
		"createdAt":  user.CreatedAt,
		"lastLogin":  time.Now(),
	}

	utils.JSONSuccess(c, userResponse, "Login success")
}
