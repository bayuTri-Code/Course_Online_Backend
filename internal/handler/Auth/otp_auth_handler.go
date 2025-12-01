package authHandler

import (
	"course_online_backend/internal/dto"
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

// @Summary Send OTP to email
// @Description Mengirimkan OTP ke email untuk proses login.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param body body dto.SendOTPRequest true "Email untuk dikirimkan OTP"
// @Success 200 {object} dto.SendOTPResponse "OTP berhasil dikirim"
// @Failure 400 {object} utils.ErrorResponse "Format email tidak valid"
// @Failure 429 {object} utils.ErrorResponse "Rate limit tercapai"
// @Failure 500 {object} utils.ErrorResponse "Gagal mengirim OTP"
// @Router /api/auth/otp/send [post]
func (h *OTPHandler) SendOTP(c *gin.Context) {
	var req dto.SendOTPRequest
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

	c.JSON(http.StatusOK, dto.SendOTPResponse{
		Status:  true,
		Message: "OTP sent successfully",
		Data: dto.SendOTPData{
			Email: email,
		},
	})

}

// @Summary Verify OTP and login
// @Tags Authentication
// @Accept json
// @Produce json
// @Param body body dto.VerifyOTPRequest true "Verify OTP"
// @Success 200 {object} dto.VerifyOTPResponse "Login Berhasil"
// @Failure 400 {object} utils.ErrorResponse "OTP tidak valid atau telah kedaluwarsa"
// @Failure 429 {object} utils.ErrorResponse "Rate limit tercapai"
// @Failure 500 {object} utils.ErrorResponse "Terjadi kesalahan pada server"
// @Router /api/auth/otp/verify [post]
func (h *OTPHandler) VerifyOTP(c *gin.Context) {
	var req dto.VerifyOTPRequest
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

	c.JSON(http.StatusOK, dto.VerifyOTPResponse{
		Status:  true,
		Message: "Login success",
		Data: dto.VerifyOTPData{
			CustomToken: customToken,
			User: dto.OTPUserResponse{
				ID:         user.ID.String(),
				FirebaseID: user.FirebaseUID,
				Email:      user.EmailAddress,
				Username:   user.Username,
				IsActive:   user.IsActive,
				CreatedAt:  user.CreatedAt,
				LastLogin:  user.LastLogin,
				Roles: func() []dto.UserRoleResponse {
					roles := make([]dto.UserRoleResponse, 0)
					for _, r := range user.Roles {
						roles = append(roles, dto.UserRoleResponse{
							ID:   r.ID.String(),
							Name: r.Name,
						})
					}
					return roles
				}(),
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
// @Param body body dto.SendOTPRequest true "Email"
// @Success 200 {object} dto.SendOTPResponse
// @Router /api/auth/otp/resend [post]
func (h *OTPHandler) ResendOTP(c *gin.Context) {
	h.SendOTP(c)
}
