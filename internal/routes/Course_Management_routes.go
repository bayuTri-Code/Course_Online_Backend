package routes

import (
	CourseManagementhandler "course_online_backend/internal/handler/Course_Management"
	"course_online_backend/internal/middleware"
	CourseManagementServices "course_online_backend/internal/services/Course_management_Services"

	firebase "firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CoursePublicRoutes(r *gin.RouterGroup, db *gorm.DB, app *firebase.App) {
	courseServices := CourseManagementServices.NewCourseService(db)
	courseHandler := CourseManagementhandler.NewCourseHandler(courseServices)

	course := r.Group("/courses")
	{
		course.GET("/", courseHandler.GetAllCourseHandler)
		course.GET("/:id", courseHandler.GetByIDCourseHandler)
	}
}

func CourseProtectedRoutes(r *gin.RouterGroup, db *gorm.DB, app *firebase.App) {
	courseServices := CourseManagementServices.NewCourseService(db)
	courseHandler := CourseManagementhandler.NewCourseHandler(courseServices)

	course := r.Group("/courses")
	course.Use(middleware.RoleMiddleware("super_admin", "admin", "instructor"))
	{
		course.POST("/", courseHandler.CreateCourseHandler)
		course.PUT("/:id", middleware.CheckCourseOwnership(db), courseHandler.UpdateCourseHandler)
		course.DELETE("/:id", middleware.CheckCourseOwnership(db), courseHandler.SoftDeleteCourseHandler)
		course.PATCH("/:id/restore", courseHandler.RestoreCourseHandler)
		course.DELETE("/:id/permanent", courseHandler.PermanentDeleteCourseHandler)
	}

	courseTypeServices := CourseManagementServices.NewCourseTypeService(db)
	courseTypeHandler := CourseManagementhandler.NewCourseTypeHandler(courseTypeServices)

	courseType := r.Group("/course-types")
	courseType.Use(middleware.RoleMiddleware("super_admin", "admin"))
	{
		courseType.POST("/", courseTypeHandler.CreateCourseTypeHandler)
		courseType.GET("/", courseTypeHandler.GetAllCourseTypeHandler)
	}
}
