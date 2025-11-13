package CourseManagementhandler

import (
	"course_online_backend/internal/dto"
	CourseManagementServices "course_online_backend/internal/services/Course_management_Services"
	"course_online_backend/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CourseHandler struct {
	service *CourseManagementServices.CourseService
}

func NewCourseHandler(service *CourseManagementServices.CourseService) *CourseHandler {
	return &CourseHandler{
		service: service,
	}
}

// CreateCourseHandler godoc
// @Summary Create a new course
// @Description Create a new course with optional thumbnail upload
// @Tags courses
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Course name"
// @Param description formData string true "Course description"
// @Param price formData number true "Course price"
// @Param is_progress_limited formData boolean false "Limit user progress (optional)"
// @Param course_type_id formData string true "Course type ID (UUID)"
// @Param thumbnail formData file false "Course thumbnail (optional)"
// @Success 201 {object} utils.StandardResponse{data=dto.CourseDetailResponse}
// @Failure 400 {object} utils.ErrorResponse "Invalid input or bad request"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/courses [post]
// @Security BearerAuth
func (h *CourseHandler) CreateCourseHandler(c *gin.Context) {
	var req dto.CreateCourseRequest

	if err := c.ShouldBind(&req); err != nil {
		utils.JSONError(c, "Invalid request data", http.StatusBadRequest, err.Error())
		return
	}

	courseTypeID, err := uuid.Parse(req.CourseTypeID)
	if err != nil {
		utils.JSONError(c, "Invalid course_type_id format", http.StatusBadRequest, err.Error())
		return
	}

	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utils.JSONUnauthorized(c, "User not authenticated")
		return
	}

	userID, ok := userIDInterface.(string)
	if !ok {
		utils.JSONError(c, "Invalid user ID format from Firebase token", http.StatusBadRequest, nil)
		return
	}

	thumbnail, _ := c.FormFile("thumbnail")

	course, err := h.service.CreateCourse(c.Request.Context(), &req, courseTypeID, userID, thumbnail)
	if err != nil {
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONCreated(c, course, "Course created successfully")
}

// GetAllCourseHandler godoc
// @Summary Get all courses
// @Description Get all courses with pagination and filters
// @Tags courses
// @Accept json
// @Produce json
// @Success 200 {object} utils.StandardResponse{data=dto.PaginationResponse{data=[]dto.CourseResponse}}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/courses [get]
func (h *CourseHandler) GetAllCourseHandler(c *gin.Context) {
	var params dto.CourseQueryParams

	if err := c.ShouldBindQuery(&params); err != nil {
		utils.JSONError(c, "Invalid query parameters", http.StatusBadRequest, err.Error())
		return
	}

	courses, err := h.service.GetAllCourse(c.Request.Context(), &params)
	if err != nil {
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONSuccess(c, courses, "Courses retrieved successfully")
}

// GetByIDCourseHandler godoc
// @Summary Get course by ID
// @Description Retrieve detailed information of a course by its ID
// @Tags courses
// @Accept json
// @Produce json
// @Param id path string true "Course ID (UUID)"
// @Success 200 {object} utils.StandardResponse{data=dto.CourseDetailResponse}
// @Failure 400 {object} utils.ErrorResponse "Invalid course ID"
// @Failure 404 {object} utils.ErrorResponse "Course not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/courses/{id} [get]
func (h *CourseHandler) GetByIDCourseHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.JSONError(c, "Invalid course ID", http.StatusBadRequest, err.Error())
		return
	}

	includeDeleted := c.Query("include_deleted") == "true"

	course, err := h.service.GetByIDCourse(c.Request.Context(), id, includeDeleted)
	if err != nil {
		if err.Error() == "course not found" {
			utils.JSONNotFound(c, "Course not found")
			return
		}
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONSuccess(c, course, "Course retrieved successfully")
}

