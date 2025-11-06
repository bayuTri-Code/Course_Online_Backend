package routes

import (
	"course_online_backend/internal/handler"
	"course_online_backend/internal/middleware"
	"course_online_backend/internal/services"
	"firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ActivityRoutes(r *gin.RouterGroup, db *gorm.DB, app *firebase.App) {
	activityService := services.NewActivityService(db)
	activityHandler := handler.NewActivityHandler(activityService)

	activity := r.Group("/history")
	activity.Use(middleware.RoleMiddleware("admin", "super_admin"))
	activity.GET("/All-activity", activityHandler.GetAllActivity)
	activity.GET("/ByUser-activity/:id", activityHandler.GetActivityByUserId)
	activity.GET("/Recent-activity", activityHandler.GetRecentActivities)
	activity.GET("/Search-activity", activityHandler.SearchActivity)
	activity.GET("/Summary-activity/:id", activityHandler.GetActivitySummary)
	activity.GET("/ByRole-activity/:role", activityHandler.GetActivityByRole)
}
