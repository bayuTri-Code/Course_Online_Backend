package LessonManagementHandler

import (
	"course_online_backend/internal/dto"
	"course_online_backend/internal/services"
	LessonManagementServices "course_online_backend/internal/services/Course_management_Services/Lesson_Management_Services"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LessonHandler struct {
	service         *LessonManagementServices.LessonService
	activityService *services.ActivityService
}

func NewLessonHandler(service *LessonManagementServices.LessonService, activityService *services.ActivityService) *LessonHandler {
	return &LessonHandler{
		service:         service,
		activityService: activityService,
	}
}

// CreateLesson godoc
// @Summary Create a new lesson
// @Description Create a new lesson in a module
// @Tags Lessons
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param module_id path string true "Module ID"
// @Param request body dto.CreateLessonRequest true "Create lesson request"
// @Success 201 {object} map[string]interface{} "success response with lesson data"
// @Failure 400 {object} map[string]interface{} "bad request"
// @Failure 404 {object} map[string]interface{} "module not found"
// @Failure 500 {object} map[string]interface{} "internal server error"
// @Router /modules/{module_id}/lessons [post]
func (h *LessonHandler) CreateLesson(c *gin.Context) {
	moduleIDStr := c.Param("module_id")
	moduleID, err := uuid.Parse(moduleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid module ID"})
		return
	}

	var req dto.CreateLessonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.ModuleID = moduleID

	result, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		if err.Error() == "module not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	

	c.JSON(http.StatusCreated, gin.H{
		"message": "Lesson created successfully",
		"data":    result,
	})

	if firebaseUID, exists := c.Get("user_id"); exists {
		go func() {
			user, err := h.activityService.GetUserByFirebaseUID(firebaseUID.(string))
			if err != nil {
				fmt.Printf("User not found for activity logging")
				return
			}
			activityMsg := fmt.Sprintf("Created lesson: %s (ID: %s) in module %s", req.Name, result.ID.String(), moduleID.String())
			_ = h.activityService.LogActivity(user.ID, activityMsg)
		}()
	}
}

// UpdateLesson godoc
// @Summary Update a lesson
// @Description Update lesson details
// @Tags Lessons
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Lesson ID"
// @Param request body dto.UpdateLessonRequest true "Update lesson request"
// @Success 200 {object} map[string]interface{} "success response with updated lesson"
// @Failure 400 {object} map[string]interface{} "bad request"
// @Failure 404 {object} map[string]interface{} "lesson not found"
// @Failure 500 {object} map[string]interface{} "internal server error"
// @Router /lessons/{id} [put]
func (h *LessonHandler) UpdateLesson(c *gin.Context) {
	idStr := c.Param("lesson_id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lesson ID"})
		return
	}

	var req dto.UpdateLessonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.Update(c.Request.Context(), id, &req)
	if err != nil {
		if err.Error() == "lesson not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	

	c.JSON(http.StatusOK, gin.H{
		"message": "Lesson updated successfully",
		"data":    result,
	})

	if firebaseUID, exists := c.Get("user_id"); exists {
    go func() {
        user, err := h.activityService.GetUserByFirebaseUID(firebaseUID.(string))
        if err != nil {
            fmt.Printf("User not found for activity logging")
            return
        }
        activityMsg := fmt.Sprintf("Updated lesson: %s (ID: %s)", result.Name, result.ID.String())
        _ = h.activityService.LogActivity(user.ID, activityMsg)
    }()
}
}

// GetLessonByID godoc
// @Summary Get lesson by ID
// @Description Get detailed information about a lesson
// @Tags Lessons
// @Produce json
// @Param id path string true "Lesson ID"
// @Success 200 {object} map[string]interface{} "success response with lesson data"
// @Failure 400 {object} map[string]interface{} "bad request"
// @Failure 404 {object} map[string]interface{} "lesson not found"
// @Failure 500 {object} map[string]interface{} "internal server error"
// @Router /lessons/{id} [get]
func (h *LessonHandler) GetLessonByID(c *gin.Context) {
	idStr := c.Param("lesson_id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lesson ID"})
		return
	}

	result, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "lesson not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Lesson retrieved successfully",
		"data":    result,
	})
}

