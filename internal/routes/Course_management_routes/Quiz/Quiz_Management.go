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

func QuizManagementRoutes(r *gin.RouterGroup, db *gorm.DB, app *firebase.App) {
	activityService := services.NewActivityService(db)
	quizService := quizservicesgo.NewQuizService(db)
	quizHandler := quizmanagementhandler.NewQuizHandler(quizService, activityService, db)

	quizRoutes := r.Group("/quizzes")
	quizRoutes.Use(middleware.FirebaseAuth(app, db))
	quizRoutes.Use(middleware.RoleMiddleware("super_admin", "admin"))
	{
		quizRoutes.POST("/courses/:courseId/quizzes", middleware.CheckCourseOwnershipDynamic(db), quizHandler.CreateQuiz)
		quizRoutes.GET("/courses/:courseId/quizzes", quizHandler.GetQuizzesByCourse)
		quizRoutes.GET("/:quizId", quizHandler.GetQuizByID)
		quizRoutes.PUT("/:quizId", middleware.CheckCourseOwnershipDynamic(db), quizHandler.UpdateQuiz)
		quizRoutes.DELETE("/:quizId", middleware.CheckCourseOwnershipDynamic(db), quizHandler.SoftDeleteQuiz)
		quizRoutes.DELETE("/:quizId/permanent", middleware.CheckCourseOwnershipDynamic(db), quizHandler.PermanentDeleteQuiz)
		quizRoutes.GET("/courses/:courseId/quizzes/deleted", quizHandler.GetDeletedQuizzes)
		quizRoutes.POST("/:quizId/restore", middleware.CheckCourseOwnershipDynamic(db), quizHandler.RestoreQuiz)
	}
}
