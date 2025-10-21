package dashboard

import (
    "course_online_backend/internal/services"
    "github.com/gin-gonic/gin"
    "net/http"
)

type InstructorDashboardHandler struct {
    DashboardService *services.DashboardService
}

func NewInstructorDashboardHandler(service *services.DashboardService) *InstructorDashboardHandler {
    return &InstructorDashboardHandler{DashboardService: service}
}

// @Summary      Get Instructor Dashboard
// @Description  Retrieve dashboard data for Instructor, including number of created courses.
// @Tags         Dashboard
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Success 200 {object} utils.BaseResponsedata{data=dto.InstructorDashboardResponse} "Success get admin dashboard"
// @Failure      401  {object}  utils.BaseResponse  "Unauthorized - invalid or missing token"
// @Failure      403  {object}  utils.BaseResponse  "Forbidden - insufficient role"
// @Failure      500  {object}  utils.BaseResponse  "Internal Server Error"
// @Router       /api/dashboard/instructor [get]
func (h *InstructorDashboardHandler) GetDashboardInstructor(c *gin.Context) {
    instructorID := c.GetUint("user_id") 
    data, err := h.DashboardService.GetInstructorDashboard(instructorID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": data})
}
