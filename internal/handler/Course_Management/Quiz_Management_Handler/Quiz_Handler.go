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
// @Summary Create quiz (multipart)
// @Description Create quiz with optional thumbnail upload
// @Tags Quiz
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param course_id path string true "Course ID"
// @Param name formData string true "Quiz name"
// @Param number formData int true "Quiz order"
// @Param min_pass_score formData int false "Minimum pass score"
// @Param is_pass_required formData bool false "Is pass required"
// @Param thumbnail formData file false "Quiz thumbnail"
// @Success 201 {object} utils.StandardResponse{data=dto.QuizResponse}
// @Router /api/quizzes/courses/{course_id}/quizzes [post]
func (h *QuizHandler) CreateQuiz(c *gin.Context) {
	courseIDParam := c.Param("course_id")
	courseUUID, err := uuid.Parse(courseIDParam)
	if err != nil {
		utils.JSONError(c, "Invalid course_id", http.StatusBadRequest, err.Error())
		return
	}

	var form dto.CreateQuizForm
	if err := c.ShouldBind(&form); err != nil {
		utils.JSONError(c, "Invalid form data", http.StatusBadRequest, err.Error())
		return
	}

	userFirebaseID, _ := c.Get("user_id")
	user, err := h.ActivityService.GetUserByFirebaseUID(userFirebaseID.(string))
	if err != nil {
		utils.JSONUnauthorized(c, "User not authenticated")
		return
	}

	if file, fileHeader, err := c.Request.FormFile("thumbnail"); err == nil {
		defer file.Close()

		minio := services.MinioHelper{}
		url, err := minio.UploadQuizThumbnail(file, fileHeader)
		if err != nil {
			utils.JSONError(c, "Failed upload thumbnail", http.StatusBadRequest, err.Error())
			return
		}
		form.ThumbnailURL = url
	}

	req := dto.CreateQuizRequest{
		CourseID:       courseUUID,
		Name:           form.Name,
		Number:         form.Number,
		MinPassScore:   form.MinPassScore,
		IsPassRequired: form.IsPassRequired,
		ThumbnailURL:   form.ThumbnailURL,
	}

	quiz, err := h.quizService.CreateQuiz(req, user.ID)
	if err != nil {
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONCreated(c, quiz, "Quiz created")

	go h.ActivityService.LogActivity(
		user.ID,
		fmt.Sprintf("Created quiz: %s", quiz.Name),
	)
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
// @Summary Update quiz (multipart)
// @Description Update quiz data and optional thumbnail
// @Tags Quiz
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param quizId path string true "Quiz ID"
// @Param name formData string false "Quiz name"
// @Param number formData int false "Quiz order"
// @Param min_pass_score formData int false "Minimum pass score"
// @Param is_pass_required formData bool false "Is pass required"
// @Param thumbnail formData file false "Quiz thumbnail"
// @Success 200 {object} utils.StandardResponse{data=dto.QuizResponse}
// @Router /api/quizzes/{quizId} [put]
func (h *QuizHandler) UpdateQuiz(c *gin.Context) {
	quizID, err := uuid.Parse(c.Param("quizId"))
	if err != nil {
		utils.JSONError(c, "Invalid quiz ID", http.StatusBadRequest, err.Error())
		return
	}

	var form dto.UpdateQuizForm
	if err := c.ShouldBind(&form); err != nil {
		utils.JSONError(c, "Invalid form data", http.StatusBadRequest, err.Error())
		return
	}

	file, fileHeader, err := c.Request.FormFile("thumbnail")
	if err == nil {
		defer file.Close()

		minio := services.MinioHelper{}
		url, err := minio.UploadQuizThumbnail(file, fileHeader)
		if err != nil {
			utils.JSONError(c, "Failed upload thumbnail", http.StatusBadRequest, err.Error())
			return
		}
		form.ThumbnailURL = url
	}

	quiz, err := h.quizService.UpdateQuiz(quizID, form)
	if err != nil {
		if err.Error() == "quiz not found" {
			utils.JSONNotFound(c, "Quiz not found")
			return
		}
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONSuccess(c, quiz, "Quiz updated successfully")
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
// @Router /api/quizzes/{quizId} [delete]
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