package handler

import (
	"course_online_backend/internal/services"
	"course_online_backend/internal/utils"
	"net/http"
	"strconv"

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
// @Summary      Get all user activities
// @Description  Retrieve all activities belonging to the authenticated user based on Firebase UID from JWT token
// @Tags         Activity
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  utils.BaseResponse{data=[]dto.ActivityResponse}  "Success get all activity"
// @Failure      401  {object}  utils.ErrorResponse  "Unauthorized"
// @Failure      500  {object}  utils.ErrorResponse  "Internal server error"
// @Router       /api/history/All-activity [get]
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

	activities, err := h.ActivityService.GetAllActivity(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.JSONSuccess(ctx, activities, "Success get activity")
}

// GetActivityByUserId godoc
// @Summary      Get Activity By User ID
// @Description  Get all activities associated with a specific user ID
// @Tags         Activity
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID (UUID)"
// @Success      200  {object}  utils.BaseResponse{data=[]dto.ActivityResponse}  "Success get activity by user"
// @Failure      400  {object}  utils.ErrorResponse  "Invalid user ID"
// @Failure      500  {object}  utils.ErrorResponse  "Internal server error"
// @Router       /api/history/ByUser-activity/{id} [get]
func (h *ActivityHandler) GetActivityByUserId(ctx *gin.Context) {
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

// GetRecentActivities godoc
// @Summary      Get Recent Activities
// @Description  Retrieve a list of the most recent user activities
// @Tags         Activity
// @Produce      json
// @Param        limit   query     int  false  "Number of recent activities to return (default 10)"
// @Success      200  {object}  utils.BaseResponse{data=[]dto.ActivityResponse}  "Success get recent activities"
// @Failure      500  {object}  utils.ErrorResponse  "Internal server error"
// @Router       /api/history/Recent-activity [get]
func (h *ActivityHandler) GetRecentActivities(ctx *gin.Context) {
	limitStr := ctx.DefaultQuery("limit", "10")
	limit, _ := strconv.Atoi(limitStr)

	data, err := h.ActivityService.GetRecentActivities(limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.JSONSuccess(ctx, data, "Success get recent activities")
}

// SearchActivity godoc
// @Summary      Search Activity
// @Description  Search activity records by keyword
// @Tags         Activity
// @Produce      json
// @Param        keyword   query     string  true  "Keyword to search activity"
// @Success      200  {object}  utils.BaseResponse{data=[]dto.ActivityResponse}  "Success search activity"
// @Failure      400  {object}  utils.ErrorResponse  "Keyword is required"
// @Failure      500  {object}  utils.ErrorResponse  "Internal server error"
// @Router       /api/history/Search-activity [get]
func (h *ActivityHandler) SearchActivity(ctx *gin.Context) {
	keyword := ctx.Query("keyword")
	if keyword == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "keyword is required"})
		return
	}

	data, err := h.ActivityService.SearchActivity(keyword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.JSONSuccess(ctx, data, "Success search activity")
}

// GetActivitySummary godoc
// @Summary      Get Activity Summary
// @Description  Get summarized activity statistics for a specific user
// @Tags         Activity
// @Produce      json
// @Param        id   path      string  true  "User ID (UUID)"
// @Success      200  {object}  utils.BaseResponse{data=[]dto.ActivityResponse}  "Success get activity summary"
// @Failure      400  {object}  utils.ErrorResponse  "Invalid user ID"
// @Failure      500  {object}  utils.ErrorResponse  "Internal server error"
// @Router       /api/history/Summary-activity/{id} [get]
func (h *ActivityHandler) GetActivitySummary(ctx *gin.Context) {
	userID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	data, err := h.ActivityService.GetActivitySummary(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.JSONSuccess(ctx, data, "Success get activity summary")
}


// GetActivityByRole godoc
// @Summary      Get Activity By Role
// @Description  Retrieve all activities performed by users with a specific role
// @Tags         Activity
// @Produce      json
// @Param        role   path      string  true  "User Role (e.g. student, instructor, admin, student)"
// @Success      200  {object}  utils.BaseResponse{data=[]dto.ActivityResponse}  "Success get activity by role"
// @Failure      500  {object}  utils.ErrorResponse  "Internal server error"
// @Router       /api/history/ByRole-activity/{role} [get]
func (h *ActivityHandler) GetActivityByRole(ctx *gin.Context) {
	role := ctx.Param("role")
	data, err := h.ActivityService.GetActivityByRole(role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.JSONSuccess(ctx, data, "Success get activity by role")
}
