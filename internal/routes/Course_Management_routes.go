package routes

import (
	CourseManagementhandler "course_online_backend/internal/handler/Course_Management"
	"course_online_backend/internal/middleware"
	CourseManagementServices "course_online_backend/internal/services/Course_management_Services"

	firebase "firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CourseManagementRoutes(r *gin.RouterGroup, db *gorm.DB, app *firebase.App) {
	courseServices := CourseManagementServices.NewCourseService(db)
	courseHandler := CourseManagementhandler.NewCourseHandler(courseServices)

	course := r.Group("/courses")

	course.GET("/", courseHandler.GetAllCourseHandler)
	course.GET("/:id", courseHandler.GetByIDCourseHandler)

	protected := course.Group("/")
	protected.Use(middleware.RoleMiddleware("super_admin", "admin", "instructor"))

	protected.POST("/", courseHandler.CreateCourseHandler)
	protected.PUT("/:id", middleware.CheckCourseOwnership(db), courseHandler.UpdateCourseHandler)
	protected.DELETE("/:id", middleware.CheckCourseOwnership(db), courseHandler.SoftDeleteCourseHandler)
	protected.PATCH("/:id/restore", courseHandler.RestoreCourseHandler)
	protected.DELETE("/:id/permanent", courseHandler.PermanentDeleteCourseHandler)


	//course type routes
	courseTypeServices := CourseManagementServices.NewCourseTypeService(db)
	courseTypeHandler := CourseManagementhandler.NewCourseTypeHandler(courseTypeServices)

	courseType := r.Group("/course-types")
	courseType.Use(middleware.RoleMiddleware("super_admin", "admin"))
	courseType.POST("/", courseTypeHandler.CreateCourseTypeHandler)
	courseType.GET("/", courseTypeHandler.GetAllCourseTypeHandler)
}
