package routes

import (
	"course_online_backend/internal/config"
	authHandler "course_online_backend/internal/handler/Auth"
	"course_online_backend/internal/middleware"
	"course_online_backend/internal/services"
	auth "course_online_backend/internal/services/Auth"
	otpemail "course_online_backend/internal/services/Auth/otp_email"

	"firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func AuthRoutesPublic(r *gin.RouterGroup, db *gorm.DB, app *firebase.App, redisClient *redis.Client) {

	activityService := services.NewActivityService(db)

	otpService := otpemail.NewOTPService(redisClient, config.DbConfig, config.SmtpConfig)
	otpAuthService := otpemail.NewOTPAuthService(db, app, otpService)
	otpHandler := authHandler.NewOTPHandler(otpService, otpAuthService, activityService)

	authService := auth.NewAuthService(db, app)
	authHandler := authHandler.NewAuthHandler(authService, activityService)

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/firebase-login", authHandler.FirebaseLoginHandler)

		authGroup.POST("/otp/send", otpHandler.SendOTP)
		authGroup.POST("/otp/verify", otpHandler.VerifyOTP)
		authGroup.POST("/otp/resend", otpHandler.ResendOTP)
	}
}

func AuthRoutesPrivate(r *gin.RouterGroup, db *gorm.DB, app *firebase.App, redisClient *redis.Client) {
	redisService := auth.NewRedisService(redisClient)
	logoutHandlerInstance := authHandler.NewLogoutHandler(redisService)

	authGroup := r.Group("/auth")
	authGroup.Use(middleware.FirebaseAuth(app, db))
	{
		authGroup.POST("/logout", logoutHandlerInstance.Logout)
	}
}
