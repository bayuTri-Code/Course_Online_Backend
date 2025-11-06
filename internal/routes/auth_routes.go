package routes

import (
	"course_online_backend/internal/handler"
	"course_online_backend/internal/services"
	"firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthRoutes(r *gin.RouterGroup, db *gorm.DB, app *firebase.App) {
	authService := services.NewAuthService(db, app)
	activityService := services.NewActivityService(db)
	authHandler := handler.NewAuthHandler(authService, activityService)

	auth := r.Group("/auth")
	auth.POST("/firebase-login", authHandler.FirebaseLoginHandler)
}
