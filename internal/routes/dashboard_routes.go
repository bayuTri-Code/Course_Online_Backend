package routes

import (
	"course_online_backend/internal/handler/dashboard"
	"course_online_backend/internal/middleware"
	"course_online_backend/internal/services"
	"firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func DashboardRoutes(r *gin.RouterGroup, db *gorm.DB, app *firebase.App) {
	dashboardService := services.NewDashboardService(db)

	superAdminHandler := dashboard.NewSuperAdminDashboardHandler(dashboardService)
	adminHandler := dashboard.NewAdminDashboardHandler(dashboardService)
	instructorHandler := dashboard.NewInstructorDashboardHandler(dashboardService)
	studentHandler := dashboard.NewStudentDashboardHandler(dashboardService)

	dashboardGroup := r.Group("/dashboard")
	dashboardGroup.GET("/super_admin", middleware.RoleMiddleware("super_admin"), superAdminHandler.GetDashboardSuperAdmin)
	dashboardGroup.GET("/admin", middleware.RoleMiddleware("admin"), adminHandler.GetDashboardAdmin)
	dashboardGroup.GET("/instructor", middleware.RoleMiddleware("instructor"), instructorHandler.GetDashboardInstructor)
	dashboardGroup.GET("/student", middleware.RoleMiddleware("student"), studentHandler.GetDashboardStudent)
}
