package CourseManagementhandler

import (
	"course_online_backend/internal/dto"
	CourseManagementServices "course_online_backend/internal/services/Course_management_Services"
	"course_online_backend/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type InstructorHandler struct {
	service *CourseManagementServices.InstructorService
}

func NewInstructorHandler(service *CourseManagementServices.InstructorService) *InstructorHandler {
	return &InstructorHandler{
		service: service,
	}
}

// SearchInstructors godoc
// @Summary Search instructors 
// @Description Search active instructors by username or email for course assignment
// @Tags instructors
// @Accept json
// @Produce json
// @Param query query string true "Search query (min 2 characters)"
// @Param limit query int false "Result limit (1-50, default 10)"
// @Success 200 {object} utils.StandardResponse{data=dto.InstructorSearchResponse}
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/instructors/search [get]
// @Security BearerAuth
func (h *InstructorHandler) SearchInstructors(c *gin.Context) {
	var req dto.InstructorSearchRequest
	
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.JSONError(c, "Invalid query parameters", http.StatusBadRequest, err.Error())
		return
	}

	instructors, err := h.service.SearchInstructors(c.Request.Context(), req.Query, req.Limit)
	if err != nil {
		utils.JSONError(c, "Failed to search instructors", http.StatusInternalServerError, err.Error())
		return
	}

	response := dto.InstructorSearchResponse{
		Total:       int64(len(instructors)),
		Instructors: instructors,
	}

	utils.JSONSuccess(c, response, "Instructors retrieved successfully")
}

// GetInstructorByID godoc
// @Summary Get instructor details
// @Description Get detailed information about a specific instructor
// @Tags instructors
// @Accept json
// @Produce json
// @Param id path string true "Instructor ID (UUID)"
// @Success 200 {object} utils.StandardResponse{data=dto.InstructorResponse}
// @Failure 400 {object} utils.ErrorResponse "Invalid ID"
// @Failure 404 {object} utils.ErrorResponse "Instructor not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/instructors/{id} [get]
// @Security BearerAuth
func (h *InstructorHandler) GetInstructorByID(c *gin.Context) {
	idParam := c.Param("instructorId")
	id, err := utils.ParseUUID(idParam)
	if err != nil {
		utils.JSONError(c, "Invalid instructor ID", http.StatusBadRequest, err.Error())
		return
	}

	instructor, err := h.service.GetInstructorByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "instructor not found or inactive" || err.Error() == "user is not an instructor" {
			utils.JSONNotFound(c, err.Error())
			return
		}
		utils.JSONError(c, "Failed to get instructor", http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONSuccess(c, instructor, "Instructor retrieved successfully")
}