package dashboard

import (
    "course_online_backend/internal/services"
    "github.com/gin-gonic/gin"
    "net/http"
)

type StudentDashboardHandler struct {
    DashboardService *services.DashboardService
}

func NewStudentDashboardHandler(service *services.DashboardService) *StudentDashboardHandler {
    return &StudentDashboardHandler{DashboardService: service}
}

// @Summary      Get Student Dashboard
// @Description  Retrieve dashboard data for Student, including enrolled courses and motivational message.
// @Tags         Dashboard
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Success 200 {object} utils.BaseResponsedata{data=dto.StudentDashboardResponse} "Success get admin dashboard"
// @Failure      401  {object}  utils.BaseResponse  "Unauthorized - invalid or missing token"
// @Failure      403  {object}  utils.BaseResponse  "Forbidden - insufficient role"
// @Failure      500  {object}  utils.BaseResponse  "Internal Server Error"
// @Router       /api/dashboard/student [get]
func (h *StudentDashboardHandler) GetDashboardStudent(c *gin.Context) {
    studentID := c.GetUint("user_id") 
    data, err := h.DashboardService.GetStudentDashboard(studentID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": data})
}
