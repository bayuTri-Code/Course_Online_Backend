package routes

import (
	"course_online_backend/internal/handler"
	"course_online_backend/internal/middleware"
	"course_online_backend/internal/services"
	"firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserRoutes(r *gin.RouterGroup, db *gorm.DB, app *firebase.App) {
	userService := services.NewUserService(db)
	activityService := services.NewActivityService(db)
	userHandler := handler.NewUserHandler(userService, activityService)

	user := r.Group("/user")
	user.Use(middleware.RoleMiddleware("super_admin"))
	user.GET("/", userHandler.GetAllUsers)
	user.GET("/:id", userHandler.GetUserByID)
	user.PUT("/:id", userHandler.UpdateUser)
	user.DELETE("/:id", userHandler.DeleteUser)
	user.DELETE("/:id/soft", userHandler.SoftDeleteUser)
	user.PATCH("/:id/restore", userHandler.RestoreUser)
}
