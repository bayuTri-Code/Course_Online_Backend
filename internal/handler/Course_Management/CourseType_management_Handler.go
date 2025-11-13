package CourseManagementhandler

import (
	"course_online_backend/internal/dto"
	CourseManagementServices "course_online_backend/internal/services/Course_management_Services"
	"course_online_backend/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CourseTypeHandler struct {
	service *CourseManagementServices.CourseTypeService
}

func NewCourseTypeHandler(service *CourseManagementServices.CourseTypeService) *CourseTypeHandler {
	return &CourseTypeHandler{service: service}
}

func (h *CourseTypeHandler) CreateCourseTypeHandler(c *gin.Context) {
	var req dto.CreateCourseTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, "Invalid request", http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.Create(&req)
	if err != nil {
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONCreated(c, result, "Course type created successfully")
}

func (h *CourseTypeHandler) GetAllCourseTypeHandler(c *gin.Context) {
	result, err := h.service.GetAll()
	if err != nil {
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONSuccess(c, result, "Course types retrieved successfully")
}