// GetLessonsByModule godoc
// @Summary Get lessons by module ID
// @Description Get all lessons in a specific module
// @Tags Lessons
// @Produce json
// @Param module_id path string true "Module ID"
// @Success 200 {object} map[string]interface{} "success response with lessons list"
// @Failure 400 {object} map[string]interface{} "bad request"
// @Failure 404 {object} map[string]interface{} "module not found"
// @Failure 500 {object} map[string]interface{} "internal server error"
// @Router /modules/{module_id}/lessons [get]
func (h *LessonHandler) GetLessonsByModule(c *gin.Context) {
	moduleIDStr := c.Param("module_id")
	moduleID, err := uuid.Parse(moduleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid module ID"})
		return
	}

	result, err := h.service.GetByModule(c.Request.Context(), moduleID)
	if err != nil {
		if err.Error() == "module not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Lessons retrieved successfully",
		"data":    result,
	})
}


// DeleteLesson godoc
// @Summary Soft delete a lesson
// @Description Soft delete a lesson by ID
// @Tags Lessons
// @Produce json
// @Security BearerAuth
// @Param id path string true "Lesson ID"
// @Success 200 {object} map[string]interface{} "success message"
// @Failure 400 {object} map[string]interface{} "bad request"
// @Failure 404 {object} map[string]interface{} "lesson not found"
// @Failure 500 {object} map[string]interface{} "internal server error"
// @Router /lessons/{id} [delete]
func (h *LessonHandler) DeleteLesson(c *gin.Context) {
	idStr := c.Param("lesson_id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lesson ID"})
		return
	}

	if err := h.service.SoftDelete(c.Request.Context(), id); err != nil {
		if err.Error() == "lesson not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	

	c.JSON(http.StatusOK, gin.H{
		"message": "Lesson deleted successfully",
	})

	if firebaseUID, exists := c.Get("user_id"); exists {
    go func() {
        user, err := h.activityService.GetUserByFirebaseUID(firebaseUID.(string))
        if err != nil {
            fmt.Printf("User not found for activity logging")
            return
        }
        activityMsg := fmt.Sprintf("Soft deleted lesson (ID: %s)", id.String())
        _ = h.activityService.LogActivity(user.ID, activityMsg)
    }()
}
}

// RestoreLesson godoc
// @Summary Restore a soft deleted lesson
// @Description Restore a soft deleted lesson by ID
// @Tags Lessons
// @Produce json
// @Security BearerAuth
// @Param id path string true "Lesson ID"
// @Success 200 {object} map[string]interface{} "success response with restored lesson"
// @Failure 400 {object} map[string]interface{} "bad request"
// @Failure 404 {object} map[string]interface{} "lesson not found"
// @Failure 500 {object} map[string]interface{} "internal server error"
// @Router /lessons/{id}/restore [patch]
func (h *LessonHandler) RestoreLesson(c *gin.Context) {
	idStr := c.Param("lesson_id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lesson ID"})
		return
	}

	result, err := h.service.Restore(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, errors.New("lesson not found")) || errors.Is(err, errors.New("lesson is not deleted")) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	

	c.JSON(http.StatusOK, gin.H{
		"message": "Lesson restored successfully",
		"data":    result,
	})

	if firebaseUID, exists := c.Get("user_id"); exists {
    go func() {
        user, err := h.activityService.GetUserByFirebaseUID(firebaseUID.(string))
        if err != nil {
            fmt.Printf("User not found for activity logging")
            return
        }
        activityMsg := fmt.Sprintf("Restored lesson: %s (ID: %s)", result.Name, result.ID.String())
        _ = h.activityService.LogActivity(user.ID, activityMsg)
    }()
}
}

