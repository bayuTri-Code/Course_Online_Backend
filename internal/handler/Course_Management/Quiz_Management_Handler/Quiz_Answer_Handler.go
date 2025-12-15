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

type StudentQuizHandler struct {
	service         *quizservicesgo.StudentQuizService
	db              *gorm.DB
	ActivityService *services.ActivityService
}

func NewStudentQuizHandler(
	service *quizservicesgo.StudentQuizService, db *gorm.DB, act *services.ActivityService) *StudentQuizHandler {
		return &StudentQuizHandler{
			service:         service,
			db:              db,
			ActivityService: act,
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

// GetQuizzesInCourse godoc
// @Summary Get quizzes for a course
// @Description Retrieve all quizzes inside a specific course for the authenticated student.
// @Tags students' quiz answers
// @Accept json
// @Produce json
// @Param courseId path string true "Course ID (UUID format)" format(uuid)
// @Success 200 {object} utils.StandardResponse{data=[]dto.QuizInCourseResponse} "Quizzes retrieved successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid course ID"
// @Failure 401 {object} utils.ErrorResponse "User not authenticated"
// @Failure 404 {object} utils.ErrorResponse "User not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/student/courses/{courseId}/quizzes [get]
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

// StartQuiz godoc
// @Summary Start a quiz
// @Description Initialize a quiz attempt and retrieve all questions with answer options.
// @Tags students' quiz answers
// @Accept json
// @Produce json
// @Param quizId path string true "Quiz ID (UUID format)" format(uuid)
// @Success 200 {object} utils.StandardResponse{data=dto.StartQuizResponse} "Quiz started successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid quiz ID"
// @Failure 401 {object} utils.ErrorResponse "User not authenticated"
// @Failure 404 {object} utils.ErrorResponse "Quiz not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/student/quizzes/{quizId}/start [get]
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

	if firebaseUID, exists := c.Get("user_id"); exists {
		go func() {
			user, err := h.ActivityService.GetUserByFirebaseUID(firebaseUID.(string))
			if err != nil {
				fmt.Printf("User not found for activity logging")
				return
			}

			activityMsg := fmt.Sprintf(
				"Started quiz: %s (QuizID: %s)",
				quiz.QuizName,
				quizID.String(),
			)

			_ = h.ActivityService.LogActivity(user.ID, activityMsg)
		}()
	}
}

// SubmitQuiz godoc
// @Summary Submit quiz answers
// @Description Submit all answers to a quiz and receive scoring and result summary.
// @Tags students' quiz answers
// @Accept json
// @Produce json
// @Param quizId path string true "Quiz ID (UUID format)" format(uuid)
// @Param request body dto.SubmitQuizRequest true "Submit quiz request payload"
// @Success 201 {object} utils.StandardResponse{data=dto.SubmitQuizResponse} "Quiz submitted successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid request data or validation error"
// @Failure 401 {object} utils.ErrorResponse "User not authenticated"
// @Failure 404 {object} utils.ErrorResponse "Quiz not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/student/quizzes/{quizId}/submit [post]
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

	if firebaseUID, exists := c.Get("user_id"); exists {
		go func() {
			user, err := h.ActivityService.GetUserByFirebaseUID(firebaseUID.(string))
			if err != nil {
				fmt.Printf("User not found for activity logging")
				return
			}

			activityMsg := fmt.Sprintf(
				"Submitted quiz: %s (QuizID: %s, Score: %d, Passed: %t)",
				result.QuizName,
				quizID.String(),
				result.ScoreAchieved,
				result.IsPassed,
			)

			_ = h.ActivityService.LogActivity(user.ID, activityMsg)
		}()
	}
}

// GetQuizAttemptsHistory godoc
// @Summary Get quiz attempt history
// @Description Retrieve all attempts the student has made for a specific quiz.
// @Tags students' quiz answers
// @Accept json
// @Produce json
// @Param quizId path string true "Quiz ID (UUID format)" format(uuid)
// @Success 200 {object} utils.StandardResponse{data=dto.QuizAttemptsHistoryResponse} "Attempts retrieved successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid quiz ID"
// @Failure 401 {object} utils.ErrorResponse "User not authenticated"
// @Failure 404 {object} utils.ErrorResponse "Quiz not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/student/quizzes/{quizId}/attempts [get]
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

// GetAttemptDetail godoc
// @Summary Get quiz attempt detail
// @Description Retrieve detailed information of a specific quiz attempt, including answers and scoring.
// @Tags students' quiz answers
// @Accept json
// @Produce json
// @Param attemptId path string true "Attempt ID (UUID format)" format(uuid)
// @Success 200 {object} utils.StandardResponse{data=dto.AttemptDetailResponse} "Attempt detail retrieved successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid attempt ID"
// @Failure 401 {object} utils.ErrorResponse "User not authenticated"
// @Failure 403 {object} utils.ErrorResponse "Access denied"
// @Failure 404 {object} utils.ErrorResponse "Attempt not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/student/quizzes/attempts/{attemptId} [get]
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

// GetAllMyAttempts godoc
// @Summary Get all quiz attempts by the authenticated student
// @Description Retrieve all quiz attempts made by the current authenticated student with pagination.
// @Tags students' quiz answers
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param order query string false "Sort order (asc|desc)"
// @Success 200 {object} utils.StandardResponse{data=dto.AllMyAttemptsResponse} "My attempts retrieved successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid query parameters"
// @Failure 401 {object} utils.ErrorResponse "User not authenticated"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/student/quizzes/my-attempts [get]
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
