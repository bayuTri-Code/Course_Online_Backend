package CourseManagementhandler

import (
	"course_online_backend/internal/dto"
	"course_online_backend/internal/services"
	CourseManagementServices "course_online_backend/internal/services/Course_management_Services"
	"course_online_backend/internal/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CourseHandler struct {
	service         *CourseManagementServices.CourseService
	ActivityService *services.ActivityService
}

func NewCourseHandler(service *CourseManagementServices.CourseService, act *services.ActivityService) *CourseHandler {
	return &CourseHandler{
		service:         service,
		ActivityService: act,
	}
}

// CreateCourseHandler godoc
// @Summary Create a new course
// @Description Create a new course with optional thumbnail upload and instructor assignment
// @Tags courses
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Course name"
// @Param description formData string true "Course description"
// @Param price formData number true "Course price"
// @Param is_progress_limited formData boolean false "Limit user progress (optional)"
// @Param course_type_id formData string true "Course type ID (UUID)"
// @Param instructor_id formData string false "Instructor ID (UUID) - optional, assign instructor to course"
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
		if err.Error() == "instructor validation failed: instructor not found or inactive" {
			utils.JSONError(c, "Invalid instructor: instructor not found or inactive", http.StatusBadRequest, nil)
			return
		}
		if err.Error() == "instructor validation failed: user is not an instructor" {
			utils.JSONError(c, "Invalid instructor: user does not have instructor role", http.StatusBadRequest, nil)
			return
		}
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONCreated(c, course, "Course created successfully")

	if firebaseUID, exists := c.Get("user_id"); exists {
		go func() {
			user, err := h.ActivityService.GetUserByFirebaseUID(firebaseUID.(string))
			if err != nil {
				fmt.Printf("User not found for activity logging")
				return
			}
			
			activityMsg := fmt.Sprintf("Created course: %s (ID: %s)", req.Name, course.ID.String())
			if req.InstructorID != "" {
				activityMsg += "with instructor assigned"
			}
			
			_ = h.ActivityService.LogActivity(user.ID, activityMsg)
		}()
	}
}

