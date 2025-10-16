package services

import (
	"context"
	"course_online_backend/internal/dto"
	"course_online_backend/internal/models"
	"fmt"
	"mime/multipart"

	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

type BiodataService struct {
	DB          *gorm.DB
	MinioClient *minio.Client
	BucketName  string
	MinioURL    string
}

func NewBiodataService(db *gorm.DB, minioClient *minio.Client, bucketName, minioURL string) *BiodataService {
	return &BiodataService{
		DB:          db,
		MinioClient: minioClient,
		BucketName:  bucketName,
		MinioURL:    minioURL,
	}
}

func (s *BiodataService) CreateBiodata(userID string, req dto.CreateBiodataRequest, file multipart.File, fileHeader *multipart.FileHeader) (*models.Biodata, error) {
	var profilePictureURL string

	if file != nil {
		objectName := fmt.Sprintf("profiles/%s_%s", userID, fileHeader.Filename)

		_, err := s.MinioClient.PutObject(
			context.Background(),
			s.BucketName,
			objectName,
			file,
			fileHeader.Size,
			minio.PutObjectOptions{ContentType: fileHeader.Header.Get("Content-Type")},
		)
		if err != nil {
			return nil, fmt.Errorf("gagal upload ke MinIO: %v", err)
		}

		profilePictureURL = fmt.Sprintf("%s/%s/%s", s.MinioURL, s.BucketName, objectName)
	}

	biodata := models.Biodata{
		UserID:         userID,
		Name:           req.Name,
		Age:            req.Age,
		School:         req.School,
		ProfilePicture: profilePictureURL,
	}

	if err := s.DB.Create(&biodata).Error; err != nil {
		return nil, err
	}

	return &biodata, nil
}
