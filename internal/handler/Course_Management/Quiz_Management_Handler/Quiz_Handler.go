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
	db *gorm.DB
}

func NewQuizHandler(quizService *quizservicesgo.QuizService,act *services.ActivityService  , db *gorm.DB) *QuizHandler {
	return &QuizHandler{
		quizService: quizService,
		ActivityService: act,
		db:          db,
	}
}

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
