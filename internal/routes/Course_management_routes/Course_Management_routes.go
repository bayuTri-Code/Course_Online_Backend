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

func CourseManagementRoutes(r *gin.RouterGroup, db *gorm.DB, app *firebase.App) {
	activityService := services.NewActivityService(db)
	courseServices := CourseManagementServices.NewCourseService(db)
	courseHandler := CourseManagementhandler.NewCourseHandler(courseServices, activityService)

	instructorService := CourseManagementServices.NewInstructorService(db)
	instructorHandler := CourseManagementhandler.NewInstructorHandler(instructorService)

	courses := r.Group("/courses")
	courses.Use(middleware.RoleMiddleware("super_admin", "admin"))
	{
		courses.POST("", courseHandler.CreateCourseHandler)
		courses.PUT("/:course_id", middleware.CheckCourseOwnershipDynamic(db), courseHandler.UpdateCourseHandler)
		courses.DELETE("/:course_id", middleware.CheckCourseOwnershipDynamic(db), courseHandler.SoftDeleteCourseHandler)
		courses.PATCH("/:course_id/restore", courseHandler.RestoreCourseHandler)
		courses.DELETE("/:course_id/permanent", courseHandler.PermanentDeleteCourseHandler)

		courses.GET("/instructors/search", instructorHandler.SearchInstructors)
		courses.GET("/instructors/:id", instructorHandler.GetInstructorByID)
	}

	categories := r.Group("/mycourses")
	categories.Use(middleware.RoleMiddleware("student"))
	{
		categories.GET("/categories/:categoryId/courses", courseHandler.GetMyCoursesByCategoryHandler)
	}

	myCourses := r.Group("/my")
	myCourses.Use(middleware.RoleMiddleware("student", "super_admin", "admin", "instructor"))
	{
		myCourses.GET("/courses", courseHandler.GetMyCoursesHandler)
	}
}
