package authHandler

import (
	"course_online_backend/internal/models"
	auth "course_online_backend/internal/services/Auth"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type LogoutHandler struct {
	logoutService *auth.RedisService
}

func NewLogoutHandler(logoutService *auth.RedisService) *LogoutHandler {
	return &LogoutHandler{
		logoutService: logoutService,
	}
}

// Logout godoc
// @Summary User logout
// @Description Logout user dan blacklist token saat ini
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/auth/logout [post]
func (h *LogoutHandler) Logout(c *gin.Context) {
	tokenInterface, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Token not found in context",
		})
		return
	}

	token, ok := tokenInterface.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid token format",
		})
		return
	}

	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not found in context",
		})
		return
	}

	user, ok := userInterface.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid user format",
		})
		return
	}

	ctx := c.Request.Context()

	err := h.logoutService.BlacklistToken(ctx, token, 24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to logout",
			"details": err.Error(),
		})
		return
	}

	_ = h.logoutService.DeleteTokenCache(ctx, token)

	_ = h.logoutService.DeleteRefreshToken(ctx, user.FirebaseUID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
		"data": gin.H{
			"logged_out_at": time.Now(),
			"user_id": user.FirebaseUID,
		},
	})
}
