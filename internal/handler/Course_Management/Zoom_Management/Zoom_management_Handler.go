package zoomhandler

import (
	"course_online_backend/internal/dto"
	"course_online_backend/internal/services"
	zoom_management_services "course_online_backend/internal/services/Course_management_Services/Zoom_Management_Services"
	"course_online_backend/internal/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ZoomHandler struct {
	zoomService     *zoom_management_services.ZoomService
	activityService *services.ActivityService
}

func NewZoomHandler(zoomService *zoom_management_services.ZoomService, activityService *services.ActivityService) *ZoomHandler {
	return &ZoomHandler{
		zoomService:     zoomService,
		activityService: activityService,
	}
}

// CreateZoomHandler godoc
// @Summary Create a new Zoom meeting
// @Description Create a new Zoom meeting for a course
// @Tags zoom
// @Accept json
// @Produce json
// @Param zoom body dto.CreateZoomRequest true "Zoom data"
// @Success 201 {object} utils.StandardResponse{data=dto.ZoomResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse "Course not found"
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/zoom [post]
// @Security BearerAuth
func (h *ZoomHandler) CreateZoomHandler(c *gin.Context) {
	var req dto.CreateZoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, "Invalid request data", http.StatusBadRequest, err.Error())
		return
	}

	zoom, err := h.zoomService.CreateZoom(c.Request.Context(), &req)
	if err != nil {
		if err.Error() == "course not found" {
			utils.JSONNotFound(c, "Course not found")
			return
		}
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONCreated(c, zoom, "Zoom meeting created successfully")

	if firebaseUID, exists := c.Get("user_id"); exists {
		go func() {
			user, err := h.activityService.GetUserByFirebaseUID(firebaseUID.(string))
			if err != nil {
				fmt.Printf("User not found for activity logging")
				return
			}
			activityMsg := fmt.Sprintf("Created Zoom meeting '%s' for course ID: %s", req.Title, req.CourseID.String())
			_ = h.activityService.LogActivity(user.ID, activityMsg)
		}()
	}
}

// GetZoomByIDHandler godoc
// @Summary Get Zoom by ID
// @Description Retrieve Zoom meeting details by ID
// @Tags zoom
// @Accept json
// @Produce json
// @Param id path string true "Zoom ID (UUID)"
// @Success 200 {object} utils.StandardResponse{data=dto.ZoomResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/zoom/{id} [get]
// @Security BearerAuth
func (h *ZoomHandler) GetZoomByIDHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.JSONError(c, "Invalid zoom ID", http.StatusBadRequest, err.Error())
		return
	}

	zoom, err := h.zoomService.GetZoomByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "zoom not found" {
			utils.JSONNotFound(c, "Zoom not found")
			return
		}
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONSuccess(c, zoom, "Zoom retrieved successfully")
}

