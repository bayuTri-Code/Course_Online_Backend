package enrollmentmanagementroutes

import (
	"course_online_backend/database/Repository"
	Enrollment_Handler "course_online_backend/internal/handler/Enrollment_Handler"
	"course_online_backend/internal/middleware"
	enrollmentService "course_online_backend/internal/services/Enrollment_Services"

	"firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func EnrollmentRoutes(r *gin.RouterGroup, db *gorm.DB, app *firebase.App) {
	enrollmentRepo := repository.NewEnrollmentRepository(db)
	service := enrollmentService.NewEnrollmentService(enrollmentRepo, db)
	handler := Enrollment_Handler.NewEnrollmentHandler(service)

	student := r.Group("/enrollments")
	student.Use(middleware.RoleMiddleware("student"))
	{
		student.POST("", handler.EnrollCourse)
		student.GET("", handler.GetMyEnrollments)
		student.DELETE("/:id", handler.UnenrollCourse)
	}

	courseEnrollmentsStudent := r.Group("/courses/:course_id/enrollments")
	courseEnrollmentsStudent.Use(middleware.RoleMiddleware("student"))
	{
		courseEnrollmentsStudent.GET("/check", handler.CheckEnrollment)
	}

	courseEnrollmentsAdmin := r.Group("/courses/:course_id/enrollments")
	courseEnrollmentsAdmin.Use(middleware.RoleMiddleware("admin", "super_admin"))
	{
		courseEnrollmentsAdmin.GET("", handler.GetCourseStudents)
	}

	admin := r.Group("/enrollments")
	admin.Use(middleware.RoleMiddleware("admin", "super_admin"))
	{
		admin.PUT("/:id/status", handler.UpdateEnrollmentStatus)
	}

	shared := r.Group("/enrollments")
	shared.Use(middleware.RoleMiddleware("student", "instructor", "admin", "super_admin"))
	{
		shared.GET("/:id", handler.GetEnrollmentDetail)
	}
}