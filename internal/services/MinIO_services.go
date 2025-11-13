package services

import (
	"context"
	"course_online_backend/database"
	"course_online_backend/internal/config"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type MinioHelper struct{}

//Biodata
func (m *MinioHelper) UploadProfilePicture(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	if database.MinioClient == nil {
		return "", fmt.Errorf("minio client not initialized")
	}

	ctx := context.Background()
	cfg := config.MinioConfig

	ext := filepath.Ext(fileHeader.Filename)
	filename := fmt.Sprintf("profiles/%s%s", uuid.New().String(), ext)

	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	_, err := database.MinioClient.PutObject(
		ctx,
		cfg.Bucket,
		filename,
		file,
		fileHeader.Size,
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to upload file to minio: %w", err)
	}

	protocol := "http"
	if cfg.UseSSL {
		protocol = "https"
	}
	fullURL := fmt.Sprintf("%s://%s/%s/%s", protocol, cfg.ImageENDPOINT, cfg.Bucket, filename)

	return fullURL, nil
}

func (m *MinioHelper) DeleteProfilePicture(profileURL string) error {
	if profileURL == "" || database.MinioClient == nil {
		return nil
	}

	objectName := m.extractObjectName(profileURL)
	if objectName == "" {
		return nil 
	}

	ctx := context.Background()
	cfg := config.MinioConfig

	err := database.MinioClient.RemoveObject(ctx, cfg.Bucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file from minio: %w", err)
	}

	return nil
}

func (m *MinioHelper) extractObjectName(url string) string {
	cfg := config.MinioConfig
	
	prefix := fmt.Sprintf("https://%s/%s/", cfg.Endpoint, cfg.Bucket)
	if strings.HasPrefix(url, prefix) {
		return strings.TrimPrefix(url, prefix)
	}
	
	prefix = fmt.Sprintf("http://%s/%s/", cfg.Endpoint, cfg.Bucket)
	if strings.HasPrefix(url, prefix) {
		return strings.TrimPrefix(url, prefix)
	}
	
	return ""
}


//course
func (m *MinioHelper) UploadCourseThumbnail(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	if database.MinioClient == nil {
		return "", fmt.Errorf("minio client not initialized")
	}

	ctx := context.Background()
	cfg := config.MinioConfig

	contentType := fileHeader.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return "", fmt.Errorf("file must be an image")
	}

	ext := filepath.Ext(fileHeader.Filename)
	filename := fmt.Sprintf("courses/thumbnails/%s%s", uuid.New().String(), ext)

	if contentType == "" {
		contentType = "application/octet-stream"
	}

	_, err := database.MinioClient.PutObject(
		ctx,
		cfg.Bucket,
		filename,
		file,
		fileHeader.Size,
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to upload thumbnail to minio: %w", err)
	}

	protocol := "http"
	if cfg.UseSSL {
		protocol = "https"
	}
	fullURL := fmt.Sprintf("%s://%s/%s/%s", protocol, cfg.ImageENDPOINT, cfg.Bucket, filename)

	return fullURL, nil
}

func (m *MinioHelper) DeleteCourseThumbnail(thumbnailURL string) error {
	if thumbnailURL == "" || database.MinioClient == nil {
		return nil
	}

	objectName := m.extractObjectName(thumbnailURL)
	if objectName == "" {
		return nil
	}

	ctx := context.Background()
	cfg := config.MinioConfig

	err := database.MinioClient.RemoveObject(ctx, cfg.Bucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete thumbnail from minio: %w", err)
	}

	return nil
}