package Course_management_routes

import (
	LessonManagementHandler "course_online_backend/internal/handler/Course_Management/Lesson_Management_Handler"
	"course_online_backend/internal/middleware"
	"course_online_backend/internal/services"
	LessonManagementServices "course_online_backend/internal/services/Course_management_Services/Lesson_Management_Services"

	firebase "firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func LessonPublicRoutes(r *gin.RouterGroup, db *gorm.DB, app *firebase.App) {
	activityService := services.NewActivityService(db)
	lessonService := LessonManagementServices.NewLessonService(db)
	lessonHandler := LessonManagementHandler.NewLessonHandler(lessonService, activityService)

	lessons := r.Group("/lessons")
	lessons.Use(middleware.FirebaseAuth(app, db))
	lessons.Use(middleware.RoleMiddleware("student", "admin", "super_admin", "instructor"))
	lessons.Use(middleware.CheckEnrollmentAccess(db))
	{
		lessons.GET("/:lesson_id", lessonHandler.GetLessonByID)
	}

	modules := r.Group("/modules/:module_id/lessons")
	modules.Use(middleware.FirebaseAuth(app, db))
	modules.Use(middleware.RoleMiddleware("student", "admin", "super_admin", "instructor"))
	modules.Use(middleware.CheckEnrollmentAccess(db))
	{
		modules.GET("", lessonHandler.GetLessonsByModule)
	}
}

func LessonRoutes(r *gin.RouterGroup, db *gorm.DB, app *firebase.App) {
	activityService := services.NewActivityService(db)
	lessonService := LessonManagementServices.NewLessonService(db)
	lessonHandler := LessonManagementHandler.NewLessonHandler(lessonService, activityService)

	lessonGroup := r.Group("/lessons")
	lessonGroup.Use(middleware.FirebaseAuth(app, db))
	lessonGroup.Use(middleware.RoleMiddleware("admin", "super_admin"))
	lessonGroup.Use(middleware.CheckCourseOwnershipDynamic(db))
	{
		lessonGroup.PUT("/:lesson_id", lessonHandler.UpdateLesson)
		lessonGroup.DELETE("/:lesson_id", lessonHandler.DeleteLesson)
		lessonGroup.PATCH("/:lesson_id/restore", lessonHandler.RestoreLesson)
		lessonGroup.DELETE("/:lesson_id/permanent", lessonHandler.PermanentDeleteLesson)
	}

	moduleGroup := r.Group("/modules/:module_id/lessons")
	moduleGroup.Use(middleware.FirebaseAuth(app, db))
	moduleGroup.Use(middleware.RoleMiddleware("admin", "super_admin"))
	moduleGroup.Use(middleware.CheckCourseOwnershipDynamic(db))
	{
		moduleGroup.POST("", lessonHandler.CreateLesson)
		moduleGroup.POST("/bulk", lessonHandler.BulkCreateLessons)
		moduleGroup.PATCH("/reorder", lessonHandler.ReorderLessons)
	}
}


