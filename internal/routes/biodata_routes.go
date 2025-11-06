package routes

import (
	"course_online_backend/internal/handler"
	"course_online_backend/internal/middleware"
	"course_online_backend/internal/services"
	"firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func BiodataRoutes(r *gin.RouterGroup, db *gorm.DB, app *firebase.App) {
	biodataService := services.NewBiodataService(db)
	activityService := services.NewActivityService(db)
	biodataHandler := handler.NewBiodataHandler(biodataService, activityService)

	biodata := r.Group("/profile")
	biodata.Use(middleware.RoleMiddleware("student", "super_admin", "admin", "instructor"))
	biodata.POST("/biodata", biodataHandler.CreateBiodata)
	biodata.GET("/mybiodata", biodataHandler.GetBiodata)
	biodata.PUT("/biodata", biodataHandler.UpdateBiodata)
	biodata.DELETE("/biodata", biodataHandler.DeleteBiodata)
	biodata.DELETE("/biodata/soft-delete", biodataHandler.SoftDeleteBiodata)
	biodata.PUT("/biodata/restore", biodataHandler.RestoreBiodata)
}
