package services

import (
	"course_online_backend/internal/dto"
	"course_online_backend/internal/models"
	"errors"
	"fmt"
	"mime/multipart"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BiodataService struct {
	DB          *gorm.DB
	MinioHelper *MinioHelper
}

func NewBiodataService(db *gorm.DB) *BiodataService {
	return &BiodataService{
		DB:          db,
		MinioHelper: &MinioHelper{},
	}
}

func (s *BiodataService) CreateBiodata(firebaseUID string, req dto.CreateBiodataRequest, file multipart.File, fileHeader *multipart.FileHeader) (*models.Biodata, error) {
	var user models.User
	if err := s.DB.Where("firebase_uid = ?", firebaseUID).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	profileURL := ""
	if file != nil && fileHeader != nil {
		uploadedURL, err := s.MinioHelper.UploadProfilePicture(file, fileHeader)
		if err != nil {
			return nil, errors.New("failed to upload profile picture: " + err.Error())
		}
		profileURL = uploadedURL
	}

	biodata := models.Biodata{
		ID:             uuid.New(),
		UserID:         user.ID,
		Name:           req.Name,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Description:    req.Description,
		Contact:        req.Contact,
		Age:            req.Age,
		School:         req.School,
		ProfilePicture: profileURL,

	}

	if err := s.DB.Create(&biodata).Error; err != nil {
		if profileURL != "" {
			s.MinioHelper.DeleteProfilePicture(profileURL)
		}
		return nil, err
	}

	return &biodata, nil
}

func (s *BiodataService) GetBiodata(firebaseUID string) (*models.Biodata, error) {
	var user models.User
	if err := s.DB.Where("firebase_uid = ?", firebaseUID).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	var biodata models.Biodata
	if err := s.DB.Where("user_id = ?", user.ID).First(&biodata).Error; err != nil {
		return nil, err
	}

	return &biodata, nil
}

func (s *BiodataService) UpdateBiodata(firebaseUID string, req dto.UpdateBiodataRequest, file multipart.File, fileHeader *multipart.FileHeader) (*models.Biodata, error) {
	var user models.User
	if err := s.DB.Where("firebase_uid = ?", firebaseUID).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	var biodata models.Biodata
	if err := s.DB.Where("user_id = ?", user.ID).First(&biodata).Error; err != nil {
		return nil, err
	}

	if req.Name != "" {
		biodata.Name = req.Name
	}
	if req.FirstName != "" {
		biodata.FirstName = req.FirstName
	}
	if req.LastName != "" {
		biodata.LastName = req.LastName
	}
	if req.Description != "" {
		biodata.Description = req.Description
	}
	if req.Contact != "" {
		biodata.Contact = req.Contact
	}
	if req.Age != 0 {
		biodata.Age = req.Age
	}
	if req.School != "" {
		biodata.School = req.School
	}

	if file != nil && fileHeader != nil {
		uploadedURL, err := s.MinioHelper.UploadProfilePicture(file, fileHeader)
		if err != nil {
			return nil, errors.New("failed to upload profile picture: " + err.Error())
		}

		oldProfileURL := biodata.ProfilePicture
		if oldProfileURL != "" {
			if err := s.MinioHelper.DeleteProfilePicture(oldProfileURL); err != nil {
				fmt.Printf("erorr to delete old file")
				fmt.Printf("failed to delete old profile picture")
			}
		}

		biodata.ProfilePicture = uploadedURL
	}

	if err := s.DB.Save(&biodata).Error; err != nil {
		return nil, err
	}

	return &biodata, nil
}

func (s *BiodataService) DeleteBiodata(firebaseUID string) error {
	var user models.User
	if err := s.DB.Where("firebase_uid = ?", firebaseUID).First(&user).Error; err != nil {
		return errors.New("user not found")
	}

	var biodata models.Biodata
	if err := s.DB.Where("user_id = ?", user.ID).First(&biodata).Error; err != nil {
		return err
	}

	if biodata.ProfilePicture != "" {
		if err := s.MinioHelper.DeleteProfilePicture(biodata.ProfilePicture); err != nil {
			fmt.Printf("failed to delete old profile picture")
		}
	}

	if err := s.DB.Where("user_id = ?", user.ID).Delete(&models.Biodata{}).Error; err != nil {
		return err
	}

	return nil
}


func (s *BiodataService) SoftDeleteBiodata(firebaseUID string) error {
	var user models.User
	if err := s.DB.Where("firebase_uid = ?", firebaseUID).First(&user).Error; err != nil {
		return errors.New("user not found")
	}

	// update deleted_at instead of permanent delete
	if err := s.DB.Where("user_id = ?", user.ID).Delete(&models.Biodata{}).Error; err != nil {
		return err
	}

	return nil
}

func (s *BiodataService) RestoreBiodata(firebaseUID string) (*models.Biodata, error) {
	var user models.User
	if err := s.DB.Unscoped().Where("firebase_uid = ?", firebaseUID).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	var biodata models.Biodata
	if err := s.DB.Unscoped().Where("user_id = ?", user.ID).First(&biodata).Error; err != nil {
		return nil, errors.New("biodata not found or already active")
	}

	if biodata.DeletedAt.Valid {
		biodata.DeletedAt = gorm.DeletedAt{} 
		if err := s.DB.Unscoped().Save(&biodata).Error; err != nil {
			return nil, err
		}
	}

	return &biodata, nil
}


