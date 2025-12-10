package quizmanagementhandler

import (
	"fmt"
	"net/http"

	"course_online_backend/internal/dto"
	"course_online_backend/internal/services"
	quizservicesgo "course_online_backend/internal/services/Course_management_Services/Quiz_Services.go"
	"course_online_backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type QuestionHandler struct {
	questionService *quizservicesgo.QuestionService
	ActivityService *services.ActivityService
}

func NewQuestionHandler(questionService *quizservicesgo.QuestionService, act *services.ActivityService) *QuestionHandler {
	return &QuestionHandler{questionService: questionService, ActivityService: act}
}

// CreateQuestion godoc
// @Summary Create a new question in a quiz
// @Description Create a single question for a specific quiz using its quiz ID
// @Tags Question
// @Accept json
// @Produce json
// @Param quizId path string true "Quiz ID (UUID format)" format(uuid)
// @Param request body dto.CreateQuestionRequest true "Create Question Request"
// @Success 201 {object} utils.StandardResponse{data=dto.QuestionResponse} "Question created successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid quiz ID or invalid request body"
// @Failure 404 {object} utils.ErrorResponse "Quiz not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/questions/quizzes/{quizId}/questions [post]
func (h *QuestionHandler) CreateQuestion(c *gin.Context) {
	quizID, err := uuid.Parse(c.Param("quizId"))
	if err != nil {
		utils.JSONError(c, "Invalid quiz ID", http.StatusBadRequest, err.Error())
		return
	}

	var req dto.CreateQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, "Invalid request data", http.StatusBadRequest, err.Error())
		return
	}

	question, err := h.questionService.CreateQuestion(quizID, req)
	if err != nil {
		if err.Error() == "quiz not found" {
			utils.JSONNotFound(c, "Quiz not found")
			return
		}
		if err.Error() == "exactly one answer must be marked as correct" {
			utils.JSONError(c, "Exactly one answer must be marked as correct", http.StatusBadRequest, nil)
			return
		}
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONCreated(c, question, "Question created successfully")

	if firebaseUID, exists := c.Get("user_id"); exists {
		go func() {
			user, err := h.ActivityService.GetUserByFirebaseUID(firebaseUID.(string))
			if err != nil {
				fmt.Printf("User not found for activity logging")
				return
			}

			activityMsg := fmt.Sprintf(
				"Created question: %s (QuizID: %s, QuestionID: %s)",
				req.QuestionTitle,
				quizID.String(),
				question.ID.String(),
			)

			_ = h.ActivityService.LogActivity(user.ID, activityMsg)
		}()
	}
}

// BulkCreateQuestions godoc
// @Summary Bulk create questions
// @Description Create multiple questions for a quiz in one request
// @Tags Question
// @Accept json
// @Produce json
// @Param quizId path string true "Quiz ID (UUID format)" format(uuid)
// @Param request body dto.BulkCreateQuestionsRequest true "Bulk Create Questions Request"
// @Success 201 {object} utils.StandardResponse{data=[]dto.QuestionResponse} "Questions created successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid quiz ID or invalid request format"
// @Failure 404 {object} utils.ErrorResponse "Quiz not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/questions/quizzes/{quizId}/questions/bulk [post]
func (h *QuestionHandler) BulkCreateQuestions(c *gin.Context) {
	quizID, err := uuid.Parse(c.Param("quizId"))
	if err != nil {
		utils.JSONError(c, "Invalid quiz ID", http.StatusBadRequest, err.Error())
		return
	}

	var req dto.BulkCreateQuestionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, "Invalid request data", http.StatusBadRequest, err.Error())
		return
	}

	questions, err := h.questionService.BulkCreateQuestions(quizID, req)
	if err != nil {
		if err.Error() == "quiz not found" {
			utils.JSONNotFound(c, "Quiz not found")
			return
		}
		if err.Error() == "exactly one answer must be marked as correct for each question" {
			utils.JSONError(c, "Exactly one answer must be marked as correct for each question", http.StatusBadRequest, nil)
			return
		}
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONCreated(c, questions, "Questions created successfully")

	if firebaseUID, exists := c.Get("user_id"); exists {
		go func() {
			user, err := h.ActivityService.GetUserByFirebaseUID(firebaseUID.(string))
			if err != nil {
				fmt.Printf("User not found for activity logging")
				return
			}

			activityMsg := fmt.Sprintf(
				"Bulk created %d questions for QuizID: %s",
				len(req.Questions),
				quizID.String(),
			)

			_ = h.ActivityService.LogActivity(user.ID, activityMsg)
		}()
	}
}

