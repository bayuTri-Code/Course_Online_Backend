package dashboard

import (
    "course_online_backend/internal/services"
    "github.com/gin-gonic/gin"
    "net/http"
)

type SuperAdminDashboardHandler struct {
    DashboardService *services.DashboardService
}

func NewSuperAdminDashboardHandler(service *services.DashboardService) *SuperAdminDashboardHandler {
    return &SuperAdminDashboardHandler{DashboardService: service}
}

// @Summary      Get Super Admin Dashboard
// @Description  Retrieve overview data for Super Admin, including total instructors and students.
// @Tags         Dashboard
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Success 200 {object} utils.BaseResponsedata{data=dto.SuperAdminDashboardResponse} "Success get admin dashboard"
// @Failure      401  {object}  utils.BaseResponse  "Unauthorized - invalid or missing token"
// @Failure      403  {object}  utils.BaseResponse  "Forbidden - insufficient role"
// @Failure      500  {object}  utils.BaseResponse  "Internal Server Error"
// @Router       /api/dashboard/super_admin [get]
func (h *SuperAdminDashboardHandler) GetDashboardSuperAdmin(c *gin.Context) {
    data, err := h.DashboardService.GetSuperAdminDashboard()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": data})
}
