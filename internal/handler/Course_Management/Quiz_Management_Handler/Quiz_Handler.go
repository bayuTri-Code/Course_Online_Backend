package quizmanagementhandler

import (
	"fmt"
	"net/http"

	"course_online_backend/internal/dto"
	"course_online_backend/internal/models"
	"course_online_backend/internal/services"
	quizservicesgo "course_online_backend/internal/services/Course_management_Services/Quiz_Services.go"
	"course_online_backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuizHandler struct {
	quizService     *quizservicesgo.QuizService
	ActivityService *services.ActivityService
	db              *gorm.DB
}

func NewQuizHandler(quizService *quizservicesgo.QuizService, act *services.ActivityService, db *gorm.DB) *QuizHandler {
	return &QuizHandler{
		quizService:     quizService,
		ActivityService: act,
		db:              db,
	}
}

// CreateQuiz godoc
// @Summary Create a new quiz
// @Description Create a new quiz inside a specific course. Only authenticated users can create quizzes.
// @Tags Quiz
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateQuizRequest true "Quiz creation data"
// @Success 201 {object} utils.StandardResponse{data=dto.QuizResponse} "Quiz created successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid request data or user ID format"
// @Failure 401 {object} utils.ErrorResponse "User not authenticated"
// @Failure 404 {object} utils.ErrorResponse "User not found or course not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/quizzes/courses/{courseId}/quizzes [post]
func (h *QuizHandler) CreateQuiz(c *gin.Context) {
	var req dto.CreateQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, "Invalid request data", http.StatusBadRequest, err.Error())
		return
	}

	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utils.JSONUnauthorized(c, "User not authenticated")
		return
	}

	firebaseUID, ok := userIDInterface.(string)
	if !ok {
		utils.JSONError(c, "Invalid user ID format from Firebase token", http.StatusBadRequest, nil)
		return
	}

	var user models.User
	if err := h.db.Where("firebase_uid = ?", firebaseUID).First(&user).Error; err != nil {
		utils.JSONError(c, "User not found", http.StatusNotFound, err.Error())
		return
	}

	quiz, err := h.quizService.CreateQuiz(req, user.ID)
	if err != nil {
		if err.Error() == "course not found" {
			utils.JSONNotFound(c, "Course not found")
			return
		}
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONCreated(c, quiz, "Quiz created successfully")

	go func() {
		firebaseUID, _ := c.Get("user_id")
		user, err := h.ActivityService.GetUserByFirebaseUID(firebaseUID.(string))
		if err == nil {
			msg := fmt.Sprintf("Created a quiz: %s (ID: %s)", quiz.Name, quiz.ID)
			_ = h.ActivityService.LogActivity(user.ID, msg)
		}
	}()
}

