package LessonManagementServices

import (
	"context"
	"course_online_backend/internal/dto"
	"course_online_backend/internal/models"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LessonService struct {
	db *gorm.DB
}

func NewLessonService(db *gorm.DB) *LessonService {
	return &LessonService{db: db}
}


func (s *LessonService) toLessonResponse(lesson *models.Lesson) dto.LessonResponse {
	var deletedAt *time.Time
	if lesson.DeletedAt.Valid {
		deletedAt = &lesson.DeletedAt.Time
	}

	return dto.LessonResponse{
		ID:           lesson.ID,
		ModuleID:     lesson.ModuleID,
		VideoID:      lesson.VideoID,
		Name:         lesson.Name,
		VideoDetails: lesson.VideoDetails,
		CourseOrder:  lesson.CourseOrder,
		CreatedAt:    lesson.CreatedAt,
		UpdatedAt:    lesson.UpdatedAt,
		DeletedAt:    deletedAt,
	}
}

func (s *LessonService) toLessonDetailResponse(lesson *models.Lesson) dto.LessonDetailResponse {
	var deletedAt *time.Time
	if lesson.DeletedAt.Valid {
		deletedAt = &lesson.DeletedAt.Time
	}

	response := dto.LessonDetailResponse{
		ID:           lesson.ID,
		ModuleID:     lesson.ModuleID,
		VideoID:      lesson.VideoID,
		Name:         lesson.Name,
		VideoDetails: lesson.VideoDetails,
		CourseOrder:  lesson.CourseOrder,
		CreatedAt:    lesson.CreatedAt,
		UpdatedAt:    lesson.UpdatedAt,
		DeletedAt:    deletedAt,
	}

	if lesson.Module.ID != uuid.Nil {
		response.Module = &dto.ModuleBasicInfo{
			ID:       lesson.Module.ID,
			CourseID: lesson.Module.CourseID,
			Name:     lesson.Module.Name,
			Number:   lesson.Module.Number,
		}
	}

	return response
}

func (s *LessonService) toLessonResponses(lessons []models.Lesson) []dto.LessonResponse {
	responses := make([]dto.LessonResponse, len(lessons))
	for i, lesson := range lessons {
		responses[i] = s.toLessonResponse(&lesson)
	}
	return responses
}

func (s *LessonService) Create(ctx context.Context, req *dto.CreateLessonRequest) (*dto.LessonResponse, error) {
	var module models.Module
	if err := s.db.First(&module, "id = ?", req.ModuleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("module not found")
		}
		return nil, err
	}

	lesson := &models.Lesson{
		ModuleID:     req.ModuleID,
		VideoID:      req.VideoID,
		Name:         req.Name,
		VideoDetails: req.VideoDetails,
		CourseOrder:  req.CourseOrder,
	}

	if err := s.db.Create(lesson).Error; err != nil {
		return nil, err
	}

	response := s.toLessonResponse(lesson)
	return &response, nil
}

func (s *LessonService) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateLessonRequest) (*dto.LessonResponse, error) {
	var lesson models.Lesson
	if err := s.db.First(&lesson, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("lesson not found")
		}
		return nil, err
	}

	if req.Name != nil {
		lesson.Name = *req.Name
	}
	if req.VideoID != nil {
		lesson.VideoID = *req.VideoID
	}
	if req.VideoDetails != nil {
		lesson.VideoDetails = *req.VideoDetails
	}
	if req.CourseOrder != nil {
		lesson.CourseOrder = *req.CourseOrder
	}

	if err := s.db.Save(&lesson).Error; err != nil {
		return nil, err
	}

	response := s.toLessonResponse(&lesson)
	return &response, nil
}

func (s *LessonService) GetByID(ctx context.Context, id uuid.UUID) (*dto.LessonDetailResponse, error) {
	var lesson models.Lesson
	if err := s.db.Preload("Module").First(&lesson, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("lesson not found")
		}
		return nil, err
	}

	response := s.toLessonDetailResponse(&lesson)
	return &response, nil
}

