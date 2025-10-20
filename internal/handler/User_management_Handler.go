package handler

import (
	"course_online_backend/internal/dto"
	"course_online_backend/internal/models"
	"course_online_backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserService *services.UserService
}

func NewUserHandler(us *services.UserService) *UserHandler {
	return &UserHandler{UserService: us}
}

// GetAllUsers godoc
// @Summary Get all users
// @Description Retrieve a list of all users with their basic info, biodata, and roles
// @Tags Users
// @Produce json
// @Success 200 {object} dto.BaseResponse{data=[]models.User} "All users retrieved successfully"
// @Failure 500 {object} utils.ErrorResponse "Failed to get users"
// @Router /api/user [get]
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.UserService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{
			Status:  false,
			Message: "failed to get users",
		})
		return
	}

	c.JSON(http.StatusOK, dto.BaseResponse{
		Status:  true,
		Message: "success",
		Data:    users,
	})
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Retrieve user details including biodata, roles, and other relations
// @Tags Users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} dto.BaseResponse{data=models.User} "User retrieved successfully"
// @Failure 404 {object} dto.BaseResponse "User not found"
// @Router /api/user/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")

	user, err := h.UserService.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.BaseResponse{
			Status:  false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.BaseResponse{
		Status:  true,
		Message: "success",
		Data:    user,
	})
}

// UpdateUser godoc
// @Summary Update user
// @Description Update user info, roles, and biodata by ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body dto.UpdateUserRequest true "User data to update"
// @Success 200 {object} dto.BaseResponse{data=models.User} "User updated successfully"
// @Failure 400 {object} dto.BaseResponse "Invalid input"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/user/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.BaseResponse{
			Status:  false,
			Message: "invalid input",
		})
		return
	}

	updatedUser := models.User{
		Username:     req.Username,
		EmailAddress: req.EmailAddress,
		IsActive:     *req.IsActive,
	}

	if len(req.RoleIDs) > 0 {
		var roles []models.Role
		h.UserService.DB.Find(&roles, "id IN ?", req.RoleIDs)
		updatedUser.Roles = roles
	}

	user, err := h.UserService.UpdateUser(id, updatedUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{
			Status:  false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.BaseResponse{
		Status:  true,
		Message: "user updated successfully",
		Data:    user,
	})
}


// DeleteUser godoc
// @Summary Delete user
// @Description Delete a user by ID
// @Tags Users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} dto.BaseResponseDelete "User deleted successfully"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/user/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := h.UserService.DeleteUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponseDelete{
			Status:  false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.BaseResponseDelete{
		Status:  true,
		Message: "user deleted successfully",
	})
}
