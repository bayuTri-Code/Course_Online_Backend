package quiz

import (
	quizmanagementhandler "course_online_backend/internal/handler/Course_Management/Quiz_Management_Handler"
	"course_online_backend/internal/middleware"
	"course_online_backend/internal/services"
	quizservicesgo "course_online_backend/internal/services/Course_management_Services/Quiz_Services.go"

	firebase "firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func StudentQuizAnswerRoutes(r *gin.RouterGroup, db *gorm.DB, app *firebase.App) {
	studentQuizService := quizservicesgo.NewStudentQuizService(db)
	activityService := services.NewActivityService(db)
	studentQuizHandler := quizmanagementhandler.NewStudentQuizHandler(studentQuizService, db, activityService)

	studentRoutes := r.Group("/student")
	studentRoutes.Use(middleware.FirebaseAuth(app, db))
	studentRoutes.Use(middleware.RoleMiddleware("student"))
	{
		studentRoutes.POST("/quizzes/:quizId/submit", middleware.CheckEnrollmentAccess(db), studentQuizHandler.SubmitQuiz)

		studentRoutes.GET("/courses/:courseId/quizzes", middleware.CheckEnrollmentAccess(db), studentQuizHandler.GetQuizzesInCourse)
		studentRoutes.GET("/quizzes/:quizId/start", middleware.CheckEnrollmentAccess(db), studentQuizHandler.StartQuiz)
		studentRoutes.GET("/quizzes/:quizId/attempts", middleware.CheckEnrollmentAccess(db), studentQuizHandler.GetQuizAttemptsHistory)

		studentRoutes.GET("/quizzes/attempts/:attemptId", studentQuizHandler.GetAttemptDetail)
		studentRoutes.GET("/my-attempts", studentQuizHandler.GetAllMyAttempts)
	}
}