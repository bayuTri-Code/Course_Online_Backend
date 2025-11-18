package CourseManagementhandler

import (
	"course_online_backend/internal/dto"
	"course_online_backend/internal/services"
	modulemanagementServices "course_online_backend/internal/services/Module_Management_Services"
	"course_online_backend/internal/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ModuleHandler struct {
	service         *modulemanagementServices.ModuleService
	ActivityService *services.ActivityService
}

func NewModuleHandler(service *modulemanagementServices.ModuleService, act *services.ActivityService) *ModuleHandler {
	return &ModuleHandler{
		service:         service,
		ActivityService: act,
	}
}

// CreateModuleHandler godoc
// @Summary Create a new module
// @Description Create a new module for a course
// @Tags modules
// @Accept json
// @Produce json
// @Param courseId path string true "Course ID (UUID)"
// @Param request body dto.CreateModuleRequest true "Module details"
// @Success 201 {object} utils.StandardResponse{data=dto.CreateModuleResponse}
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 404 {object} utils.ErrorResponse "Course not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/courses/{courseId}/modules [post]
// @Security BearerAuth
func (h *ModuleHandler) CreateModuleHandler(c *gin.Context) {
	courseIDParam := c.Param("id")
	courseID, err := uuid.Parse(courseIDParam)
	if err != nil {
		utils.JSONError(c, "Invalid course ID", http.StatusBadRequest, err.Error())
		return
	}

	var req dto.CreateModuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, "Invalid request data", http.StatusBadRequest, err.Error())
		return
	}

	req.CourseID = courseID

	module, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		if err.Error() == "course not found" {
			utils.JSONNotFound(c, "Course not found")
			return
		}
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONCreated(c, module, "Module created successfully")

	if firebaseUID, exists := c.Get("user_id"); exists {
		go func() {
			user, err := h.ActivityService.GetUserByFirebaseUID(firebaseUID.(string))
			if err != nil {
				fmt.Printf("User not found for activity logging")
				return
			}
			activityMsg := fmt.Sprintf("Created module: %s (ID: %s) From Course: %s", req.Name, module.ID.String(), courseID.String())
			_ = h.ActivityService.LogActivity(user.ID, activityMsg)
		}()
	}
}

// UpdateModuleHandler godoc
// @Summary Update a module
// @Description Update module information
// @Tags modules
// @Accept json
// @Produce json
// @Param id path string true "Module ID (UUID)"
// @Param request body dto.UpdateModuleRequest true "Module details"
// @Success 200 {object} utils.StandardResponse{data=dto.CreateModuleResponse}
// @Failure 400 {object} utils.ErrorResponse "Invalid request data"
// @Failure 404 {object} utils.ErrorResponse "Module not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/modules/{id} [put]
// @Security BearerAuth
func (h *ModuleHandler) UpdateModuleHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.JSONError(c, "Invalid module ID", http.StatusBadRequest, err.Error())
		return
	}

	var req dto.UpdateModuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, "Invalid request data", http.StatusBadRequest, err.Error())
		return
	}

	module, err := h.service.Update(c.Request.Context(), id, &req)
	if err != nil {
		if err.Error() == "module not found" {
			utils.JSONNotFound(c, "Module not found")
			return
		}
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONSuccess(c, module, "Module updated successfully")

	if firebaseUID, exists := c.Get("user_id"); exists {
		go func() {
			user, err := h.ActivityService.GetUserByFirebaseUID(firebaseUID.(string))
			if err != nil {
				fmt.Printf("User not found for activity logging")
				return
			}
			activityMsg := fmt.Sprintf("Updated module: %s (ID: %s)", module.Name, module.ID.String())
			_ = h.ActivityService.LogActivity(user.ID, activityMsg)
		}()
	}
}

// GetModuleByIDHandler godoc
// @Summary Get module by ID
// @Description Retrieve detailed information of a module by its ID
// @Tags modules
// @Accept json
// @Produce json
// @Param id path string true "Module ID (UUID)"
// @Success 200 {object} utils.StandardResponse{data=dto.ModuleDetailResponse}
// @Failure 400 {object} utils.ErrorResponse "Invalid module ID"
// @Failure 404 {object} utils.ErrorResponse "Module not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/modules/{id} [get]
func (h *ModuleHandler) GetModuleByIDHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.JSONError(c, "Invalid module ID", http.StatusBadRequest, err.Error())
		return
	}

	module, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "module not found" {
			utils.JSONNotFound(c, "Module not found")
			return
		}
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONSuccess(c, module, "Module retrieved successfully")
}

