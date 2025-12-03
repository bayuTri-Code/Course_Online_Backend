package Course_management_routes

import (
	ModuleMgmthandler "course_online_backend/internal/handler/Course_Management/Module_Management_Handler"
	"course_online_backend/internal/middleware"
	"course_online_backend/internal/services"
	modulemanagementServices "course_online_backend/internal/services/Course_management_Services/Module_Management_Services"

	firebase "firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ModuleManagementRoutes(r *gin.RouterGroup, db *gorm.DB, app *firebase.App) {
	activityService := services.NewActivityService(db)
	moduleServices := modulemanagementServices.NewModuleService(db)
	moduleHandler := ModuleMgmthandler.NewModuleHandler(moduleServices, activityService)

	publicCourseModules := r.Group("/courses/:course_id/modules")
	publicCourseModules.Use(middleware.FirebaseAuth(app, db))
	publicCourseModules.Use(middleware.RoleMiddleware("student", "admin", "super_admin", "instructor"))
	publicCourseModules.Use(middleware.CheckEnrollmentAccess(db))
	{
		publicCourseModules.GET("", moduleHandler.GetModulesByCourseHandler)
	}

	publicModules := r.Group("/modules")
	{
		publicModules.GET("/:module_id", moduleHandler.GetModuleByIDHandler)
	}

	securedCourseModules := r.Group("/courses/:course_id/modules")
	securedCourseModules.Use(middleware.FirebaseAuth(app, db))
	securedCourseModules.Use(middleware.RoleMiddleware("admin", "super_admin"))
	securedCourseModules.Use(middleware.CheckCourseOwnershipDynamic(db))
	{
		securedCourseModules.POST("", moduleHandler.CreateModuleHandler)
	}

	securedModules := r.Group("/modules")
	securedModules.Use(middleware.FirebaseAuth(app, db))
	securedModules.Use(middleware.RoleMiddleware("admin", "super_admin"))
	securedModules.Use(middleware.CheckCourseOwnershipDynamic(db))
	{
		securedModules.PUT("/:module_id", moduleHandler.UpdateModuleHandler)
		securedModules.DELETE("/:module_id", moduleHandler.SoftDeleteModuleHandler)
		securedModules.PATCH("/:module_id/restore", moduleHandler.RestoreModuleHandler)
		securedModules.DELETE("/:module_id/permanent", moduleHandler.PermanentDeleteModuleHandler)
	}
}