// GetQuizzesByCourse godoc
// @Summary Get all quizzes in a course
// @Description Retrieve all active (non-deleted) quizzes belonging to a specific course by course ID
// @Tags Quiz
// @Accept json
// @Produce json
// @Param courseId path string true "Course ID (UUID format)" format(uuid)
// @Success 200 {object} utils.StandardResponse{data=[]dto.QuizResponse} "List of quizzes retrieved successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid course ID format"
// @Failure 404 {object} utils.ErrorResponse "Course not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/quizzes/course/{courseId}/quizzes [get]
func (h *QuizHandler) GetQuizzesByCourse(c *gin.Context) {
	courseID, err := uuid.Parse(c.Param("courseId"))
	if err != nil {
		utils.JSONError(c, "Invalid course ID", http.StatusBadRequest, err.Error())
		return
	}

	quizzes, err := h.quizService.GetQuizzesByCourse(courseID)
	if err != nil {
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONSuccess(c, quizzes, "Quizzes retrieved successfully")
}

// GetQuizByID godoc
// @Summary Get quiz details by ID
// @Description Retrieve detailed information of a specific quiz by its quiz ID
// @Tags Quiz
// @Accept json
// @Produce json
// @Param quizId path string true "Quiz ID (UUID format)" format(uuid)
// @Success 200 {object} utils.StandardResponse{data=dto.QuizResponse} "Quiz details retrieved successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid quiz ID format"
// @Failure 404 {object} utils.ErrorResponse "Quiz not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/quizzes/{quizId} [get]
func (h *QuizHandler) GetQuizByID(c *gin.Context) {
	quizID, err := uuid.Parse(c.Param("quizId"))
	if err != nil {
		utils.JSONError(c, "Invalid quiz ID", http.StatusBadRequest, err.Error())
		return
	}

	quiz, err := h.quizService.GetQuizByID(quizID)
	if err != nil {
		if err.Error() == "quiz not found" {
			utils.JSONNotFound(c, "Quiz not found")
			return
		}
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONSuccess(c, quiz, "Quiz retrieved successfully")
}

// UpdateQuiz godoc
// @Summary Update quiz information
// @Description Update quiz details such as name, description, time limit, passing score, and other fields by quiz ID
// @Tags Quiz
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param quizId path string true "Quiz ID (UUID format)" format(uuid)
// @Param request body dto.UpdateQuizRequest true "Updated quiz data"
// @Success 200 {object} utils.StandardResponse{data=dto.QuizResponse} "Quiz updated successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid quiz ID format or invalid request data"
// @Failure 404 {object} utils.ErrorResponse "Quiz not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/quizzes/{quizId} [put]
func (h *QuizHandler) UpdateQuiz(c *gin.Context) {
	quizID, err := uuid.Parse(c.Param("quizId"))
	if err != nil {
		utils.JSONError(c, "Invalid quiz ID", http.StatusBadRequest, err.Error())
		return
	}

	var req dto.UpdateQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, "Invalid request data", http.StatusBadRequest, err.Error())
		return
	}

	quiz, err := h.quizService.UpdateQuiz(quizID, req)
	if err != nil {
		if err.Error() == "quiz not found" {
			utils.JSONNotFound(c, "Quiz not found")
			return
		}
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONSuccess(c, quiz, "Quiz updated successfully")

	go func() {
		firebaseUID, _ := c.Get("user_id")
		user, err := h.ActivityService.GetUserByFirebaseUID(firebaseUID.(string))
		if err == nil {
			msg := fmt.Sprintf("Updated quiz: %s (ID: %s)", quiz.Name, quiz.ID)
			_ = h.ActivityService.LogActivity(user.ID, msg)
		}
	}()
}

// SoftDeleteQuiz godoc
// @Summary Soft delete a quiz
// @Description Mark a quiz as deleted (soft delete) without permanently removing it from the database. The quiz can be restored later.
// @Tags Quiz
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param quizId path string true "Quiz ID (UUID format)" format(uuid)
// @Success 200 {object} utils.StandardResponse{data=interface{}} "Quiz soft deleted successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid quiz ID format"
// @Failure 404 {object} utils.ErrorResponse "Quiz not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/quizzes/{quizId}/soft-delete [delete]
func (h *QuizHandler) SoftDeleteQuiz(c *gin.Context) {
	quizID, err := uuid.Parse(c.Param("quizId"))
	if err != nil {
		utils.JSONError(c, "Invalid quiz ID", http.StatusBadRequest, err.Error())
		return
	}

	if err := h.quizService.SoftDeleteQuiz(quizID); err != nil {
		if err.Error() == "quiz not found" {
			utils.JSONNotFound(c, "Quiz not found")
			return
		}
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONSuccess(c, nil, "Quiz soft deleted successfully")

	go func() {
		firebaseUID, _ := c.Get("user_id")
		user, err := h.ActivityService.GetUserByFirebaseUID(firebaseUID.(string))
		if err == nil {
			msg := fmt.Sprintf("Soft deleted quiz ID: %s", quizID.String())
			_ = h.ActivityService.LogActivity(user.ID, msg)
		}
	}()
}

// PermanentDeleteQuiz godoc
// @Summary Permanently delete a quiz
// @Description Permanently remove a quiz from the database. This action cannot be undone.
// @Tags Quiz
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param quizId path string true "Quiz ID (UUID format)" format(uuid)
// @Success 200 {object} utils.StandardResponse{data=interface{}} "Quiz permanently deleted successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid quiz ID format"
// @Failure 404 {object} utils.ErrorResponse "Quiz not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/quizzes/{quizId}/permanent [delete]
func (h *QuizHandler) PermanentDeleteQuiz(c *gin.Context) {
	quizID, err := uuid.Parse(c.Param("quizId"))
	if err != nil {
		utils.JSONError(c, "Invalid quiz ID", http.StatusBadRequest, err.Error())
		return
	}

	if err := h.quizService.PermanentDeleteQuiz(quizID); err != nil {
		if err.Error() == "quiz not found" {
			utils.JSONNotFound(c, "Quiz not found")
			return
		}
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONSuccess(c, nil, "Quiz permanently deleted")

	go func() {
		firebaseUID, _ := c.Get("user_id")
		user, err := h.ActivityService.GetUserByFirebaseUID(firebaseUID.(string))
		if err == nil {
			msg := fmt.Sprintf("Permanently deleted quiz ID: %s", quizID.String())
			_ = h.ActivityService.LogActivity(user.ID, msg)
		}
	}()
}

// GetDeletedQuizzes godoc
// @Summary Get all soft-deleted quizzes
// @Description Retrieve all quizzes that have been soft-deleted for a specific course. These quizzes can be restored.
// @Tags Quiz
// @Accept json
// @Produce json
// @Param courseId path string true "Course ID (UUID format)" format(uuid)
// @Success 200 {object} utils.StandardResponse{data=[]dto.QuizResponse} "Deleted quizzes retrieved successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid course ID format"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/quizzes/courses/{courseId}/quizzes/deleted [get]
func (h *QuizHandler) GetDeletedQuizzes(c *gin.Context) {
	courseID, err := uuid.Parse(c.Param("courseId"))
	if err != nil {
		utils.JSONError(c, "Invalid course ID", http.StatusBadRequest, err.Error())
		return
	}

	quizzes, err := h.quizService.GetDeletedQuizzes(courseID)
	if err != nil {
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONSuccess(c, quizzes, "Deleted quizzes retrieved successfully")
}

// RestoreQuiz godoc
// @Summary Restore a soft-deleted quiz
// @Description Restore a previously soft-deleted quiz back to active state, making it available again in the course
// @Tags Quiz
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param quizId path string true "Quiz ID (UUID format)" format(uuid)
// @Success 200 {object} utils.StandardResponse{data=interface{}} "Quiz restored successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid quiz ID format or quiz is not in deleted state"
// @Failure 404 {object} utils.ErrorResponse "Quiz not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/quizzes/{quizId}/restore [put]
func (h *QuizHandler) RestoreQuiz(c *gin.Context) {
	quizID, err := uuid.Parse(c.Param("quizId"))
	if err != nil {
		utils.JSONError(c, "Invalid quiz ID", http.StatusBadRequest, err.Error())
		return
	}

	if err := h.quizService.RestoreQuiz(quizID); err != nil {
		if err.Error() == "quiz not found" {
			utils.JSONNotFound(c, "Quiz not found")
			return
		}
		if err.Error() == "quiz is not deleted" {
			utils.JSONError(c, "Quiz is not deleted", http.StatusBadRequest, nil)
			return
		}
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONSuccess(c, nil, "Quiz restored successfully")

	go func() {
		firebaseUID, _ := c.Get("user_id")
		user, err := h.ActivityService.GetUserByFirebaseUID(firebaseUID.(string))
		if err == nil {
			msg := fmt.Sprintf("Restored quiz ID: %s", quizID.String())
			_ = h.ActivityService.LogActivity(user.ID, msg)
		}
	}()
}