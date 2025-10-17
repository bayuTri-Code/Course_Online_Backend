package handler

import (
	"course_online_backend/internal/dto"
	"course_online_backend/internal/services"
	"course_online_backend/internal/utils"
	"net/http"
	"strings"
	"time"

	_ "course_online_backend/cmd/api/docs"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	Service  *services.AuthService
	Activity *services.ActivityService
}

func NewAuthHandler(service *services.AuthService, activityServices *services.ActivityService ) *AuthHandler {
	return &AuthHandler{Service: service,
	Activity: activityServices,
	}
}

// FirebaseLoginHandler godoc
// @Summary Login user using Firebase ID Token
// @Description This endpoint allows users to log in using a Firebase authentication token (e.g., Google Sign-In).
// The token can be sent either through the Authorization header (Bearer <token>) or in the JSON body.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param Authorization header string false "Bearer Firebase ID Token"
// @Param body body dto.FirebaseLoginRequest false "Firebase token in JSON body (optional)"
// @Success 200 {object} dto.UserLoginResponse "Login successful"
// @Failure 400 {object} utils.ErrorResponse "Invalid request body"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized or invalid Firebase token"
// @Router /api/auth/firebase-login [post]
func (h *AuthHandler) FirebaseLoginHandler(c *gin.Context) {
	var req dto.FirebaseLoginRequest
	authHeader := c.GetHeader("Authorization")

	if authHeader != "" {
		req.Token = strings.Replace(authHeader, "Bearer ", "", 1)
	} else if req.Token != "" {
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing Firebase token"})
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
	go func() {
		_ = h.Activity.LogActivity(user.ID, "User logged in")
	}()
}
