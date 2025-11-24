package Course_management_routes


import (
	CourseManagementhandler "course_online_backend/internal/handler/Course_Management"
	"course_online_backend/internal/middleware"
	"course_online_backend/internal/services"
	CourseManagementServices "course_online_backend/internal/services/Course_management_Services"

	firebase "firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


func CourseTypePublic(r *gin.RouterGroup, db *gorm.DB, app *firebase.App) {
	activityService := services.NewActivityService(db)
	courseTypeServicesPublic := CourseManagementServices.NewCourseTypeService(db)
	courseTypeHandlerPublic := CourseManagementhandler.NewCourseTypeHandler(courseTypeServicesPublic, activityService)

	courseTypes := r.Group("/course-types")
	{
		courseTypes.GET("", courseTypeHandlerPublic.GetAllCourseTypeHandler)
	}
}

func CourseTypeRoutes(r *gin.RouterGroup, db *gorm.DB, app *firebase.App) {
	activityService := services.NewActivityService(db)
	courseTypeServices := CourseManagementServices.NewCourseTypeService(db)
	courseTypeHandler := CourseManagementhandler.NewCourseTypeHandler(courseTypeServices, activityService)

	courseTypes := r.Group("/course-types")
	courseTypes.Use(middleware.RoleMiddleware("super_admin", "admin"))
	{
		courseTypes.POST("", courseTypeHandler.CreateCourseTypeHandler)
	}
}