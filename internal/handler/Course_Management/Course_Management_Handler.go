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

// Create godoc
// @Summary Create a new course
// @Description Create a new course with optional thumbnail
// @Tags courses
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Course name"
// @Param description formData string true "Course description"
// @Param price formData number true "Course price"
// @Param is_progress_limited formData boolean false "Is progress limited"
// @Param course_type_id formData string true "Course type ID"
// @Param thumbnail formData file false "Course thumbnail"
// @Success 201 {object} utils.StandardResponse{data=dto.CourseDetailResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /courses [post]
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

// GetAll godoc
// @Summary Get all courses
// @Description Get all courses with pagination and filters
// @Tags courses
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param search query string false "Search by name or description"
// @Param course_type_id query string false "Filter by course type ID"
// @Param created_by query string false "Filter by creator ID"
// @Param min_price query number false "Minimum price"
// @Param max_price query number false "Maximum price"
// @Param sort_by query string false "Sort by field" Enums(name, price, created_at)
// @Param sort_order query string false "Sort order" Enums(asc, desc)
// @Param include_deleted query boolean false "Include soft deleted items"
// @Success 200 {object} utils.StandardResponse{data=dto.PaginationResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /courses [get]
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

// GetByID godoc
// @Summary Get course by ID
// @Description Get detailed course information by ID
// @Tags courses
// @Accept json
// @Produce json
// @Param id path string true "Course ID"
// @Param include_deleted query boolean false "Include if soft deleted"
// @Success 200 {object} utils.StandardResponse{data=dto.CourseDetailResponse}
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /courses/{id} [get]
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

// Update godoc
// @Summary Update a course
// @Description Update course information with optional thumbnail
// @Tags courses
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Course ID"
// @Param name formData string false "Course name"
// @Param description formData string false "Course description"
// @Param price formData number false "Course price"
// @Param is_progress_limited formData boolean false "Is progress limited"
// @Param course_type_id formData string false "Course type ID"
// @Param thumbnail formData file false "Course thumbnail"
// @Success 200 {object} utils.StandardResponse{data=dto.CourseDetailResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /courses/{id} [put]
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

// Delete godoc
// @Summary Soft delete a course
// @Description Soft delete a course (can be restored)
// @Tags courses
// @Accept json
// @Produce json
// @Param id path string true "Course ID"
// @Success 200 {object} utils.StandardResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /courses/{id} [delete]
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

// Restore godoc
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
// @Router /courses/{id}/restore [post]
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

// PermanentDelete godoc
// @Summary Permanently delete a course
// @Description Permanently delete a course (cannot be restored)
// @Tags courses
// @Accept json
// @Produce json
// @Param id path string true "Course ID"
// @Success 200 {object} utils.StandardResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /courses/{id}/permanent [delete]
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
