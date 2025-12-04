package handlers

import (
	"net/http"

	enrollmentdto "course_online_backend/internal/dto/enrollment_dto"
	"course_online_backend/internal/models"
	enrollmentservices "course_online_backend/internal/services/Enrollment_Services"
	"course_online_backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EnrollmentHandler struct {
	service *enrollmentservices.EnrollmentService
}

func NewEnrollmentHandler(service *enrollmentservices.EnrollmentService) *EnrollmentHandler {
	return &EnrollmentHandler{service: service}
}

// EnrollCourse godoc
// @Summary Enroll in a free course
// @Description Student enrolls in a free course. Only courses with status "published", free price, not full, and not past enrollment deadline can be enrolled.
// @Tags Enrollment
// @Accept json
// @Produce json
// @Param body body enrollmentdto.EnrollFreeCourseRequest true "Enrollment Request"
// @Success 201 {object} enrollmentdto.EnrollmentResponse "Successfully enrolled"
// @Failure 400 {object} utils.ErrorResponse "Invalid course ID, invalid request body, course not available, deadline passed, or course full"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized - user not authenticated"
// @Failure 402 {object} utils.ErrorResponse "Payment required - course is not free"
// @Failure 403 {object} utils.ErrorResponse "Only students can enroll"
// @Failure 404 {object} utils.ErrorResponse "Course not found"
// @Failure 409 {object} utils.ErrorResponse "User already enrolled"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /enrollments [post]
func (h *EnrollmentHandler) EnrollCourseHandler(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		utils.JSONError(c, "user not authenticated", http.StatusUnauthorized, nil)
		return
	}

	user, ok := userInterface.(models.User)
	if !ok {
		utils.JSONError(c, "invalid user data", http.StatusInternalServerError, nil)
		return
	}

	rolesInterface, exists := c.Get("roles")
	if !exists {
		utils.JSONError(c, "user roles not found", http.StatusUnauthorized, nil)
		return
	}

	roles, ok := rolesInterface.([]string)
	if !ok {
		utils.JSONError(c, "invalid roles format", http.StatusInternalServerError, nil)
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
		utils.JSONError(c, "only students can enroll in courses", http.StatusForbidden, nil)
		return
	}

	var req enrollmentdto.EnrollFreeCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, "invalid request body", http.StatusBadRequest, err.Error())
		return
	}

	courseID, err := uuid.Parse(req.CourseID)
	if err != nil {
		utils.JSONError(c, "invalid course_id format", http.StatusBadRequest, err.Error())
		return
	}

	enrollment, err := h.service.EnrollCourse(user.ID, courseID)
	if err != nil {
		switch err.Error() {
		case "course not found":
			utils.JSONError(c, "course not found", http.StatusNotFound, nil)
		case "course is not available for enrollment":
			utils.JSONError(c, "course is not available for enrollment", http.StatusBadRequest, nil)
		case "this course is not free, please use payment endpoint":
			utils.JSONError(c, "this course requires payment", http.StatusPaymentRequired, nil)
		case "enrollment deadline has passed":
			utils.JSONError(c, "enrollment deadline has passed", http.StatusBadRequest, nil)
		case "course is full, max capacity reached":
			utils.JSONError(c, "course is full, maximum capacity reached", http.StatusBadRequest, nil)
		case "you are already enrolled in this course":
			utils.JSONError(c, "you are already enrolled in this course", http.StatusConflict, nil)
		default:
			utils.JSONError(c, "failed to enroll in course", http.StatusInternalServerError, err.Error())
		}
		return
	}

	response := enrollmentdto.ToEnrollmentResponse(enrollment)
	utils.JSONSuccess(c, response, "successfully enrolled in course")
}

// CheckEnrollment godoc
// @Summary Check enrollment status
// @Description Check whether the authenticated student is enrolled in the specified course
// @Tags Enrollment
// @Accept json
// @Produce json
// @Param id path string true "Course ID"
// @Success 200 {object} enrollmentdto.CheckEnrollmentResponse "Enrollment status retrieved"
// @Failure 400 {object} utils.ErrorResponse "Invalid course ID format"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 500 {object} utils.ErrorResponse "Failed to check enrollment"
// @Security BearerAuth
// @Router /courses/{id}/enrollments/check [get]
func (h *EnrollmentHandler) CheckEnrollmentHandler(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		utils.JSONError(c, "user not authenticated", http.StatusUnauthorized, nil)
		return
	}

	user, ok := userInterface.(models.User)
	if !ok {
		utils.JSONError(c, "invalid user data", http.StatusInternalServerError, nil)
		return
	}

	courseIDStr := c.Param("course_id")
	courseID, err := uuid.Parse(courseIDStr)
	if err != nil {
		utils.JSONError(c, "invalid course_id format", http.StatusBadRequest, err.Error())
		return
	}

	isEnrolled, err := h.service.CheckEnrollment(user.ID, courseID)
	if err != nil {
		utils.JSONError(c, "failed to check enrollment", http.StatusInternalServerError, err.Error())
		return
	}

	response := enrollmentdto.ToCheckEnrollmentResponse(isEnrolled, courseID)
	utils.JSONSuccess(c, response, "enrollment status retrieved")
}

