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
