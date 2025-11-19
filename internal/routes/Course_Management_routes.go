package routes

import (
	CourseManagementhandler "course_online_backend/internal/handler/Course_Management"
	"course_online_backend/internal/middleware"
	"course_online_backend/internal/services"
	CourseManagementServices "course_online_backend/internal/services/Course_management_Services"

	ModuleMgmthandler "course_online_backend/internal/handler/Course_Management/Module_Management_Handler"
	modulemanagementServices "course_online_backend/internal/services/Course_management_Services/Module_Management_Services"

	firebase "firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CoursePublicRoutes(r *gin.RouterGroup, db *gorm.DB, app *firebase.App) {
	activityService := services.NewActivityService(db)
	courseServices := CourseManagementServices.NewCourseService(db)
	courseHandler := CourseManagementhandler.NewCourseHandler(courseServices, activityService)

	course := r.Group("/courses")
	{
		course.GET("/category/:categoryId", courseHandler.GetCoursesByCategoryHandler)
		course.GET("/popular", courseHandler.GetPopularCoursesHandler)
		course.GET("/latest", courseHandler.GetLatestCoursesHandler)
		course.GET("/instructor/:instructorId", courseHandler.GetCoursesByInstructorHandler)
		course.GET("/:id/related", courseHandler.GetRelatedCoursesHandler)
		course.GET("/stats", courseHandler.GetCourseStatsHandler)
	}
}

func CourseProtectedRoutes(r *gin.RouterGroup, db *gorm.DB, app *firebase.App) {
	activityService := services.NewActivityService(db)

	courseServices := CourseManagementServices.NewCourseService(db)
	courseHandler := CourseManagementhandler.NewCourseHandler(courseServices, activityService)

	course := r.Group("/courses")
	course.Use(middleware.RoleMiddleware("super_admin", "admin", "instructor"))
	{
		course.POST("/", courseHandler.CreateCourseHandler)
		course.PUT("/:id", middleware.CheckCourseOwnership(db), courseHandler.UpdateCourseHandler)
		course.DELETE("/:id", middleware.CheckCourseOwnership(db), courseHandler.SoftDeleteCourseHandler)
		course.PATCH("/:id/restore", courseHandler.RestoreCourseHandler)
		course.DELETE("/:id/permanent", courseHandler.PermanentDeleteCourseHandler)
	}

	//course Browsing Routes for all roles
	courses := r.Group("/courses")
	courses.Use(middleware.RoleMiddleware("student", "super_admin", "admin", "instructor"))
	{
		course.GET("/", courseHandler.GetAllCourseHandler)
		course.GET("/:id", courseHandler.GetByIDCourseHandler)
		courses.GET("/my-courses", courseHandler.GetMyCoursesHandler)
		course.GET("/course-types", courseHandler.GetAllCourseTypesHandler)
	}

	// Course Type Routes
	courseTypeServices := CourseManagementServices.NewCourseTypeService(db)
	courseTypeHandler := CourseManagementhandler.NewCourseTypeHandler(courseTypeServices, activityService)

	courseType := r.Group("/course-types")
	courseType.Use(middleware.RoleMiddleware("super_admin", "admin"))
	{
		courseType.POST("/", courseTypeHandler.CreateCourseTypeHandler)
		courseType.GET("/", courseTypeHandler.GetAllCourseTypeHandler)
	}

	//module routes
	moduleServices := modulemanagementServices.NewModuleService(db)
	moduleHandler := ModuleMgmthandler.NewModuleHandler(moduleServices, activityService)

	module := r.Group("/courses")
	module.Use(middleware.RoleMiddleware("super_admin", "admin", "instructor"))
	{
		module.GET("/:id/modules", moduleHandler.GetModulesByCourseHandler)

		module.POST("/:id/modules", moduleHandler.CreateModuleHandler)

		module.GET("/modules/:id", moduleHandler.GetModuleByIDHandler)
		module.PUT("/modules/:id", moduleHandler.UpdateModuleHandler)
		module.DELETE("/modules/:id", moduleHandler.SoftDeleteModuleHandler)
		module.PATCH("/modules/:id/restore", moduleHandler.RestoreModuleHandler)
		module.DELETE("/modules/:id/permanent", moduleHandler.PermanentDeleteModuleHandler)
	}	
}
