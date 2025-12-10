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

func QuizQuestionAndAnswerRoutes(r *gin.RouterGroup, db *gorm.DB, app *firebase.App) {
	questionService := quizservicesgo.NewQuestionService(db)
	activityService := services.NewActivityService(db)
	questionHandler := quizmanagementhandler.NewQuestionHandler(questionService, activityService)

	questionRoutes := r.Group("/questions")
	questionRoutes.Use(middleware.FirebaseAuth(app, db))
	questionRoutes.Use(middleware.RoleMiddleware("super_admin", "admin"))
	{
		questionRoutes.POST("/quizzes/:quizId/questions", middleware.CheckCourseOwnershipDynamic(db), questionHandler.CreateQuestion)
		questionRoutes.POST("/quizzes/:quizId/questions/bulk", middleware.CheckCourseOwnershipDynamic(db), questionHandler.BulkCreateQuestions)
		questionRoutes.PUT("/:questionId", middleware.CheckCourseOwnershipDynamic(db), questionHandler.UpdateQuestion)

		questionRoutes.GET("/quizzes/:quizId/questions", questionHandler.GetQuestionsByQuiz)
		questionRoutes.GET("/:questionId", questionHandler.GetQuestionByID)
		questionRoutes.GET("/quizzes/:quizId/questions/deleted", questionHandler.GetDeletedQuestions)

		questionRoutes.DELETE("/:questionId", middleware.CheckCourseOwnershipDynamic(db), questionHandler.SoftDeleteQuestion)
		questionRoutes.DELETE("/:questionId/permanent", middleware.CheckCourseOwnershipDynamic(db), questionHandler.PermanentDeleteQuestion)
		questionRoutes.POST("/:questionId/restore", middleware.CheckCourseOwnershipDynamic(db), questionHandler.RestoreQuestion)
	}
}