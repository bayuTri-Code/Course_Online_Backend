package routes

import (
	"course_online_backend/database/Repository"
	Enrollment_Handler "course_online_backend/internal/handler/Enrollment_Handler"
	"course_online_backend/internal/middleware"
	enrollmentService "course_online_backend/internal/services/Enrollment_Services"

	"firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func EnrollmentProtectedRoutes(r *gin.RouterGroup, db *gorm.DB, app *firebase.App) {

	enrollmentRepo := repository.NewEnrollmentRepository(db)

	service := enrollmentService.NewEnrollmentService(enrollmentRepo, db)

	handler := Enrollment_Handler.NewEnrollmentHandler(service)


	student := r.Group("/enrollments")
	student.Use(middleware.RoleMiddleware("student"))
	{
		student.POST("/free", handler.EnrollFreeCourse)
		student.GET("/check/:courseID", handler.CheckEnrollment)
		student.GET("/my-courses", handler.GetMyEnrollments)
		student.DELETE("/:id", handler.UnenrollCourse)
	}


	details := r.Group("/enrollments")
	details.Use(middleware.RoleMiddleware("student", "instructor", "admin", "super_admin"))
	{
		details.GET("/:id", handler.GetEnrollmentDetail)
		details.PUT("/:id/status", handler.UpdateEnrollmentStatus)
	}

	
	instructor := r.Group("/enrollments")
	instructor.Use(middleware.RoleMiddleware("instructor","admin", "super_admin"))
	{
		instructor.GET("/course/:courseID/students", handler.GetCourseStudents)
	}

	
	superAdmin := r.Group("/enrollments")
	superAdmin.Use(middleware.RoleMiddleware("super_admin"))
	{
	}
}
