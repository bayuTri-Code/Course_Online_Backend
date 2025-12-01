package zoom_management_services

import (
	"context"
	"course_online_backend/internal/dto"
	"course_online_backend/internal/models"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ZoomService struct {
	db *gorm.DB
}

func NewZoomService(db *gorm.DB) *ZoomService {
	return &ZoomService{db: db}
}

func (s *ZoomService) mapToZoomResponse(zoom *models.Zoom) dto.ZoomResponse {
	resp := dto.ZoomResponse{
		ID:          zoom.ID,
		Title:       zoom.Title,
		Link:        zoom.Link,
		Description: zoom.Description,
		ScheduledAt: zoom.ScheduledAt,
		Duration:    zoom.Duration,
		CourseID:    zoom.CourseID,
		CreatedAt:   zoom.CreatedAt,
		UpdatedAt:   zoom.UpdatedAt,
	}

	if zoom.Course.ID != uuid.Nil {
		resp.Course = &dto.ZoomCourseResponse{
			ID:          zoom.Course.ID,
			Name:        zoom.Course.Name,
			Description: zoom.Course.Description,
			Thumbnail:   zoom.Course.Thumbnail,
		}
	}

	return resp
}

func (s *ZoomService) CreateZoom(ctx context.Context, req *dto.CreateZoomRequest) (*dto.ZoomResponse, error) {
	var course models.Course
	if err := s.db.First(&course, req.CourseID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("course not found")
		}
		return nil, err
	}

	zoom := &models.Zoom{
		Title:       req.Title,
		Link:        req.Link,
		Description: req.Description,
		ScheduledAt: req.ScheduledAt,
		Duration:    req.Duration,
		CourseID:    req.CourseID,
	}

	if err := s.db.Create(zoom).Error; err != nil {
		return nil, err
	}

	return s.GetZoomByID(ctx, zoom.ID)
}

func (s *ZoomService) GetZoomByID(ctx context.Context, id uuid.UUID) (*dto.ZoomResponse, error) {
	var zoom models.Zoom
	if err := s.db.Preload("Course").First(&zoom, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("zoom not found")
		}
		return nil, err
	}

	response := s.mapToZoomResponse(&zoom)
	return &response, nil
}

func (s *ZoomService) GetZoomsByCourseID(ctx context.Context, courseID uuid.UUID) (*dto.ZoomListResponse, error) {
	var course models.Course
	if err := s.db.First(&course, courseID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("course not found")
		}
		return nil, err
	}

	var zooms []models.Zoom
	if err := s.db.Preload("Course").
		Where("course_id = ?", courseID).
		Order("scheduled_at ASC, created_at ASC").
		Find(&zooms).Error; err != nil {
		return nil, err
	}

	responses := make([]dto.ZoomResponse, len(zooms))
	for i, zoom := range zooms {
		responses[i] = s.mapToZoomResponse(&zoom)
	}

	return &dto.ZoomListResponse{
		Total: int64(len(zooms)),
		Zooms: responses,
	}, nil
}

func (s *ZoomService) GetAllZooms(ctx context.Context) ([]dto.ZoomResponse, error) {
	var zooms []models.Zoom
	if err := s.db.Preload("Course").
		Order("scheduled_at ASC, created_at ASC").
		Find(&zooms).Error; err != nil {
		return nil, err
	}

	responses := make([]dto.ZoomResponse, len(zooms))
	for i, zoom := range zooms {
		responses[i] = s.mapToZoomResponse(&zoom)
	}

	return responses, nil
}

func (s *ZoomService) UpdateZoom(ctx context.Context, id uuid.UUID, req *dto.UpdateZoomRequest) (*dto.ZoomResponse, error) {
	var zoom models.Zoom
	if err := s.db.First(&zoom, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("zoom not found")
		}
		return nil, err
	}

	zoom.Title = req.Title
	zoom.Link = req.Link
	zoom.Description = req.Description
	zoom.ScheduledAt = req.ScheduledAt
	zoom.Duration = req.Duration

	if err := s.db.Save(&zoom).Error; err != nil {
		return nil, err
	}

	return s.GetZoomByID(ctx, zoom.ID)
}

func (s *ZoomService) DeleteZoom(ctx context.Context, id uuid.UUID) error {
	var zoom models.Zoom
	if err := s.db.First(&zoom, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("zoom not found")
		}
		return err
	}

	return s.db.Delete(&zoom).Error
}

func (s *ZoomService) GetUpcomingZoomsByCourseID(ctx context.Context, courseID uuid.UUID) (*dto.ZoomListResponse, error) {
	var zooms []models.Zoom
	now := time.Now()
	
	if err := s.db.Preload("Course").
		Where("course_id = ? AND scheduled_at > ?", courseID, now).
		Order("scheduled_at ASC").
		Find(&zooms).Error; err != nil {
		return nil, err
	}

	responses := make([]dto.ZoomResponse, len(zooms))
	for i, zoom := range zooms {
		responses[i] = s.mapToZoomResponse(&zoom)
	}

	return &dto.ZoomListResponse{
		Total: int64(len(zooms)),
		Zooms: responses,
	}, nil
}
