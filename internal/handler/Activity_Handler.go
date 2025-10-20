package handler

import (
	"course_online_backend/internal/services"
	"course_online_backend/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ActivityHandler struct {
	ActivityService *services.ActivityService
}

func NewActivityHandler(svc *services.ActivityService) *ActivityHandler {
	return &ActivityHandler{ActivityService: svc}
}

// GetAllActivity godoc
// @Summary Get all user activities
// @Description Retrieve all activities belonging to the authenticated user. The user is identified using Firebase UID from the JWT token.
// @Tags Activity
// @Security BearerAuth
// @Produce json
// @Success 200 {object} utils.BaseResponse{data=[]dto.ActivityResponse} "Success get All activity"
// @Failure 401 {object} utils.BaseResponse "Unauthorized"
// @Failure 500 {object} utils.BaseResponse "Internal Server Error"
// @Router /api/history/activity [get]
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

	activities, err := h.ActivityService.GeAllActivity(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.JSONSuccess(ctx, activities, "Success get activity")
}

// GetActivityByUser godoc
// @Summary      Get Activity By User ID
// @Description  Mengambil semua aktivitas berdasarkan ID user yang diberikan
// @Tags         Activity
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID (UUID)"
// @Success      200  {object}  dto.BaseResponseActivityByUserId  "Success get activity by user"
// @Failure      400  {object}  utils.ErrorResponse    "Invalid user ID"
// @Failure      500  {object}  utils.ErrorResponse    "Internal server error"
// @Router       /activity/user/{id} [get]
func (h *ActivityHandler) GetActivityByUser(ctx *gin.Context) {
	userIDStr := ctx.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	data, err := h.ActivityService.GetActivityByUserID(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.JSONSuccess(ctx, data, "Success get activity by user")
}
