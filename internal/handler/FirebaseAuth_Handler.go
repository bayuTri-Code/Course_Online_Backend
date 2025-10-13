package handlers

import (
	"course_online_backend/internal/dto"
	"course_online_backend/internal/services"
	"course_online_backend/internal/utils"
	"net/http"
	"time"

	_ "course_online_backend/cmd/api/docs"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	Service *services.AuthService
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{Service: service}
}

// FirebaseLogin godoc
// @Summary Login with Firebase Token
// @Description Login user with token from Firebase Authentication
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.TokenClaims true "Firebase Token"
// @Success 200 {object} utils.StandardResponse "Login success"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 401 {object} utils.ErrorResponse "Invalid token"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/auth/firebase-login [post]
func (h *AuthHandler) FirebaseLoginHandler(c *gin.Context) {
	var req dto.TokenClaims

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
