package handler

import (
	enrollmentdto "course_online_backend/internal/dto/enrollment_dto"
	"course_online_backend/internal/models"
	enrollmentservices "course_online_backend/internal/services/Enrollment_Services"
	"course_online_backend/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EnrollmentHandler struct {
	enrollmentService enrollmentservices.EnrollmentService
}

func NewEnrollmentHandler(enrollmentService enrollmentservices.EnrollmentService) *EnrollmentHandler {
	return &EnrollmentHandler{
		enrollmentService: enrollmentService,
	}
}

// EnrollFreeCourse - POST /api/enrollments/free
func (h *EnrollmentHandler) EnrollFreeCourse(c *gin.Context) {
	// 1. Get user from context
	userInterface, exists := c.Get("user")
	if !exists {
		utils.JSONUnauthorized(c, "User not authenticated")
		return
	}

	user, ok := userInterface.(models.User)
	if !ok {
		utils.JSONError(c, "Invalid user data", http.StatusInternalServerError, nil)
		return
	}

	// 2. Validate student role
	rolesInterface, exists := c.Get("roles")
	if !exists {
		utils.JSONUnauthorized(c, "User roles not found")
		return
	}

	roles, ok := rolesInterface.([]string)
	if !ok {
		utils.JSONError(c, "Invalid roles format", http.StatusInternalServerError, nil)
		return
	}

	isStudent := false
	for _, role := range roles {
		if role == "student" {
			isStudent = true
			break
		}
	}

	if !isStudent {
		utils.JSONError(c, "Only students can enroll in courses", http.StatusForbidden, nil)
		return
	}

	// 3. Bind request to DTO
	var req enrollmentdto.EnrollFreeCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, "Invalid request body", http.StatusBadRequest, err.Error())
		return
	}

	// 4. Parse course ID
	courseID, err := uuid.Parse(req.CourseID)
	if err != nil {
		utils.JSONError(c, "Invalid course ID format", http.StatusBadRequest, nil)
		return
	}

	// 5. Call service
	enrollment, err := h.enrollmentService.EnrollFreeCourse(user.ID, courseID)
	if err != nil {
		switch err.Error() {
		case "course not found":
			utils.JSONNotFound(c, "Course not found")
		case "course is not available for enrollment":
			utils.JSONError(c, "Course is not available for enrollment", http.StatusBadRequest, nil)
		case "this course is not free, please use payment endpoint":
			utils.JSONError(c, "This course requires payment", http.StatusPaymentRequired, nil)
		case "enrollment deadline has passed":
			utils.JSONError(c, "Enrollment deadline has passed", http.StatusBadRequest, nil)
		case "course is full, max capacity reached":
			utils.JSONError(c, "Course is full, maximum capacity reached", http.StatusBadRequest, nil)
		case "you are already enrolled in this course":
			utils.JSONError(c, "You are already enrolled in this course", http.StatusConflict, nil)
		default:
			utils.JSONError(c, "Failed to enroll in course", http.StatusInternalServerError, err.Error())
		}
		return
	}

	// 6. Transform to DTO response
	response := enrollmentdto.ToEnrollmentResponse(enrollment)
	utils.JSONCreated(c, response, "Successfully enrolled in course")
}

// CheckEnrollment - GET /api/enrollments/check/:courseID
func (h *EnrollmentHandler) CheckEnrollment(c *gin.Context) {
	// 1. Get user
	userInterface, exists := c.Get("user")
	if !exists {
		utils.JSONUnauthorized(c, "User not authenticated")
		return
	}

	user, ok := userInterface.(models.User)
	if !ok {
		utils.JSONError(c, "Invalid user data", http.StatusInternalServerError, nil)
		return
	}

	// 2. Parse course ID
	courseIDStr := c.Param("courseID")
	courseID, err := uuid.Parse(courseIDStr)
	if err != nil {
		utils.JSONError(c, "Invalid course ID format", http.StatusBadRequest, nil)
		return
	}

	// 3. Check enrollment
	isEnrolled, err := h.enrollmentService.CheckEnrollment(user.ID, courseID)
	if err != nil {
		utils.JSONError(c, "Failed to check enrollment", http.StatusInternalServerError, err.Error())
		return
	}

	// 4. Transform to DTO response
	response := enrollmentdto.ToCheckEnrollmentResponse(isEnrolled, courseID)
	utils.JSONSuccess(c, response, "Enrollment status retrieved")
}