// GetMyEnrollments godoc
// @Summary Get student's enrollments
// @Description Retrieve all enrollments belonging to the authenticated student
// @Tags Enrollment
// @Accept json
// @Produce json
// @Success 200 {array} enrollmentdto.EnrollmentListResponse "List of student enrollments"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 500 {object} utils.ErrorResponse "Failed to retrieve enrollments"
// @Security BearerAuth
// @Router /enrollments [get]
func (h *EnrollmentHandler) GetMyEnrollmentsHandler(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		utils.JSONError(c, "user not authenticated", http.StatusUnauthorized, nil)
		return
	}

	user, ok := userInterface.(models.User)
	if !ok {
		utils.JSONError(c, "invalid user data", http.StatusInternalServerError, nil)
		return
	}

	enrollments, err := h.service.GetMyEnrollments(user.ID)
	if err != nil {
		utils.JSONError(c, "failed to get enrollments", http.StatusInternalServerError, err.Error())
		return
	}

	response := enrollmentdto.ToEnrollmentListResponses(enrollments)
	utils.JSONSuccess(c, response, "my enrollments retrieved successfully")
}

// GetEnrollmentDetail godoc
// @Summary Get enrollment detail
// @Description Retrieve specific enrollment detail. Only the owner of the enrollment can access this data.
// @Tags Enrollment
// @Accept json
// @Produce json
// @Param id path string true "Enrollment ID"
// @Success 200 {object} enrollmentdto.EnrollmentDetailResponse "Detailed enrollment information"
// @Failure 400 {object} utils.ErrorResponse "Invalid enrollment ID"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 403 {object} utils.ErrorResponse "Forbidden - cannot access others' enrollment"
// @Failure 404 {object} utils.ErrorResponse "Enrollment not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /enrollments/{id} [get]
func (h *EnrollmentHandler) GetEnrollmentDetailHandler(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		utils.JSONError(c, "user not authenticated", http.StatusUnauthorized, nil)
		return
	}

	user, ok := userInterface.(models.User)
	if !ok {
		utils.JSONError(c, "invalid user data", http.StatusInternalServerError, nil)
		return
	}

	enrollmentIDStr := c.Param("enrollment_id")
	enrollmentID, err := uuid.Parse(enrollmentIDStr)
	if err != nil {
		utils.JSONError(c, "invalid enrollment_id format", http.StatusBadRequest, err.Error())
		return
	}

	roles := c.MustGet("roles").([]string)

	enrollment, err := h.service.GetEnrollmentDetail(user.ID, enrollmentID, roles)
	if err != nil {
		switch err.Error() {
		case "enrollment not found":
			utils.JSONError(c, "enrollment not found", http.StatusNotFound, nil)
		case "unauthorized to access this enrollment":
			utils.JSONError(c, "unauthorized to access this enrollment", http.StatusForbidden, nil)
		default:
			utils.JSONError(c, "failed to get enrollment detail", http.StatusInternalServerError, err.Error())
		}
		return
	}

	response := enrollmentdto.ToEnrollmentDetailResponse(enrollment, false)
	utils.JSONSuccess(c, response, "enrollment detail retrieved successfully")
}

