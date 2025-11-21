package Course_management_routes


import (
	CourseManagementhandler "course_online_backend/internal/handler/Course_Management"
	"course_online_backend/internal/services"
	CourseManagementServices "course_online_backend/internal/services/Course_management_Services"

	firebase "firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CourseBrowsingRoutes(r *gin.RouterGroup, db *gorm.DB, app *firebase.App) {
	activityService := services.NewActivityService(db)
	courseServices := CourseManagementServices.NewCourseService(db)
	courseHandler := CourseManagementhandler.NewCourseHandler(courseServices, activityService)

	courses := r.Group("/courses")
	{
		courses.GET("", courseHandler.GetAllCourseHandler)
		courses.GET("/popular", courseHandler.GetPopularCoursesHandler)
		courses.GET("/latest", courseHandler.GetLatestCoursesHandler)
		courses.GET("/stats", courseHandler.GetCourseStatsHandler)
		courses.GET("/:id", courseHandler.GetByIDCourseHandler)
		courses.GET("/:id/related", courseHandler.GetRelatedCoursesHandler)
	}

	categories := r.Group("/categories/:categoryId/courses")
	{
		categories.GET("", courseHandler.GetCoursesByCategoryHandler)
	}

	instructors := r.Group("/instructors/:instructorId/courses")
	{
		instructors.GET("", courseHandler.GetCoursesByInstructorHandler)
	}
}