// GetZoomsByCourseIDHandler godoc
// @Summary Get all Zoom meetings by Course ID
// @Description Retrieve all Zoom meetings for a specific course
// @Tags zoom
// @Accept json
// @Produce json
// @Param courseId path string true "Course ID (UUID)"
// @Success 200 {object} utils.StandardResponse{data=dto.ZoomListResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/zoom/course/{courseId} [get]
// @Security BearerAuth
func (h *ZoomHandler) GetZoomsByCourseIDHandler(c *gin.Context) {
	courseIDParam := c.Param("courseId")
	courseID, err := uuid.Parse(courseIDParam)
	if err != nil {
		utils.JSONError(c, "Invalid course ID", http.StatusBadRequest, err.Error())
		return
	}

	zooms, err := h.zoomService.GetZoomsByCourseID(c.Request.Context(), courseID)
	if err != nil {
		if err.Error() == "course not found" {
			utils.JSONNotFound(c, "Course not found")
			return
		}
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONSuccess(c, zooms, "Zoom meetings retrieved successfully")
}

// GetUpcomingZoomsByCourseIDHandler godoc
// @Summary Get upcoming Zoom meetings by Course ID
// @Description Retrieve upcoming Zoom meetings for a specific course
// @Tags zoom
// @Accept json
// @Produce json
// @Param courseId path string true "Course ID (UUID)"
// @Success 200 {object} utils.StandardResponse{data=dto.ZoomListResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/zoom/course/{courseId}/upcoming [get]
// @Security BearerAuth
func (h *ZoomHandler) GetUpcomingZoomsByCourseIDHandler(c *gin.Context) {
	courseIDParam := c.Param("courseId")
	courseID, err := uuid.Parse(courseIDParam)
	if err != nil {
		utils.JSONError(c, "Invalid course ID", http.StatusBadRequest, err.Error())
		return
	}

	zooms, err := h.zoomService.GetUpcomingZoomsByCourseID(c.Request.Context(), courseID)
	if err != nil {
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONSuccess(c, zooms, "Upcoming zoom meetings retrieved successfully")
}

// GetAllZoomsHandler godoc
// @Summary Get all Zoom meetings
// @Description Retrieve all Zoom meetings
// @Tags zoom
// @Accept json
// @Produce json
// @Success 200 {object} utils.StandardResponse{data=[]dto.ZoomResponse}
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/zoom [get]
// @Security BearerAuth
func (h *ZoomHandler) GetAllZoomsHandler(c *gin.Context) {
	zooms, err := h.zoomService.GetAllZooms(c.Request.Context())
	if err != nil {
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONSuccess(c, zooms, "Zoom meetings retrieved successfully")
}

// UpdateZoomHandler godoc
// @Summary Update Zoom meeting
// @Description Update Zoom meeting by ID
// @Tags zoom
// @Accept json
// @Produce json
// @Param id path string true "Zoom ID (UUID)"
// @Param zoom body dto.UpdateZoomRequest true "Updated zoom data"
// @Success 200 {object} utils.StandardResponse{data=dto.ZoomResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/zoom/{id} [put]
// @Security BearerAuth
func (h *ZoomHandler) UpdateZoomHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.JSONError(c, "Invalid zoom ID", http.StatusBadRequest, err.Error())
		return
	}

	var req dto.UpdateZoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, "Invalid request data", http.StatusBadRequest, err.Error())
		return
	}

	zoom, err := h.zoomService.UpdateZoom(c.Request.Context(), id, &req)
	if err != nil {
		if err.Error() == "zoom not found" {
			utils.JSONNotFound(c, "Zoom not found")
			return
		}
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONSuccess(c, zoom, "Zoom meeting updated successfully")

	if firebaseUID, exists := c.Get("user_id"); exists {
		go func() {
			user, err := h.activityService.GetUserByFirebaseUID(firebaseUID.(string))
			if err != nil {
				fmt.Printf("User not found for activity logging")
				return
			}
			activityMsg := fmt.Sprintf("Updated Zoom meeting '%s' (ID: %s)", req.Title, id.String())
			_ = h.activityService.LogActivity(user.ID, activityMsg)
		}()
	}
}

// DeleteZoomHandler godoc
// @Summary Delete Zoom meeting
// @Description Delete Zoom meeting by ID
// @Tags zoom
// @Accept json
// @Produce json
// @Param id path string true "Zoom ID (UUID)"
// @Success 200 {object} utils.StandardResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/zoom/{id} [delete]
// @Security BearerAuth
func (h *ZoomHandler) DeleteZoomHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.JSONError(c, "Invalid zoom ID", http.StatusBadRequest, err.Error())
		return
	}

	if err := h.zoomService.DeleteZoom(c.Request.Context(), id); err != nil {
		if err.Error() == "zoom not found" {
			utils.JSONNotFound(c, "Zoom not found")
			return
		}
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONSuccess(c, nil, "Zoom meeting deleted successfully")

	if firebaseUID, exists := c.Get("user_id"); exists {
		go func() {
			user, err := h.activityService.GetUserByFirebaseUID(firebaseUID.(string))
			if err != nil {
				fmt.Printf("User not found for activity logging")
				return
			}
			activityMsg := fmt.Sprintf("Deleted Zoom meeting (ID: %s)", id.String())
			_ = h.activityService.LogActivity(user.ID, activityMsg)
		}()
	}
}