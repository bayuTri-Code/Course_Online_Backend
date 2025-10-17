package handler

import (
	"course_online_backend/internal/dto"
	"course_online_backend/internal/services"
	"course_online_backend/internal/utils"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BiodataHandler struct {
	BiodataService *services.BiodataService
}

func NewBiodataHandler(svc *services.BiodataService) *BiodataHandler {
	return &BiodataHandler{BiodataService: svc}
}

// @Summary Create Biodata
// @Tags Profile
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "User name"
// @Param age formData int true "User age"
// @Param school formData string true "User school"
// @Param profile_picture formData file false "Profile picture upload"
// @Success 200 {object} dto.BaseResponseBiodata "Biodata successfully created"
// @Failure 400 {object} utils.ErrorResponse "Invalid request body"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Router /api/profile/biodata [post]
func (c *BiodataHandler) CreateBiodata(ctx *gin.Context) {
	firebaseUID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req dto.CreateBiodataRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var file multipart.File
	var fileHeader *multipart.FileHeader
	if f, fh, err := ctx.Request.FormFile("profile_picture"); err == nil {
		file = f
		fileHeader = fh
	}

	biodata, err := c.BiodataService.CreateBiodata(firebaseUID.(string), req, file, fileHeader)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	createBiodata := dto.BiodataResponse{
		ID:             biodata.ID.String(),
		UserID:         biodata.UserID.String(),
		Name:           biodata.Name,
		Age:            biodata.Age,
		School:         biodata.School,
		ProfilePicture: biodata.ProfilePicture,
	}

	utils.JSONCreated(ctx, createBiodata, "succses")
}

// @Summary Get Biodata
// @Tags Profile
// @Accept json
// @Produce json
// @Success 200 {object} dto.BaseResponseBiodata "Get Biodata successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 404 {object} utils.ErrorResponse "Biodata not found"
// @Router /api/profile/biodata [get]
func (c *BiodataHandler) GetBiodata(ctx *gin.Context) {
	firebaseUID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	biodata, err := c.BiodataService.GetBiodata(firebaseUID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	GetBiodata := dto.BiodataResponse{
		ID:             biodata.ID.String(),
		UserID:         biodata.UserID.String(),
		Name:           biodata.Name,
		Age:            biodata.Age,
		School:         biodata.School,
		ProfilePicture: biodata.ProfilePicture,
	}

	utils.JSONSuccess(ctx, GetBiodata, "Succses get data")
}

// @Summary Update Biodata
// @Tags Profile
// @Accept multipart/form-data
// @Produce json
// @Param name formData string false "User name (optional)"
// @Param age formData int false "User age (optional)"
// @Param school formData string false "User school (optional)"
// @Param profile_picture formData file false "Profile picture upload (optional)"
// @Success 200 {object} dto.BaseResponseBiodata "Biodata successfully updated"
// @Failure 400 {object} utils.ErrorResponse "Invalid request body"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 404 {object} utils.ErrorResponse "Biodata not found"
// @Router /api/profile/biodata [put]
func (c *BiodataHandler) UpdateBiodata(ctx *gin.Context) {
	firebaseUID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req dto.UpdateBiodataRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var file multipart.File
	var fileHeader *multipart.FileHeader
	if f, fh, err := ctx.Request.FormFile("profile_picture"); err == nil {
		file = f
		fileHeader = fh
	}

	biodata, err := c.BiodataService.UpdateBiodata(firebaseUID.(string), req, file, fileHeader)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updateBiodata := dto.BiodataResponse{
		ID:             biodata.ID.String(),
		UserID:         biodata.UserID.String(),
		Name:           biodata.Name,
		Age:            biodata.Age,
		School:         biodata.School,
		ProfilePicture: biodata.ProfilePicture,
	}

	utils.JSONSuccess(ctx, updateBiodata, "Succes Update biodata")
}

// @Summary Delete Biodata
// @Tags Profile
// @Accept json
// @Produce json
// @Success 200 {object} utils.BaseResponse "Biodata successfully deleted"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 404 {object} utils.ErrorResponse "Biodata not found"
// @Router /api/profile/biodata [delete]
func (c *BiodataHandler) DeleteBiodata(ctx *gin.Context) {
	firebaseUID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := c.BiodataService.DeleteBiodata(firebaseUID.(string)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, utils.BaseResponse{
		Status:  true,
		Message: "Biodata successfully deleted",
	})

}
