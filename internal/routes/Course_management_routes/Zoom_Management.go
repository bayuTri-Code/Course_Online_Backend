package Course_management_routes

import (
	zoommanagement "course_online_backend/internal/handler/Course_Management/Zoom_Management"
	"course_online_backend/internal/middleware"
	"course_online_backend/internal/services"
	zoommanagementservices "course_online_backend/internal/services/Course_management_Services/Zoom_Management_Services"

	firebase "firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ZoomRoutes(r *gin.RouterGroup, db *gorm.DB, app *firebase.App) {
	zoomService := zoommanagementservices.NewZoomService(db)
	activityService := services.NewActivityService(db)
	zoomHandler := zoommanagement.NewZoomHandler(zoomService, activityService)

	zoom := r.Group("/zoom")
	zoom.Use(middleware.RoleMiddleware("super_admin", "admin", "instructor"))
	{
		zoom.POST("", zoomHandler.CreateZoomHandler)
		zoom.GET("", zoomHandler.GetAllZoomsHandler)
		zoom.GET("/:id", zoomHandler.GetZoomByIDHandler)
		zoom.PUT("/:id", zoomHandler.UpdateZoomHandler)
		zoom.DELETE("/:id", zoomHandler.DeleteZoomHandler)
		
		zoom.GET("/course/:course_id", zoomHandler.GetZoomsByCourseIDHandler)
		zoom.GET("/course/:course_id/upcoming", zoomHandler.GetUpcomingZoomsByCourseIDHandler)
	}
}