package authHandler

import (
	"course_online_backend/internal/services"
	otpemail "course_online_backend/internal/services/Auth/otp_email"
	"log"
	"net/http"
	"strings"

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
// @Description Mengirimkan OTP ke email untuk proses login.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param body body SendOTPRequest true "Email untuk dikirimkan OTP"
// @Success 200 {object} dto.SendOTPResponse "OTP berhasil dikirim"
// @Failure 400 {object} utils.ErrorResponse "Format email tidak valid"
// @Failure 429 {object} utils.ErrorResponse "Rate limit tercapai"
// @Failure 500 {object} utils.ErrorResponse "Gagal mengirim OTP"
// @Router /api/auth/otp/send [post]
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
		log.Printf("[OTP_HANDLER] Failed to send OTP to %s: %v", email, err)
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
// @Success 200 {object} dto.VerifyOTPResponse "Login Berhasil"
// @Failure 400 {object} utils.ErrorResponse "OTP tidak valid atau telah kedaluwarsa"
// @Failure 429 {object} utils.ErrorResponse "Rate limit tercapai"
// @Failure 500 {object} utils.ErrorResponse "Terjadi kesalahan pada server"
// @Router /api/auth/otp/verify [post]
func (h *OTPHandler) VerifyOTP(c *gin.Context) {
	var req VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid request"})
		return
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))

	user, customToken, err := h.OTPAuthService.VerifyOTPAndLogin(c.Request.Context(), email, req.OTP)
	if err != nil {
		log.Printf("[OTP_HANDLER] Failed to verify OTP for %s: %v", email, err)
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Login success",
		"data": gin.H{
			"customToken": customToken,
			"user": gin.H{             
				"id":         user.ID,
				"firebaseId": user.FirebaseUID,
				"email":      user.EmailAddress,
				"username":   user.Username,
				"roles":      user.Roles,
				"isActive":   user.IsActive,
				"createdAt":  user.CreatedAt,
				"lastLogin":  user.LastLogin,
			},
		},
	})

	go func() {
		if err := h.ActivityService.LogActivity(user.ID, "Login via OTP"); err != nil {
			log.Printf("[OTP_HANDLER] Failed to log activity for user %s: %v", user.ID, err)
		}
	}()
}

// @Summary Resend OTP
// @Tags Authentication
// @Accept json
// @Produce json
// @Param body body SendOTPRequest true "Email"
// @Success 200 {object} map[string]interface{}
// @Router /api/auth/otp/resend [post]
func (h *OTPHandler) ResendOTP(c *gin.Context) {
	h.SendOTP(c)
}