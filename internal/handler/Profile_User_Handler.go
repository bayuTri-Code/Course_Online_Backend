package handler

import (
	"course_online_backend/internal/dto"
	"course_online_backend/internal/services"
	"net/http"
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

type BiodataController struct {
	BiodataService *services.BiodataService
}

func NewBiodataHandler(svc *services.BiodataService) *BiodataController {
	return &BiodataController{BiodataService: svc}
}

// @Summary Create Biodata
// @Tags Profile
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true
// @Param age formData int true
// @Param school formData string true
// @Param profile_picture formData file false
// @Success 200 {object} dto.BiodataResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /profile/biodata [post]
func (c *BiodataController) CreateBiodata(ctx *gin.Context) {
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

	ctx.JSON(http.StatusOK, dto.BiodataResponse{
		ID:             biodata.ID.String(),
		UserID:         biodata.UserID.String(),
		Name:           biodata.Name,
		Age:            biodata.Age,
		School:         biodata.School,
		ProfilePicture: biodata.ProfilePicture,
	})
}

// @Summary Get Biodata
// @Tags Profile
// @Produce json
// @Success 200 {object} dto.BiodataResponse
// @Failure 401 {object} map[string]string
// @Router /profile/biodata [get]
func (c *BiodataController) GetBiodata(ctx *gin.Context) {
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

	ctx.JSON(http.StatusOK, dto.BiodataResponse{
		ID:             biodata.ID.String(),
		UserID:         biodata.UserID.String(),
		Name:           biodata.Name,
		Age:            biodata.Age,
		School:         biodata.School,
		ProfilePicture: biodata.ProfilePicture,
	})
}

// @Summary Update Biodata
// @Tags Profile
// @Accept multipart/form-data
// @Produce json
// @Param name formData string false
// @Param age formData int false
// @Param school formData string false
// @Param profile_picture formData file false
// @Success 200 {object} dto.BiodataResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /profile/biodata [put]
func (c *BiodataController) UpdateBiodata(ctx *gin.Context) {
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

	ctx.JSON(http.StatusOK, dto.BiodataResponse{
		ID:             biodata.ID.String(),
		UserID:         biodata.UserID.String(),
		Name:           biodata.Name,
		Age:            biodata.Age,
		School:         biodata.School,
		ProfilePicture: biodata.ProfilePicture,
	})
}

// @Summary Delete Biodata
// @Tags Profile
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /profile/biodata [delete]
func (c *BiodataController) DeleteBiodata(ctx *gin.Context) {
	firebaseUID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := c.BiodataService.DeleteBiodata(firebaseUID.(string)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Biodata deleted successfully"})
}
