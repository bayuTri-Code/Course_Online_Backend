package enrollmentmanagementroutes

import (
	Enrollment_Handler "course_online_backend/internal/handler/Enrollment_Handler"
	"course_online_backend/internal/middleware"
	enrollmentService "course_online_backend/internal/services/Enrollment_Services"

	"firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func EnrollmentRoutes(r *gin.RouterGroup, db *gorm.DB, app *firebase.App) {
	service := enrollmentService.NewEnrollmentService(db)
	handler := Enrollment_Handler.NewEnrollmentHandler(service)

	student := r.Group("/enrollments")
	student.Use(middleware.RoleMiddleware("student"))
	{
		student.POST("", handler.EnrollCourseHandler)                     
		student.GET("", handler.GetMyEnrollmentsHandler)                  
		student.DELETE("/:enrollment_id", handler.UnenrollCourseHandler)  
	}

	courseEnrollmentsStudent := r.Group("/courses/:course_id/enrollments")
	courseEnrollmentsStudent.Use(middleware.RoleMiddleware("student"))
	{
		courseEnrollmentsStudent.GET("/check", handler.CheckEnrollmentHandler) 
	}

	
	admin := r.Group("/enrollments")
	admin.Use(middleware.RoleMiddleware("admin", "super_admin"))
	{
		admin.PUT("/:enrollment_id/status", handler.UpdateEnrollmentStatusHandler)
	}

	courseEnrollmentsAdmin := r.Group("/courses/:course_id/enrollments")
	courseEnrollmentsAdmin.Use(middleware.RoleMiddleware("admin", "super_admin"))
	{
		courseEnrollmentsAdmin.GET("", handler.GetCourseStudentsHandler)
	}


	shared := r.Group("/enrollments")
	shared.Use(middleware.RoleMiddleware("student", "instructor", "admin", "super_admin"))
	{
		shared.GET("/:enrollment_id", handler.GetEnrollmentDetailHandler) 
	}
}