// GetMyEnrollments - GET /api/enrollments/my-courses
func (h *EnrollmentHandler) GetMyEnrollments(c *gin.Context) {
	// 1. Get user
	userInterface, exists := c.Get("user")
	if !exists {
		utils.JSONUnauthorized(c, "User not authenticated")
		return
	}

	user, ok := userInterface.(models.User)
	if !ok {
		utils.JSONError(c, "Invalid user data", http.StatusInternalServerError, nil)
		return
	}

	// 2. Get enrollments
	enrollments, err := h.enrollmentService.GetMyEnrollments(user.ID)
	if err != nil {
		utils.JSONError(c, "Failed to get enrollments", http.StatusInternalServerError, err.Error())
		return
	}

	// 3. Transform to DTO response
	response := enrollmentdto.ToEnrollmentListResponses(enrollments)
	utils.JSONSuccess(c, response, "My enrollments retrieved successfully")
}

// GetEnrollmentDetail - GET /api/enrollments/:id
func (h *EnrollmentHandler) GetEnrollmentDetail(c *gin.Context) {
	// 1. Get user
	userInterface, exists := c.Get("user")
	if !exists {
		utils.JSONUnauthorized(c, "User not authenticated")
		return
	}

	user, ok := userInterface.(models.User)
	if !ok {
		utils.JSONError(c, "Invalid user data", http.StatusInternalServerError, nil)
		return
	}

	// 2. Parse enrollment ID
	enrollmentIDStr := c.Param("id")
	enrollmentID, err := uuid.Parse(enrollmentIDStr)
	if err != nil {
		utils.JSONError(c, "Invalid enrollment ID format", http.StatusBadRequest, nil)
		return
	}

	// 3. Get enrollment detail
	enrollment, err := h.enrollmentService.GetEnrollmentDetail(user.ID, enrollmentID)
	if err != nil {
		if err.Error() == "enrollment not found" {
			utils.JSONNotFound(c, "Enrollment not found")
		} else if err.Error() == "unauthorized to access this enrollment" {
			utils.JSONError(c, "Unauthorized to access this enrollment", http.StatusForbidden, nil)
		} else {
			utils.JSONError(c, "Failed to get enrollment detail", http.StatusInternalServerError, err.Error())
		}
		return
	}

	// 4. Transform to DTO response (no user info for own enrollment)
	response := enrollmentdto.ToEnrollmentDetailResponse(enrollment, false)
	utils.JSONSuccess(c, response, "Enrollment detail retrieved successfully")
}