// UpdateCourseHandler godoc
// @Summary Update a course
// @Description Update course information, optionally including thumbnail upload and instructor assignment/removal
// @Tags courses
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Course ID (UUID)"
// @Param name formData string false "Course name"
// @Param description formData string false "Course description"
// @Param price formData number false "Course price"
// @Param is_progress_limited formData boolean false "Limit user progress (optional)"
// @Param course_type_id formData string false "Course type ID (UUID)"
// @Param instructor_id formData string false "Instructor ID (UUID) - empty string to remove instructor"
// @Param thumbnail formData file false "Course thumbnail (optional)"
// @Success 200 {object} utils.StandardResponse{data=dto.CourseDetailResponse}
// @Failure 400 {object} utils.ErrorResponse "Invalid request data"
// @Failure 404 {object} utils.ErrorResponse "Course not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/courses/{id} [put]
// @Security BearerAuth
func (h *CourseHandler) UpdateCourseHandler(c *gin.Context) {
	idParam := c.Param("course_id")
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
		if err.Error() == "instructor validation failed: instructor not found or inactive" {
			utils.JSONError(c, "Invalid instructor: instructor not found or inactive", http.StatusBadRequest, nil)
			return
		}
		if err.Error() == "instructor validation failed: user is not an instructor" {
			utils.JSONError(c, "Invalid instructor: user does not have instructor role", http.StatusBadRequest, nil)
			return
		}
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONSuccess(c, course, "Course updated successfully")

	if firebaseUID, exists := c.Get("user_id"); exists {
		go func() {
			user, err := h.ActivityService.GetUserByFirebaseUID(firebaseUID.(string))
			if err != nil {
				fmt.Printf("User not found for activity logging")
				return
			}
			
			courseName := course.Name
			if req.Name != "" {
				courseName = req.Name
			}
			
			activityMsg := fmt.Sprintf("Updated course: %s (ID: %s)", courseName, id.String())
			if req.InstructorID != nil {
				if *req.InstructorID == "" {
					activityMsg += " - removed instructor"
				} else {
					activityMsg += " - changed instructor"
				}
			}
			
			_ = h.ActivityService.LogActivity(user.ID, activityMsg)
		}()
	}
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
	idParam := c.Param("course_id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.JSONError(c, "Invalid course ID", http.StatusBadRequest, err.Error())
		return
	}

	course, err := h.service.GetByIDCourse(c.Request.Context(), id, false)
	if err != nil {
		if err.Error() == "course not found" {
			utils.JSONNotFound(c, "Course not found")
			return
		}
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
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

	if firebaseUID, exists := c.Get("user_id"); exists {
		go func() {
			user, err := h.ActivityService.GetUserByFirebaseUID(firebaseUID.(string))
			if err != nil {
				fmt.Printf("User not found for activity logging")
				return
			}
			activityMsg := fmt.Sprintf("Soft deleted course: %s (ID: %s)", course.Name, id.String())
			_ = h.ActivityService.LogActivity(user.ID, activityMsg)
		}()
	}
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
	idParam := c.Param("course_id")
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

	if firebaseUID, exists := c.Get("user_id"); exists {
		go func() {
			user, err := h.ActivityService.GetUserByFirebaseUID(firebaseUID.(string))
			if err != nil {
				fmt.Printf("User not found for activity logging")
				return
			}
			activityMsg := fmt.Sprintf("Restored course: %s (ID: %s)", course.Name, id.String())
			_ = h.ActivityService.LogActivity(user.ID, activityMsg)
		}()
	}
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
	idParam := c.Param("course_id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.JSONError(c, "Invalid course ID", http.StatusBadRequest, err.Error())
		return
	}

	course, err := h.service.GetByIDCourse(c.Request.Context(), id, true)
	if err != nil {
		if err.Error() == "course not found" {
			utils.JSONNotFound(c, "Course not found")
			return
		}
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
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

	if firebaseUID, exists := c.Get("user_id"); exists {
		go func() {
			user, err := h.ActivityService.GetUserByFirebaseUID(firebaseUID.(string))
			if err != nil {
				fmt.Printf("User not found for activity logging")
				return
			}
			activityMsg := fmt.Sprintf("Permanently deleted course: %s (ID: %s)", course.Name, id.String())
			_ = h.ActivityService.LogActivity(user.ID, activityMsg)
		}()
	}
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
	idParam := c.Param("course_id")
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

// GetCoursesByCategoryHandler godoc
// @Summary Get courses by category
// @Description Get all courses in a specific category with pagination
// @Tags courses-browsing
// @Accept json
// @Produce json
// @Param categoryId path string true "Category ID (UUID)"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} utils.StandardResponse{data=dto.PaginationResponse{data=[]dto.CourseResponse}}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/courses/categories/{categoryId} [get]
func (h *CourseHandler) GetCoursesByCategoryHandler(c *gin.Context) {
	categoryIDParam := c.Param("categoryId")
	categoryID, err := uuid.Parse(categoryIDParam)
	if err != nil {
		utils.JSONError(c, "Invalid category ID", http.StatusBadRequest, err.Error())
		return
	}

	var params dto.SimplePaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		utils.JSONError(c, "Invalid query parameters", http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.GetCoursesByCategory(c.Request.Context(), categoryID, &params)
	if err != nil {
		if err.Error() == "category not found" {
			utils.JSONNotFound(c, "Category not found")
			return
		}
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONSuccess(c, result, "Courses retrieved successfully")
}


// GetMyCoursesByCategoryHandler godoc
// @Summary Get my enrolled courses by category
// @Description Get courses that the user has enrolled in, filtered by category
// @Tags courses-browsing
// @Accept json
// @Produce json
// @Param categoryId path string true "Category ID (UUID)"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} utils.StandardResponse{data=dto.PaginationResponse{data=[]dto.CourseResponse}}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/mycourses/categories/:categoryId/courses [get]
// @Security BearerAuth
func (h *CourseHandler) GetMyCoursesByCategoryHandler(c *gin.Context) {
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

	categoryIDParam := c.Param("categoryId")
	categoryID, err := uuid.Parse(categoryIDParam)
	if err != nil {
		utils.JSONError(c, "Invalid category ID", http.StatusBadRequest, err.Error())
		return
	}

	var params dto.SimplePaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		utils.JSONError(c, "Invalid query parameters", http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.GetMyCoursesByCategory(c.Request.Context(), userID, categoryID, &params)
	if err != nil {
		if err.Error() == "user not found" {
			utils.JSONNotFound(c, "User not found")
			return
		}
		if err.Error() == "category not found" {
			utils.JSONNotFound(c, "Category not found")
			return
		}
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONSuccess(c, result, "My enrolled courses by category retrieved successfully")
}

// GetMyAssignedCoursesHandler godoc
// @Summary Get my assigned courses (for instructor)
// @Description Get all courses assigned to the logged-in instructor
// @Tags courses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} utils.StandardResponse{data=dto.MyCoursesInstructorResponse}
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/instructor/my-courses [get]
func (h *CourseHandler) GetMyAssignedCoursesHandler(c *gin.Context) {
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

	instructorID, err := h.service.GetInstructorIDByFirebaseUID(c.Request.Context(), firebaseUID)
	if err != nil {
		if err.Error() == "instructor not found" {
			utils.JSONNotFound(c, "Instructor not found")
			return
		}
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	var params dto.SimplePaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		utils.JSONError(c, "Invalid query parameters", http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.GetMyAssignedCourses(c.Request.Context(), instructorID, &params)
	if err != nil {
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONSuccess(c, result, "My assigned courses retrieved successfully")
}

// GetPopularCoursesHandler godoc
// @Summary Get popular courses
// @Description Get most popular courses sorted by enrollment count
// @Tags courses-browsing
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} utils.StandardResponse{data=dto.PaginationResponse{data=[]dto.CourseResponse}}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/courses/popular [get]
func (h *CourseHandler) GetPopularCoursesHandler(c *gin.Context) {
	var params dto.SimplePaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		utils.JSONError(c, "Invalid query parameters", http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.GetPopularCourses(c.Request.Context(), &params)
	if err != nil {
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONSuccess(c, result, "Popular courses retrieved successfully")
}

// GetLatestCoursesHandler godoc
// @Summary Get latest courses
// @Description Get newest courses sorted by creation date
// @Tags courses-browsing
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} utils.StandardResponse{data=dto.PaginationResponse{data=[]dto.CourseResponse}}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/courses/latest [get]
func (h *CourseHandler) GetLatestCoursesHandler(c *gin.Context) {
	var params dto.SimplePaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		utils.JSONError(c, "Invalid query parameters", http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.GetLatestCourses(c.Request.Context(), &params)
	if err != nil {
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONSuccess(c, result, "Latest courses retrieved successfully")
}


// GetAllCourseTypesHandler godoc
// @Summary Get all course types
// @Description Get all course categories with course count
// @Tags courses-browsing
// @Accept json
// @Produce json
// @Success 200 {object} utils.StandardResponse{data=[]dto.CourseTypeWithCountResponse}
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/courses/course-types [get]
func (h *CourseHandler) GetAllCourseTypesHandler(c *gin.Context) {
	courseTypes, err := h.service.GetAllCourseTypes(c.Request.Context())
	if err != nil {
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONSuccess(c, courseTypes, "Course types retrieved successfully")
}



// GetRelatedCoursesHandler godoc
// @Summary Get related courses
// @Description Get courses related to a specific course (same category)
// @Tags courses-browsing
// @Accept json
// @Produce json
// @Param id path string true "Course ID (UUID)"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(5)
// @Success 200 {object} utils.StandardResponse{data=dto.PaginationResponse{data=[]dto.CourseResponse}}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/courses/{id}/related [get]
func (h *CourseHandler) GetRelatedCoursesHandler(c *gin.Context) {
	idParam := c.Param("course_id")
	courseID, err := uuid.Parse(idParam)
	if err != nil {
		utils.JSONError(c, "Invalid course ID", http.StatusBadRequest, err.Error())
		return
	}

	var params dto.SimplePaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		utils.JSONError(c, "Invalid query parameters", http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.GetRelatedCourses(c.Request.Context(), courseID, &params)
	if err != nil {
		if err.Error() == "course not found" {
			utils.JSONNotFound(c, "Course not found")
			return
		}
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONSuccess(c, result, "Related courses retrieved successfully")
}

// GetCourseStatsHandler godoc
// @Summary Get course statistics
// @Description Get overall course statistics including total courses, enrollments, and category breakdown
// @Tags courses-browsing
// @Accept json
// @Produce json
// @Success 200 {object} utils.StandardResponse{data=dto.CourseStatsResponse}
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/courses/stats [get]
// @Security BearerAuth
func (h *CourseHandler) GetCourseStatsHandler(c *gin.Context) {
	stats, err := h.service.GetCourseStats(c.Request.Context())
	if err != nil {
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONSuccess(c, stats, "Course statistics retrieved successfully")
}
// GetMyCoursesHandler godoc
// @Summary Get my courses
// @Description Get courses based on user role: Super Admin (all courses), Admin (created courses), Instructor (assigned courses), Student (enrolled courses)
// @Tags courses-browsing
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param search query string false "Search keyword"
// @Param course_type_id query string false "Filter by course type"
// @Param sort_by query string false "Sort field" Enums(name, price, created_at)
// @Param sort_order query string false "Sort order" Enums(asc, desc)
// @Param include_deleted query bool false "Include deleted courses"
// @Success 200 {object} utils.StandardResponse{data=dto.PaginationResponse{data=[]dto.CourseResponse}}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/courses/my-courses [get]
// @Security BearerAuth
func (h *CourseHandler) GetMyCoursesHandler(c *gin.Context) {
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

	roleValue, exists := c.Get("roles")
	if !exists {
		utils.JSONUnauthorized(c, "User role not found")
		return
	}

	var userRole string
	switch roles := roleValue.(type) {
	case string:
		userRole = roles
	case []string:
		if len(roles) > 0 {
			for _, r := range []string{"super_admin", "admin", "instructor", "student"} {
				for _, role := range roles {
					if role == r {
						userRole = r
						break
					}
				}
				if userRole != "" {
					break
				}
			}
		}
	default:
		utils.JSONError(c, "Invalid role format", http.StatusBadRequest, nil)
		return
	}

	var params dto.CourseQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		utils.JSONError(c, "Invalid query parameters", http.StatusBadRequest, err.Error())
		return
	}

	courses, err := h.service.GetMyCourses(c.Request.Context(), userID, userRole, &params)
	if err != nil {
		if err.Error() == "user not found" {
			utils.JSONNotFound(c, "User not found")
			return
		}
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONSuccess(c, courses, "My courses retrieved successfully")
}