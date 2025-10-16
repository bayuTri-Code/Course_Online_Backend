package handler

import (
	"course_online_backend/internal/dto"
	"course_online_backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BiodataController struct {
	BiodataService *services.BiodataService
}

func NewBiodataHandler(biodataService *services.BiodataService) *BiodataController {
	return &BiodataController{BiodataService: biodataService}
}

func (c *BiodataController) CreateBiodata(ctx *gin.Context) {
	userID := ctx.GetString("uid") 

	var req dto.CreateBiodataRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, fileHeader, _ := ctx.Request.FormFile("profile_picture")

	biodata, err := c.BiodataService.CreateBiodata(userID, req, file, fileHeader)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.BiodataResponse{
		ID:             biodata.ID.String(),
		UserID:         biodata.UserID,
		Name:           biodata.Name,
		Age:            biodata.Age,
		School:         biodata.School,
		ProfilePicture: biodata.ProfilePicture,
	})
}