// UnenrollCourse - DELETE /api/enrollments/:id
func (h *EnrollmentHandler) UnenrollCourse(c *gin.Context) {
	// 1. Get user
	userInterface, exists := c.Get("user")
	if !exists {
		utils.JSONUnauthorized(c, "User not authenticated")
		return
	}

	user, ok := userInterface.(models.User)
	if !ok {
		utils.JSONError(c, "Invalid user data", http.StatusInternalServerError, nil)
		return
	}

	// 2. Validate student role
	rolesInterface, exists := c.Get("roles")
	if !exists {
		utils.JSONUnauthorized(c, "User roles not found")
		return
	}

	roles, ok := rolesInterface.([]string)
	if !ok {
		utils.JSONError(c, "Invalid roles format", http.StatusInternalServerError, nil)
		return
	}

	isStudent := false
	for _, role := range roles {
		if role == "student" {
			isStudent = true
			break
		}
	}

	if !isStudent {
		utils.JSONError(c, "Only students can unenroll from courses", http.StatusForbidden, nil)
		return
	}

	// 3. Parse enrollment ID
	enrollmentIDStr := c.Param("id")
	enrollmentID, err := uuid.Parse(enrollmentIDStr)
	if err != nil {
		utils.JSONError(c, "Invalid enrollment ID format", http.StatusBadRequest, nil)
		return
	}

	// 4. Unenroll
	err = h.enrollmentService.UnenrollCourse(user.ID, enrollmentID)
	if err != nil {
		if err.Error() == "enrollment not found" {
			utils.JSONNotFound(c, "Enrollment not found")
		} else if err.Error() == "unauthorized to unenroll this course" {
			utils.JSONError(c, "Unauthorized to unenroll this course", http.StatusForbidden, nil)
		} else if err.Error() == "cannot unenroll from completed course" {
			utils.JSONError(c, "Cannot unenroll from completed course", http.StatusBadRequest, nil)
		} else {
			utils.JSONError(c, "Failed to unenroll from course", http.StatusInternalServerError, err.Error())
		}
		return
	}

	// 5. Success response
	utils.JSONSuccess(c, nil, "Successfully unenrolled from course")
}

// UpdateEnrollmentStatus - PUT /api/enrollments/:id/status
func (h *EnrollmentHandler) UpdateEnrollmentStatus(c *gin.Context) {
	// 1. Get user
	userInterface, exists := c.Get("user")
	if !exists {
		utils.JSONUnauthorized(c, "User not authenticated")
		return
	}

	user, ok := userInterface.(models.User)
	if !ok {
		utils.JSONError(c, "Invalid user data", http.StatusInternalServerError, nil)
		return
	}

	// 2. Parse enrollment ID
	enrollmentIDStr := c.Param("id")
	enrollmentID, err := uuid.Parse(enrollmentIDStr)
	if err != nil {
		utils.JSONError(c, "Invalid enrollment ID format", http.StatusBadRequest, nil)
		return
	}

	// 3. Bind request to DTO
	var req enrollmentdto.UpdateEnrollmentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, "Invalid request body", http.StatusBadRequest, err.Error())
		return
	}

	// 4. Update status
	err = h.enrollmentService.UpdateEnrollmentStatus(user.ID, enrollmentID, req.Status)
	if err != nil {
		if err.Error() == "enrollment not found" {
			utils.JSONNotFound(c, "Enrollment not found")
		} else if err.Error() == "unauthorized to update this enrollment" {
			utils.JSONError(c, "Unauthorized to update this enrollment", http.StatusForbidden, nil)
		} else if err.Error() == "invalid status, must be: active, dropped, or completed" {
			utils.JSONError(c, err.Error(), http.StatusBadRequest, nil)
		} else {
			utils.JSONError(c, "Failed to update enrollment status", http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.JSONSuccess(c, nil, "Enrollment status updated successfully")
}

func (h *EnrollmentHandler) GetCourseStudents(c *gin.Context) {
	rolesInterface, exists := c.Get("roles")
	if !exists {
		utils.JSONUnauthorized(c, "User roles not found")
		return
	}

	roles, ok := rolesInterface.([]string)
	if !ok {
		utils.JSONError(c, "Invalid roles format", http.StatusInternalServerError, nil)
		return
	}

	courseIDStr := c.Param("courseID")
	courseID, err := uuid.Parse(courseIDStr)
	if err != nil {
		utils.JSONError(c, "Invalid course ID format", http.StatusBadRequest, nil)
		return
	}

	enrollments, err := h.enrollmentService.GetCourseStudents(courseID, roles)
	if err != nil {
		if err.Error() == "only instructors and admins can view course students" {
			utils.JSONError(c, err.Error(), http.StatusForbidden, nil)
		} else {
			utils.JSONError(c, "Failed to get course students", http.StatusInternalServerError, err.Error())
		}
		return
	}

	response := enrollmentdto.ToEnrollmentDetailResponses(enrollments, true)
	utils.JSONSuccess(c, response, "Course students retrieved successfully")
}