// BulkCreateQuestions godoc
// @Summary Bulk create questions
// @Description Create multiple questions for a quiz in one request
// @Tags Question
// @Accept json
// @Produce json
// @Param quizId path string true "Quiz ID (UUID format)" format(uuid)
// @Param request body dto.BulkCreateQuestionsRequest true "Bulk Create Questions Request"
// @Success 201 {object} utils.StandardResponse{data=[]dto.QuestionResponse} "Questions created successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid quiz ID or invalid request format"
// @Failure 404 {object} utils.ErrorResponse "Quiz not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/questions/quizzes/{quizId}/questions/bulk [post]
func (h *QuestionHandler) GetQuestionsByQuiz(c *gin.Context) {
	quizID, err := uuid.Parse(c.Param("quizId"))
	if err != nil {
		utils.JSONError(c, "Invalid quiz ID", http.StatusBadRequest, err.Error())
		return
	}

	questions, err := h.questionService.GetQuestionsByQuiz(quizID)
	if err != nil {
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONSuccess(c, questions, "Questions retrieved successfully")
}

// GetQuestionByID godoc
// @Summary Get a question by ID
// @Description Retrieve a single question using its question ID
// @Tags Question
// @Accept json
// @Produce json
// @Param questionId path string true "Question ID (UUID format)" format(uuid)
// @Success 200 {object} utils.StandardResponse{data=dto.QuestionResponse} "Question retrieved successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid question ID"
// @Failure 404 {object} utils.ErrorResponse "Question not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/questions/{questionId} [get]
func (h *QuestionHandler) GetQuestionByID(c *gin.Context) {
	questionID, err := uuid.Parse(c.Param("questionId"))
	if err != nil {
		utils.JSONError(c, "Invalid question ID", http.StatusBadRequest, err.Error())
		return
	}

	question, err := h.questionService.GetQuestionByID(questionID)
	if err != nil {
		if err.Error() == "question not found" {
			utils.JSONNotFound(c, "Question not found")
			return
		}
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONSuccess(c, question, "Question retrieved successfully")
}

// UpdateQuestion godoc
// @Summary Update a question
// @Description Update a question's title or answers using its ID
// @Tags Question
// @Accept json
// @Produce json
// @Param questionId path string true "Question ID (UUID format)" format(uuid)
// @Param request body dto.UpdateQuestionRequest true "Update Question Request"
// @Success 200 {object} utils.StandardResponse{data=dto.QuestionResponse} "Question updated successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid question ID or invalid request format"
// @Failure 404 {object} utils.ErrorResponse "Question not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/questions/{questionId} [put]
func (h *QuestionHandler) UpdateQuestion(c *gin.Context) {
	questionID, err := uuid.Parse(c.Param("questionId"))
	if err != nil {
		utils.JSONError(c, "Invalid question ID", http.StatusBadRequest, err.Error())
		return
	}

	var req dto.UpdateQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, "Invalid request data", http.StatusBadRequest, err.Error())
		return
	}

	question, err := h.questionService.UpdateQuestion(questionID, req)
	if err != nil {
		if err.Error() == "question not found" {
			utils.JSONNotFound(c, "Question not found")
			return
		}
		if err.Error() == "exactly one answer must be marked as correct" {
			utils.JSONError(c, "Exactly one answer must be marked as correct", http.StatusBadRequest, nil)
			return
		}
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONSuccess(c, question, "Question updated successfully")

	if firebaseUID, exists := c.Get("user_id"); exists {
		go func() {
			user, err := h.ActivityService.GetUserByFirebaseUID(firebaseUID.(string))
			if err != nil {
				fmt.Printf("User not found for activity logging")
				return
			}

			activityMsg := fmt.Sprintf(
				"Updated question: %s (QuestionID: %s)",
				req.QuestionTitle,
				questionID.String(),
			)

			_ = h.ActivityService.LogActivity(user.ID, activityMsg)
		}()
	}
}

// SoftDeleteQuestion godoc
// @Summary Soft delete a question
// @Description Soft delete a question by marking is_deleted = true
// @Tags Question
// @Accept json
// @Produce json
// @Param questionId path string true "Question ID (UUID format)" format(uuid)
// @Success 200 {object} utils.StandardResponse "Question soft deleted successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid question ID"
// @Failure 404 {object} utils.ErrorResponse "Question not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/questions/{questionId} [delete]
func (h *QuestionHandler) SoftDeleteQuestion(c *gin.Context) {
	questionID, err := uuid.Parse(c.Param("questionId"))
	if err != nil {
		utils.JSONError(c, "Invalid question ID", http.StatusBadRequest, err.Error())
		return
	}

	if err := h.questionService.SoftDeleteQuestion(questionID); err != nil {
		if err.Error() == "question not found" {
			utils.JSONNotFound(c, "Question not found")
			return
		}
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONSuccess(c, nil, "Question soft deleted successfully")

	if firebaseUID, exists := c.Get("user_id"); exists {
		go func() {
			user, err := h.ActivityService.GetUserByFirebaseUID(firebaseUID.(string))
			if err != nil {
				fmt.Printf("User not found for activity logging")
				return
			}

			activityMsg := fmt.Sprintf("Soft deleted question (QuestionID: %s)", questionID.String())

			_ = h.ActivityService.LogActivity(user.ID, activityMsg)
		}()
	}
}

// PermanentDeleteQuestion godoc
// @Summary Permanently delete a question
// @Description Completely remove a question from the database
// @Tags Question
// @Accept json
// @Produce json
// @Param questionId path string true "Question ID (UUID format)" format(uuid)
// @Success 200 {object} utils.StandardResponse "Question permanently deleted"
// @Failure 400 {object} utils.ErrorResponse "Invalid question ID"
// @Failure 404 {object} utils.ErrorResponse "Question not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/questions/{questionId}/permanent [delete]
func (h *QuestionHandler) PermanentDeleteQuestion(c *gin.Context) {
	questionID, err := uuid.Parse(c.Param("questionId"))
	if err != nil {
		utils.JSONError(c, "Invalid question ID", http.StatusBadRequest, err.Error())
		return
	}

	if err := h.questionService.PermanentDeleteQuestion(questionID); err != nil {
		if err.Error() == "question not found" {
			utils.JSONNotFound(c, "Question not found")
			return
		}
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONSuccess(c, nil, "Question permanently deleted")

	if firebaseUID, exists := c.Get("user_id"); exists {
		go func() {
			user, err := h.ActivityService.GetUserByFirebaseUID(firebaseUID.(string))
			if err != nil {
				fmt.Printf("User not found for activity logging")
				return
			}

			activityMsg := fmt.Sprintf("Permanently deleted question (QuestionID: %s)", questionID.String())

			_ = h.ActivityService.LogActivity(user.ID, activityMsg)
		}()
	}
}

// GetDeletedQuestions godoc
// @Summary Get soft-deleted questions from a quiz
// @Description Retrieve all questions that were soft-deleted (is_deleted = true) for a specific quiz
// @Tags Question
// @Accept json
// @Produce json
// @Param quizId path string true "Quiz ID (UUID format)" format(uuid)
// @Success 200 {object} utils.StandardResponse{data=[]dto.QuestionResponse} "Deleted questions retrieved successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid quiz ID"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/questions/quizzes/{quizId}/questions/deleted [get]
func (h *QuestionHandler) GetDeletedQuestions(c *gin.Context) {
	quizID, err := uuid.Parse(c.Param("quizId"))
	if err != nil {
		utils.JSONError(c, "Invalid quiz ID", http.StatusBadRequest, err.Error())
		return
	}

	questions, err := h.questionService.GetDeletedQuestions(quizID)
	if err != nil {
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONSuccess(c, questions, "Deleted questions retrieved successfully")
}

// RestoreQuestion godoc
// @Summary Restore a soft-deleted question
// @Description Restore a previously soft-deleted question (is_deleted = false)
// @Tags Question
// @Accept json
// @Produce json
// @Param questionId path string true "Question ID (UUID format)" format(uuid)
// @Success 200 {object} utils.StandardResponse "Question restored successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid question ID or question is not deleted"
// @Failure 404 {object} utils.ErrorResponse "Question not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/questions/{questionId}/restore [post]
func (h *QuestionHandler) RestoreQuestion(c *gin.Context) {
	questionID, err := uuid.Parse(c.Param("questionId"))
	if err != nil {
		utils.JSONError(c, "Invalid question ID", http.StatusBadRequest, err.Error())
		return
	}

	if err := h.questionService.RestoreQuestion(questionID); err != nil {
		if err.Error() == "question not found" {
			utils.JSONNotFound(c, "Question not found")
			return
		}
		if err.Error() == "question is not deleted" {
			utils.JSONError(c, "Question is not deleted", http.StatusBadRequest, nil)
			return
		}
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONSuccess(c, nil, "Question restored successfully")

	if firebaseUID, exists := c.Get("user_id"); exists {
		go func() {
			user, err := h.ActivityService.GetUserByFirebaseUID(firebaseUID.(string))
			if err != nil {
				fmt.Printf("User not found for activity logging")
				return
			}

			activityMsg := fmt.Sprintf("Restored question (QuestionID: %s)", questionID.String())

			_ = h.ActivityService.LogActivity(user.ID, activityMsg)
		}()
	}
}