// UnenrollCourse godoc
// @Summary Unenroll from a course
// @Description Student removes themselves from an active enrollment. Cannot unenroll if course is completed.
// @Tags Enrollment
// @Accept json
// @Produce json
// @Param id path string true "Enrollment ID"
// @Success 200 {object} utils.StandardResponse "Successfully unenrolled"
// @Failure 400 {object} utils.ErrorResponse "Invalid enrollment ID or course already completed"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 403 {object} utils.ErrorResponse "Only students can unenroll or unauthorized to unenroll"
// @Failure 404 {object} utils.ErrorResponse "Enrollment not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /enrollments/{id} [delete]
func (h *EnrollmentHandler) UnenrollCourseHandler(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		utils.JSONError(c, "user not authenticated", http.StatusUnauthorized, nil)
		return
	}

	user, ok := userInterface.(models.User)
	if !ok {
		utils.JSONError(c, "invalid user data", http.StatusInternalServerError, nil)
		return
	}

	rolesInterface, exists := c.Get("roles")
	if !exists {
		utils.JSONError(c, "user roles not found", http.StatusUnauthorized, nil)
		return
	}

	roles, ok := rolesInterface.([]string)
	if !ok {
		utils.JSONError(c, "invalid roles format", http.StatusInternalServerError, nil)
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
		utils.JSONError(c, "only students can unenroll from courses", http.StatusForbidden, nil)
		return
	}

	enrollmentIDStr := c.Param("enrollment_id")
	enrollmentID, err := uuid.Parse(enrollmentIDStr)
	if err != nil {
		utils.JSONError(c, "invalid enrollment_id format", http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.UnenrollCourse(user.ID, enrollmentID)
	if err != nil {
		switch err.Error() {
		case "enrollment not found":
			utils.JSONError(c, "enrollment not found", http.StatusNotFound, nil)
		case "unauthorized to unenroll this course":
			utils.JSONError(c, "unauthorized to unenroll this course", http.StatusForbidden, nil)
		case "cannot unenroll from completed course":
			utils.JSONError(c, "cannot unenroll from completed course", http.StatusBadRequest, nil)
		default:
			utils.JSONError(c, "failed to unenroll from course", http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.JSONSuccess(c, nil, "successfully unenrolled from course")
}

// UpdateEnrollmentStatus godoc
// @Summary Update enrollment status
// @Description Update enrollment status to active, dropped, or completed. Only authorized roles (admin, super_admin, instructor) may update.
// @Tags Enrollment
// @Accept json
// @Produce json
// @Param id path string true "Enrollment ID"
// @Param body body enrollmentdto.UpdateEnrollmentStatusRequest true "Status Update Request"
// @Success 200 {object} utils.StandardResponse "Status updated successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid request body or invalid status"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 403 {object} utils.ErrorResponse "Forbidden - insufficient role or unauthorized user"
// @Failure 404 {object} utils.ErrorResponse "Enrollment not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /enrollments/{id}/status [put]
func (h *EnrollmentHandler) UpdateEnrollmentStatusHandler(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		utils.JSONError(c, "user not authenticated", http.StatusUnauthorized, nil)
		return
	}

	user, ok := userInterface.(models.User)
	if !ok {
		utils.JSONError(c, "invalid user data", http.StatusInternalServerError, nil)
		return
	}

	enrollmentIDStr := c.Param("enrollment_id")
	enrollmentID, err := uuid.Parse(enrollmentIDStr)
	if err != nil {
		utils.JSONError(c, "invalid enrollment_id format", http.StatusBadRequest, err.Error())
		return
	}

	var req enrollmentdto.UpdateEnrollmentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, "invalid request body", http.StatusBadRequest, err.Error())
		return
	}

	roles := c.MustGet("roles").([]string)


	err = h.service.UpdateEnrollmentStatus(user.ID, enrollmentID, req.Status, roles)
	if err != nil {
		switch err.Error() {
		case "enrollment not found":
			utils.JSONError(c, "enrollment not found", http.StatusNotFound, nil)
		case "unauthorized to update this enrollment":
			utils.JSONError(c, "unauthorized to update this enrollment", http.StatusForbidden, nil)
		case "invalid status, must be: active, dropped, or completed":
			utils.JSONError(c, err.Error(), http.StatusBadRequest, nil)
		default:
			utils.JSONError(c, "failed to update enrollment status", http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.JSONSuccess(c, nil, "enrollment status updated successfully")
}

// GetCourseStudents godoc
// @Summary Get course students
// @Description Admin/Instructor gets list of students enrolled in a course
// @Tags Enrollment
// @Accept json
// @Produce json
// @Param id path string true "Course ID"
// @Success 200 {array} enrollmentdto.EnrollmentDetailResponse "List of enrolled students (detail responses)"
// @Failure 400 {object} utils.ErrorResponse "Invalid course ID"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 403 {object} utils.ErrorResponse "Forbidden"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /courses/{id}/enrollments [get]
func (h *EnrollmentHandler) GetCourseStudentsHandler(c *gin.Context) {
	rolesInterface, exists := c.Get("roles")
	if !exists {
		utils.JSONError(c, "user roles not found", http.StatusUnauthorized, nil)
		return
	}

	roles, ok := rolesInterface.([]string)
	if !ok {
		utils.JSONError(c, "invalid roles format", http.StatusInternalServerError, nil)
		return
	}

	courseIDStr := c.Param("course_id")
	courseID, err := uuid.Parse(courseIDStr)
	if err != nil {
		utils.JSONError(c, "invalid course_id format", http.StatusBadRequest, err.Error())
		return
	}

	enrollments, err := h.service.GetCourseStudents(courseID, roles)
	if err != nil {
		if err.Error() == "only instructors and admins can view course students" {
			utils.JSONError(c, err.Error(), http.StatusForbidden, nil)
		} else {
			utils.JSONError(c, "failed to get course students", http.StatusInternalServerError, err.Error())
		}
		return
	}

	response := enrollmentdto.ToEnrollmentDetailResponses(enrollments, true)
	utils.JSONSuccess(c, response, "course students retrieved successfully")
}