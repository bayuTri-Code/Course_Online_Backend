package dashboard

import (
    "course_online_backend/internal/services"
    "github.com/gin-gonic/gin"
    "net/http"
)

type AdminDashboardHandler struct {
    DashboardService *services.DashboardService
}

func NewAdminDashboardHandler(service *services.DashboardService) *AdminDashboardHandler {
    return &AdminDashboardHandler{DashboardService: service}
}

// @Summary      Get Admin Dashboard
// @Description  Retrieve overview data for Admin, including total users and total courses.
// @Tags         Dashboard
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Success 200 {object} utils.BaseResponsedata{data=dto.AdminDashboardResponse} "Success get admin dashboard"
// @Failure      401  {object}  utils.BaseResponse  "Unauthorized - invalid or missing token"
// @Failure      403  {object}  utils.BaseResponse  "Forbidden - insufficient role"
// @Failure      500  {object}  utils.BaseResponse  "Internal Server Error"
// @Router       /api/dashboard/admin [get]
func (h *AdminDashboardHandler) GetDashboardAdmin(c *gin.Context) {
    data, err := h.DashboardService.GetAdminDashboard()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": data})
}
