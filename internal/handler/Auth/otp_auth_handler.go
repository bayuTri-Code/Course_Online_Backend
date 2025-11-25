package authHandler

import (
	"course_online_backend/internal/services"
	otpemail "course_online_backend/internal/services/Auth/otp_email"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type OTPHandler struct {
	OTPService      *otpemail.OTPService
	OTPAuthService  *otpemail.OTPAuthService
	ActivityService *services.ActivityService
}

func NewOTPHandler(otpService *otpemail.OTPService, otpAuthService *otpemail.OTPAuthService, activityService *services.ActivityService) *OTPHandler {
	return &OTPHandler{
		OTPService:      otpService,
		OTPAuthService:  otpAuthService,
		ActivityService: activityService,
	}
}

type SendOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifyOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required,min=4,max=8"`
}

// @Summary Send OTP to email
// @Tags Authentication
// @Accept json
// @Produce json
// @Param body body SendOTPRequest true "Email"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 429 {object} map[string]interface{}
// @Router /auth/otp/send [post]
func (h *OTPHandler) SendOTP(c *gin.Context) {
	var req SendOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid email format"})
		return
	}
	email := strings.ToLower(strings.TrimSpace(req.Email))
	if err := h.OTPService.CheckRateLimit(c.Request.Context(), email); err != nil {
		c.JSON(http.StatusTooManyRequests, gin.H{"status": false, "message": err.Error()})
		return
	}
	if err := h.OTPService.SendOTP(c.Request.Context(), email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to send OTP"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "OTP sent successfully",
		"data": gin.H{
			"email": email,
		},
	})
}

// @Summary Verify OTP and login
// @Tags Authentication
// @Accept json
// @Produce json
// @Param body body VerifyOTPRequest true "Verify OTP"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /auth/otp/verify [post]
func (h *OTPHandler) VerifyOTP(c *gin.Context) {
	var req VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid request"})
		return
	}
	email := strings.ToLower(strings.TrimSpace(req.Email))
	user, token, err := h.OTPAuthService.VerifyOTPAndLogin(c.Request.Context(), email, req.OTP)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
		return
	}
	resp := gin.H{
		"id":         user.ID,
		"firebaseId": user.FirebaseUID,
		"email":      user.EmailAddress,
		"username":   user.Username,
		"roles":      user.Roles,
		"isActive":   user.IsActive,
		"createdAt":  user.CreatedAt,
		"lastLogin":  time.Now(),
		"token":      token,
	}
	c.JSON(http.StatusOK, gin.H{"status": true, "message": "Login success", "data": resp})
	go func() {
		_ = h.ActivityService.LogActivity(user.ID, "Login via OTP")
	}()
}

// @Summary Resend OTP
// @Tags Authentication
// @Accept json
// @Produce json
// @Param body body SendOTPRequest true "Email"
// @Success 200 {object} map[string]interface{}
// @Router /auth/otp/resend [post]
func (h *OTPHandler) ResendOTP(c *gin.Context) {
	h.SendOTP(c)
}
