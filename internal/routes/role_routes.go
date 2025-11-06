package routes

import (
	"course_online_backend/internal/handler"
	"course_online_backend/internal/middleware"
	"course_online_backend/internal/services"
	"firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RoleRoutes(r *gin.RouterGroup, db *gorm.DB, app *firebase.App) {
	userRoleServices := services.NewUserRoleService(db)
	userRoleHandler := handler.NewUserRoleHandler(userRoleServices)

	userRole := r.Group("/admin")
	userRole.Use(middleware.RoleMiddleware("super_admin", "admin"))
	userRole.PUT("/users/:user_id/role", userRoleHandler.AssignRole)
	userRole.GET("/role", userRoleHandler.GetAllRoles)
	userRole.GET("/role/assignable", userRoleHandler.GetAssignableRoles)
	userRole.GET("/role/:id", userRoleHandler.GetRoleByID)
}