// UpdateCourseHandler godoc
// @Summary Update a course
// @Description Update course information, optionally including thumbnail upload
// @Tags courses
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Course ID (UUID)"
// @Param name formData string false "Course name"
// @Param description formData string false "Course description"
// @Param price formData number false "Course price"
// @Param is_progress_limited formData boolean false "Limit user progress (optional)"
// @Param course_type_id formData string false "Course type ID (UUID)"
// @Param thumbnail formData file false "Course thumbnail (optional)"
// @Success 200 {object} utils.StandardResponse{data=dto.CourseDetailResponse}
// @Failure 400 {object} utils.ErrorResponse "Invalid request data"
// @Failure 404 {object} utils.ErrorResponse "Course not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/courses/{id} [put]
// @Security BearerAuth
func (h *CourseHandler) UpdateCourseHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.JSONError(c, "Invalid course ID", http.StatusBadRequest, err.Error())
		return
	}

	var req dto.UpdateCourseRequest
	if err := c.ShouldBind(&req); err != nil {
		utils.JSONError(c, "Invalid request data", http.StatusBadRequest, err.Error())
		return
	}

	thumbnail, _ := c.FormFile("thumbnail")

	course, err := h.service.Update(c.Request.Context(), id, &req, thumbnail)
	if err != nil {
		if err.Error() == "course not found" {
			utils.JSONNotFound(c, "Course not found")
			return
		}
		if err.Error() == "course type not found" {
			utils.JSONError(c, "Course type not found", http.StatusBadRequest, nil)
			return
		}
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONSuccess(c, course, "Course updated successfully")
}

// SoftDeleteCourseHandler godoc
// @Summary Soft delete a course
// @Description Soft delete a course (it can be restored later)
// @Tags courses
// @Accept json
// @Produce json
// @Param id path string true "Course ID (UUID)"
// @Success 200 {object} utils.StandardResponse "Course deleted successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid course ID"
// @Failure 404 {object} utils.ErrorResponse "Course not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/courses/{id} [delete]
// @Security BearerAuth
func (h *CourseHandler) SoftDeleteCourseHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.JSONError(c, "Invalid course ID", http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.SoftDeleteCourse(c.Request.Context(), id); err != nil {
		if err.Error() == "course not found" {
			utils.JSONNotFound(c, "Course not found")
			return
		}
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONSuccess(c, nil, "Course deleted successfully")
}

// RestoreCourseHandler godoc
// @Summary Restore a soft deleted course
// @Description Restore a previously soft deleted course
// @Tags courses
// @Accept json
// @Produce json
// @Param id path string true "Course ID"
// @Success 200 {object} utils.StandardResponse{data=dto.CourseDetailResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/courses/{id}/restore [patch]
// @Security BearerAuth
func (h *CourseHandler) RestoreCourseHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.JSONError(c, "Invalid course ID", http.StatusBadRequest, err.Error())
		return
	}

	course, err := h.service.RestoreCourse(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "course not found" {
			utils.JSONNotFound(c, "Course not found")
			return
		}
		if err.Error() == "course is not deleted" {
			utils.JSONError(c, "Course is not deleted", http.StatusBadRequest, nil)
			return
		}
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONSuccess(c, course, "Course restored successfully")
}

// PermanentDeleteCourseHandler godoc
// @Summary Permanently delete a course
// @Description Permanently delete a course (cannot be restored)
// @Tags courses
// @Accept json
// @Produce json
// @Param id path string true "Course ID (UUID)"
// @Success 200 {object} utils.StandardResponse "Course permanently deleted"
// @Failure 400 {object} utils.ErrorResponse "Invalid course ID"
// @Failure 404 {object} utils.ErrorResponse "Course not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/courses/{id}/permanent [delete]
// @Security BearerAuth
func (h *CourseHandler) PermanentDeleteCourseHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.JSONError(c, "Invalid course ID", http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.PermanentDeleteCourse(c.Request.Context(), id); err != nil {
		if err.Error() == "course not found" {
			utils.JSONNotFound(c, "Course not found")
			return
		}
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONSuccess(c, nil, "Course permanently deleted")
}
