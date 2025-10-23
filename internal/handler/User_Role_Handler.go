package handler

import (
	"course_online_backend/internal/models"
	"course_online_backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserRoleHandler struct {
	Service *services.UserRoleService
}

func NewUserRoleHandler(service *services.UserRoleService) *UserRoleHandler {
	return &UserRoleHandler{Service: service}
}

type AssignRoleRequest struct {
	RoleID uuid.UUID `json:"role_id" binding:"required"`
}

// @Summary Assign role to a user
// @Description Assign a role to a specific user based on current user's permissions
// @Tags Roles
// @Accept json
// @Produce json
// @Param user_id path string true "Target User ID"
// @Param request body AssignRoleRequest true "Role ID to assign"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /admin/users/{user_id}/role [put]
// @Security BearerAuth
func (h *UserRoleHandler) AssignRole(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	currentUser := user.(models.User)

	userIDParam := ctx.Param("user_id")
	targetUserID, err := uuid.Parse(userIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var req AssignRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.Service.AssignRole(targetUserID, req.RoleID, currentUser.ID)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "role assigned successfully",
		"user_id": targetUserID,
		"role_id": req.RoleID,
	})
}


// @Summary Get all roles
// @Description Retrieve all roles in the system (super_admin & admin only)
// @Tags Roles
// @Produce json
// @Success 200 {array} models.Role
// @Router /roles [get]
// @Security BearerAuth
func (h *UserRoleHandler) GetAllRoles(c *gin.Context) {
	roles, err := h.Service.GetAllRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve roles"})
		return
	}
	c.JSON(http.StatusOK, roles)
}

// @Summary Get assignable roles
// @Description Get roles that the current user is allowed to assign
// @Tags Roles
// @Produce json
// @Success 200 {array} models.Role
// @Router /roles/assignable [get]
// @Security BearerAuth
func (h *UserRoleHandler) GetAssignableRoles(c *gin.Context) {
	user, _ := c.Get("user")
	currentUser := user.(models.User)

	roles, err := h.Service.GetAssignableRoles(currentUser)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, roles)
}

// @Summary Get role by ID
// @Description Retrieve a single role by its ID
// @Tags Roles
// @Produce json
// @Param id path string true "Role ID"
// @Success 200 {object} models.Role
// @Failure 404 {object} map[string]string
// @Router /roles/{id} [get]
// @Security BearerAuth
func (h *UserRoleHandler) GetRoleByID(c *gin.Context) {
	roleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role ID"})
		return
	}

	role, err := h.Service.GetRoleByID(roleID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "role not found"})
		return
	}
	c.JSON(http.StatusOK, role)
}