// PermanentDeleteLesson godoc
// @Summary Permanently delete a lesson
// @Description Hard delete a lesson by ID (cannot be restored)
// @Tags Lessons
// @Produce json
// @Security BearerAuth
// @Param id path string true "Lesson ID"
// @Success 200 {object} map[string]interface{} "success message"
// @Failure 400 {object} map[string]interface{} "bad request"
// @Failure 404 {object} map[string]interface{} "lesson not found"
// @Failure 500 {object} map[string]interface{} "internal server error"
// @Router /lessons/{id}/permanent [delete]
func (h *LessonHandler) PermanentDeleteLesson(c *gin.Context) {
	idStr := c.Param("lesson_id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lesson ID"})
		return
	}

	if err := h.service.PermanentDelete(c.Request.Context(), id); err != nil {
		if err.Error() == "lesson not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	

	c.JSON(http.StatusOK, gin.H{
		"message": "Lesson permanently deleted",
	})

	if firebaseUID, exists := c.Get("user_id"); exists {
    go func() {
        user, err := h.activityService.GetUserByFirebaseUID(firebaseUID.(string))
        if err != nil {
            fmt.Printf("User not found for activity logging")
            return
        }
        activityMsg := fmt.Sprintf("Permanently deleted lesson (ID: %s)", id.String())
        _ = h.activityService.LogActivity(user.ID, activityMsg)
    }()
}
}

// BulkCreateLessons godoc
// @Summary Bulk create lessons
// @Description Create multiple lessons at once
// @Tags Lessons
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param module_id path string true "Module ID"
// @Param request body dto.BulkCreateLessonRequest true "Bulk create lessons request"
// @Success 201 {object} map[string]interface{} "success response with created lessons"
// @Failure 400 {object} map[string]interface{} "bad request"
// @Failure 404 {object} map[string]interface{} "module not found"
// @Failure 500 {object} map[string]interface{} "internal server error"
// @Router /modules/{module_id}/lessons/bulk [post]
func (h *LessonHandler) BulkCreateLessons(c *gin.Context) {
	moduleIDStr := c.Param("module_id")
	moduleID, err := uuid.Parse(moduleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid module ID"})
		return
	}

	var req dto.BulkCreateLessonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.BulkCreate(c.Request.Context(), moduleID, &req)
	if err != nil {
		if err.Error() == "module not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	

	c.JSON(http.StatusCreated, gin.H{
		"message": "Lessons created successfully",
		"data":    result,
	})
	if firebaseUID, exists := c.Get("user_id"); exists {
    go func() {
        user, err := h.activityService.GetUserByFirebaseUID(firebaseUID.(string))
        if err != nil {
            fmt.Printf("User not found for activity logging")
            return
        }
        activityMsg := fmt.Sprintf("Bulk created %d lessons in module %s", len(result), moduleID.String())
        _ = h.activityService.LogActivity(user.ID, activityMsg)
    }()
}
}

// ReorderLessons godoc
// @Summary Reorder lessons
// @Description Update course order for multiple lessons
// @Tags Lessons
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param module_id path string true "Module ID"
// @Param request body dto.ReorderLessonsRequest true "Reorder lessons request"
// @Success 200 {object} map[string]interface{} "success message"
// @Failure 400 {object} map[string]interface{} "bad request"
// @Failure 404 {object} map[string]interface{} "module not found"
// @Failure 500 {object} map[string]interface{} "internal server error"
// @Router /modules/{module_id}/lessons/reorder [patch]
func (h *LessonHandler) ReorderLessons(c *gin.Context) {
	moduleIDStr := c.Param("module_id")
	moduleID, err := uuid.Parse(moduleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid module ID"})
		return
	}

	var req dto.ReorderLessonsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Reorder(c.Request.Context(), moduleID, &req); err != nil {
		if err.Error() == "module not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	

	c.JSON(http.StatusOK, gin.H{
		"message": "Lessons reordered successfully",
	})

	if firebaseUID, exists := c.Get("user_id"); exists {
    go func() {
        user, err := h.activityService.GetUserByFirebaseUID(firebaseUID.(string))
        if err != nil {
            fmt.Printf("User not found for activity logging")
            return
        }
        activityMsg := fmt.Sprintf("Reordered lessons in module %s", moduleID.String())
        _ = h.activityService.LogActivity(user.ID, activityMsg)
    }()
}
}