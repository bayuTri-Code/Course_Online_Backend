package quizmanagementhandler

import (
	"net/http"

	"course_online_backend/internal/dto"
	"course_online_backend/internal/models"
	quizservicesgo "course_online_backend/internal/services/Course_management_Services/Quiz_Services.go"
	"course_online_backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StudentQuizHandler struct {
	service *quizservicesgo.StudentQuizService
	db      *gorm.DB
}

func NewStudentQuizHandler(service *quizservicesgo.StudentQuizService, db *gorm.DB) *StudentQuizHandler {
	return &StudentQuizHandler{
		service: service,
		db:      db,
	}
}

func (h *StudentQuizHandler) getUserIDFromContext(c *gin.Context) (uuid.UUID, error) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, gorm.ErrRecordNotFound
	}

	firebaseUID, ok := userIDInterface.(string)
	if !ok {
		return uuid.Nil, gorm.ErrInvalidData
	}

	var user models.User
	if err := h.db.Select("id").Where("firebase_uid = ?", firebaseUID).First(&user).Error; err != nil {
		return uuid.Nil, err
	}

	return user.ID, nil
}

func (h *StudentQuizHandler) GetQuizzesInCourse(c *gin.Context) {
	courseID, err := uuid.Parse(c.Param("courseId"))
	if err != nil {
		utils.JSONError(c, "Invalid course ID", http.StatusBadRequest, err.Error())
		return
	}

	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.JSONUnauthorized(c, "User not authenticated")
		} else {
			utils.JSONError(c, "User not found", http.StatusNotFound, err.Error())
		}
		return
	}

	quizzes, err := h.service.GetQuizzesInCourse(courseID, userID)
	if err != nil {
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONSuccess(c, quizzes, "Quizzes retrieved successfully")
}

func (h *StudentQuizHandler) StartQuiz(c *gin.Context) {
	quizID, err := uuid.Parse(c.Param("quizId"))
	if err != nil {
		utils.JSONError(c, "Invalid quiz ID", http.StatusBadRequest, err.Error())
		return
	}

	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.JSONUnauthorized(c, "User not authenticated")
		} else {
			utils.JSONError(c, "User not found", http.StatusNotFound, err.Error())
		}
		return
	}

	quiz, err := h.service.StartQuiz(quizID, userID)
	if err != nil {
		if err.Error() == "quiz not found" {
			utils.JSONNotFound(c, "Quiz not found")
			return
		}
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONSuccess(c, quiz, "Quiz started successfully")
}

func (h *StudentQuizHandler) SubmitQuiz(c *gin.Context) {
	quizID, err := uuid.Parse(c.Param("quizId"))
	if err != nil {
		utils.JSONError(c, "Invalid quiz ID", http.StatusBadRequest, err.Error())
		return
	}

	var req dto.SubmitQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, "Invalid request data", http.StatusBadRequest, err.Error())
		return
	}

	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.JSONUnauthorized(c, "User not authenticated")
		} else {
			utils.JSONError(c, "User not found", http.StatusNotFound, err.Error())
		}
		return
	}

	result, err := h.service.SubmitQuiz(quizID, userID, req)
	if err != nil {
		if err.Error() == "quiz not found" {
			utils.JSONNotFound(c, "Quiz not found")
			return
		}
		if err.Error() == "you must answer all questions" {
			utils.JSONError(c, "You must answer all questions", http.StatusBadRequest, nil)
			return
		}
		if err.Error() == "duplicate answer for question" {
			utils.JSONError(c, "Duplicate answer for question", http.StatusBadRequest, nil)
			return
		}
		if err.Error() == "invalid question_id" {
			utils.JSONError(c, "Invalid question_id", http.StatusBadRequest, nil)
			return
		}
		if err.Error() == "invalid answer_id" {
			utils.JSONError(c, "Invalid answer_id", http.StatusBadRequest, nil)
			return
		}
		if err.Error() == "answer does not belong to the question" {
			utils.JSONError(c, "Answer does not belong to the question", http.StatusBadRequest, nil)
			return
		}
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONCreated(c, result, "Quiz submitted successfully")
}

func (h *StudentQuizHandler) GetQuizAttemptsHistory(c *gin.Context) {
	quizID, err := uuid.Parse(c.Param("quizId"))
	if err != nil {
		utils.JSONError(c, "Invalid quiz ID", http.StatusBadRequest, err.Error())
		return
	}

	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.JSONUnauthorized(c, "User not authenticated")
		} else {
			utils.JSONError(c, "User not found", http.StatusNotFound, err.Error())
		}
		return
	}

	history, err := h.service.GetQuizAttemptsHistory(quizID, userID)
	if err != nil {
		if err.Error() == "quiz not found" {
			utils.JSONNotFound(c, "Quiz not found")
			return
		}
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONSuccess(c, history, "Attempts retrieved successfully")
}

func (h *StudentQuizHandler) GetAttemptDetail(c *gin.Context) {
	attemptID, err := uuid.Parse(c.Param("attemptId"))
	if err != nil {
		utils.JSONError(c, "Invalid attempt ID", http.StatusBadRequest, err.Error())
		return
	}

	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.JSONUnauthorized(c, "User not authenticated")
		} else {
			utils.JSONError(c, "User not found", http.StatusNotFound, err.Error())
		}
		return
	}

	detail, err := h.service.GetAttemptDetail(attemptID, userID)
	if err != nil {
		if err.Error() == "attempt not found" {
			utils.JSONNotFound(c, "Attempt not found")
			return
		}
		if err.Error() == "access denied: this attempt belongs to another user" {
			utils.JSONError(c, "Access denied: This attempt belongs to another user", http.StatusForbidden, nil)
			return
		}
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONSuccess(c, detail, "Attempt detail retrieved successfully")
}

func (h *StudentQuizHandler) GetAllMyAttempts(c *gin.Context) {
	var params dto.AllMyAttemptsQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		utils.JSONError(c, "Invalid query parameters", http.StatusBadRequest, err.Error())
		return
	}

	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.JSONUnauthorized(c, "User not authenticated")
		} else {
			utils.JSONError(c, "User not found", http.StatusNotFound, err.Error())
		}
		return
	}

	attempts, err := h.service.GetAllMyAttempts(userID, params)
	if err != nil {
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONSuccess(c, attempts, "My attempts retrieved successfully")
}