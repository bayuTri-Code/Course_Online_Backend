package CourseManagementhandler

import (
	"course_online_backend/internal/dto"
	"course_online_backend/internal/services"
	CourseManagementServices "course_online_backend/internal/services/Course_management_Services"
	"course_online_backend/internal/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CourseTypeHandler struct {
	service         *CourseManagementServices.CourseTypeService
	ActivityService *services.ActivityService
}

func NewCourseTypeHandler(service *CourseManagementServices.CourseTypeService, act *services.ActivityService) *CourseTypeHandler {
	return &CourseTypeHandler{
		service:         service,
		ActivityService: act,
	}
}

// CreateCourseTypeHandler godoc
// @Summary Create a new course type
// @Description Create a new course type/category
// @Tags course-types
// @Accept json
// @Produce json
// @Param request body dto.CreateCourseTypeRequestSwagger true "Course type details"
// @Success 201 {object} utils.StandardResponse{data=dto.CourseTypeResponse}
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/course-types [post]
// @Security BearerAuth
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

	if firebaseUID, exists := c.Get("user_id"); exists {
		go func() {
			user, err := h.ActivityService.GetUserByFirebaseUID(firebaseUID.(string))
			if err != nil {
				fmt.Printf("User not found for activity logging")
				return
			}
			activityMsg := fmt.Sprintf("Created course type: %s (ID: %s)", req.Name, result.ID.String())
			_ = h.ActivityService.LogActivity(user.ID, activityMsg)
		}()
	}
}

// GetAllCourseTypeHandler godoc
// @Summary Get all course types
// @Description Retrieve all course types/categories
// @Tags course-types
// @Accept json
// @Produce json
// @Success 200 {object} utils.StandardResponse{data=[]dto.CourseTypeResponse}
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/course-types [get]
func (h *CourseTypeHandler) GetAllCourseTypeHandler(c *gin.Context) {
	result, err := h.service.GetAll()
	if err != nil {
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONSuccess(c, result, "Course types retrieved successfully")
}