// GetModulesByCourseHandler godoc
// @Summary Get modules by course ID
// @Description Retrieve all modules for a specific course
// @Tags modules
// @Accept json
// @Produce json
// @Param courseId path string true "Course ID (UUID)"
// @Success 200 {object} utils.StandardResponse{data=[]dto.CreateModuleResponse}
// @Failure 400 {object} utils.ErrorResponse "Invalid course ID"
// @Failure 404 {object} utils.ErrorResponse "Course not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/courses/{courseId}/modules [get]
func (h *ModuleHandler) GetModulesByCourseHandler(c *gin.Context) {
	courseIDParam := c.Param("id")
	courseID, err := uuid.Parse(courseIDParam)
	if err != nil {
		utils.JSONError(c, "Invalid course ID", http.StatusBadRequest, err.Error())
		return
	}

	modules, err := h.service.GetByCourse(c.Request.Context(), courseID)
	if err != nil {
		if err.Error() == "course not found" {
			utils.JSONNotFound(c, "Course not found")
			return
		}
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONSuccess(c, modules, "Modules retrieved successfully")
}

// SoftDeleteModuleHandler godoc
// @Summary Soft delete a module
// @Description Soft delete a module (it can be restored later)
// @Tags modules
// @Accept json
// @Produce json
// @Param id path string true "Module ID (UUID)"
// @Success 200 {object} utils.StandardResponse "Module deleted successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid module ID"
// @Failure 404 {object} utils.ErrorResponse "Module not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/modules/{id} [delete]
// @Security BearerAuth
func (h *ModuleHandler) SoftDeleteModuleHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.JSONError(c, "Invalid module ID", http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.SoftDelete(c.Request.Context(), id); err != nil {
		if err.Error() == "module not found" {
			utils.JSONNotFound(c, "Module not found")
			return
		}
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONSuccess(c, nil, "Module deleted successfully")

	if firebaseUID, exists := c.Get("user_id"); exists {
		go func() {
			user, err := h.ActivityService.GetUserByFirebaseUID(firebaseUID.(string))
			if err != nil {
				fmt.Printf("User not found for activity logging")
				return
			}
			activityMsg := fmt.Sprintf("soft Deleted module: (ID: %s)", id.String())
			_ = h.ActivityService.LogActivity(user.ID, activityMsg)
		}()
	}
}

// RestoreModuleHandler godoc
// @Summary Restore a soft deleted module
// @Description Restore a previously soft deleted module
// @Tags modules
// @Accept json
// @Produce json
// @Param id path string true "Module ID (UUID)"
// @Success 200 {object} utils.StandardResponse{data=dto.CreateModuleResponse}
// @Failure 400 {object} utils.ErrorResponse "Invalid module ID"
// @Failure 404 {object} utils.ErrorResponse "Module not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/modules/{id}/restore [patch]
// @Security BearerAuth
func (h *ModuleHandler) RestoreModuleHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.JSONError(c, "Invalid module ID", http.StatusBadRequest, err.Error())
		return
	}

	module, err := h.service.Restore(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "module not found" {
			utils.JSONNotFound(c, "Module not found")
			return
		}
		if err.Error() == "module is not deleted" {
			utils.JSONError(c, "Module is not deleted", http.StatusBadRequest, nil)
			return
		}
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONSuccess(c, module, "Module restored successfully")

	if firebaseUID, exists := c.Get("user_id"); exists {
		go func() {
			user, err := h.ActivityService.GetUserByFirebaseUID(firebaseUID.(string))
			if err != nil {
				fmt.Printf("User not found for activity logging")
				return
			}
			activityMsg := fmt.Sprintf("Restore module: (ID: %s)", id.String())
			_ = h.ActivityService.LogActivity(user.ID, activityMsg)
		}()
	}
}

// PermanentDeleteModuleHandler godoc
// @Summary Permanently delete a module
// @Description Permanently delete a module (cannot be restored)
// @Tags modules
// @Accept json
// @Produce json
// @Param id path string true "Module ID (UUID)"
// @Success 200 {object} utils.StandardResponse "Module permanently deleted"
// @Failure 400 {object} utils.ErrorResponse "Invalid module ID"
// @Failure 404 {object} utils.ErrorResponse "Module not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/modules/{id}/permanent [delete]
// @Security BearerAuth
func (h *ModuleHandler) PermanentDeleteModuleHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.JSONError(c, "Invalid module ID", http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.PermanentDelete(c.Request.Context(), id); err != nil {
		if err.Error() == "module not found" {
			utils.JSONNotFound(c, "Module not found")
			return
		}
		utils.JSONError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utils.JSONSuccess(c, nil, "Module permanently deleted")

	if firebaseUID, exists := c.Get("user_id"); exists {
		go func() {
			user, err := h.ActivityService.GetUserByFirebaseUID(firebaseUID.(string))
			if err != nil {
				fmt.Printf("User not found for activity logging")
				return
			}
			activityMsg := fmt.Sprintf("Deleted Permanent module: (ID: %s)", id.String())
			_ = h.ActivityService.LogActivity(user.ID, activityMsg)
		}()
	}
}
