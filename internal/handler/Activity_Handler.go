package handler

import (
	"course_online_backend/internal/services"
	"course_online_backend/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ActivityHandler struct {
	ActivityService *services.ActivityService
}

func NewActivityHandler(svc *services.ActivityService) *ActivityHandler {
	return &ActivityHandler{ActivityService: svc}
}

func (h *ActivityHandler) GetAllActivity(ctx *gin.Context) {
	firebaseUID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := h.ActivityService.GetUserByFirebaseUID(firebaseUID.(string))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	activities, err := h.ActivityService.GetActivity(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.JSONSuccess(ctx, activities, "Success get activity")
}