func (s *LessonService) GetByModule(ctx context.Context, moduleID uuid.UUID) ([]dto.LessonResponse, error) {
	var module models.Module
	if err := s.db.First(&module, "id = ?", moduleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("module not found")
		}
		return nil, err
	}

	var lessons []models.Lesson
	if err := s.db.Where("module_id = ?", moduleID).Order("course_order ASC").Find(&lessons).Error; err != nil {
		return nil, err
	}

	return s.toLessonResponses(lessons), nil
}

// func (s *LessonService) GetAll(ctx context.Context, moduleID *uuid.UUID, page, pageSize int, search string) (*dto.LessonListResponse, error) {
// 	if page < 1 {
// 		page = 1
// 	}
// 	if pageSize < 1 {
// 		pageSize = 10
// 	}

// 	query := s.db.Model(&models.Lesson{})

// 	if moduleID != nil {
// 		query = query.Where("module_id = ?", *moduleID)
// 	}

// 	if search != "" {
// 		searchPattern := fmt.Sprintf("%%%s%%", search)
// 		query = query.Where("name ILIKE ? OR video_details ILIKE ?", searchPattern, searchPattern)
// 	}

// 	var total int64
// 	if err := query.Count(&total).Error; err != nil {
// 		return nil, err
// 	}

// 	var lessons []models.Lesson
// 	offset := (page - 1) * pageSize
// 	if err := query.Order("course_order ASC").Offset(offset).Limit(pageSize).Find(&lessons).Error; err != nil {
// 		return nil, err
// 	}

// 	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

// 	return &dto.LessonListResponse{
// 		Lessons:    s.toLessonResponses(lessons),
// 		Total:      total,
// 		Page:       page,
// 		PageSize:   pageSize,
// 		TotalPages: totalPages,
// 	}, nil
// }

func (s *LessonService) SoftDelete(ctx context.Context, id uuid.UUID) error {
	var lesson models.Lesson
	if err := s.db.First(&lesson, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("lesson not found")
		}
		return err
	}

	return s.db.Delete(&lesson).Error
}

func (s *LessonService) Restore(ctx context.Context, id uuid.UUID) (*dto.LessonResponse, error) {
	var lesson models.Lesson
	if err := s.db.Unscoped().First(&lesson, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("lesson not found")
		}
		return nil, err
	}

	if lesson.DeletedAt.Time.IsZero() {
		return nil, errors.New("lesson is not deleted")
	}

	if err := s.db.Model(&lesson).Unscoped().Update("deleted_at", nil).Error; err != nil {
		return nil, err
	}

	response := s.toLessonResponse(&lesson)
	return &response, nil
}

func (s *LessonService) PermanentDelete(ctx context.Context, id uuid.UUID) error {
	var lesson models.Lesson
	if err := s.db.Unscoped().First(&lesson, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("lesson not found")
		}
		return err
	}

	return s.db.Unscoped().Delete(&lesson).Error
}


func (s *LessonService) BulkCreate(ctx context.Context, moduleID uuid.UUID, req *dto.BulkCreateLessonRequest) ([]dto.LessonResponse, error) {
	var module models.Module
	if err := s.db.First(&module, "id = ?", moduleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("module not found")
		}
		return nil, err
	}

	lessons := make([]models.Lesson, len(req.Lessons))
	for i, item := range req.Lessons {
		lessons[i] = models.Lesson{
			ModuleID:     moduleID,
			VideoID:      item.VideoID,
			Name:         item.Name,
			VideoDetails: item.VideoDetails,
			CourseOrder:  item.CourseOrder,
		}
	}

	if err := s.db.Create(&lessons).Error; err != nil {
		return nil, err
	}

	return s.toLessonResponses(lessons), nil
}

func (s *LessonService) Reorder(ctx context.Context, moduleID uuid.UUID, req *dto.ReorderLessonsRequest) error {
	var module models.Module
	if err := s.db.First(&module, "id = ?", moduleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("module not found")
		}
		return err
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		for _, order := range req.LessonOrders {
			var lesson models.Lesson
			if err := tx.First(&lesson, "id = ? AND module_id = ?", order.LessonID, moduleID).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return fmt.Errorf("lesson %s not found in module %s", order.LessonID, moduleID)
				}
				return err
			}

			if err := tx.Model(&lesson).Update("course_order", order.CourseOrder).Error; err != nil {
				return err
			}
		}
		return nil
	})
}