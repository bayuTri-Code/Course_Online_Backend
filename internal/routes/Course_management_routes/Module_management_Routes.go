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

	courseModules := r.Group("/courses/:id/modules")
	courseModules.Use(middleware.RoleMiddleware("super_admin", "admin", "instructor"))
	{
		courseModules.GET("", moduleHandler.GetModulesByCourseHandler)
		courseModules.POST("", moduleHandler.CreateModuleHandler)
	}

	modules := r.Group("/modules")
	modules.Use(middleware.RoleMiddleware("super_admin", "admin", "instructor"))
	{
		modules.GET("/:id", moduleHandler.GetModuleByIDHandler)
		modules.PUT("/:id", moduleHandler.UpdateModuleHandler)
		modules.DELETE("/:id", moduleHandler.SoftDeleteModuleHandler)
		modules.PATCH("/:id/restore", moduleHandler.RestoreModuleHandler)
		modules.DELETE("/:id/permanent", moduleHandler.PermanentDeleteModuleHandler)
	}
}
