package ModuleManagementServices

import (
	"context"
	"course_online_backend/internal/dto"
	"course_online_backend/internal/models"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ModuleService struct {
	db *gorm.DB
}

func NewModuleService(db *gorm.DB) *ModuleService {
	return &ModuleService{db: db}
}

func (s *ModuleService) Create(ctx context.Context, req *dto.CreateModuleRequest) (*dto.ModuleResponse, error) {
	var course models.Course
	if err := s.db.First(&course, "id = ?", req.CourseID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("course not found")
		}
		return nil, err
	}

	module := &models.Module{
		CourseID: req.CourseID,
		Name:     req.Name,
		Number:   req.Number,
	}

	if err := s.db.Create(module).Error; err != nil {
		return nil, err
	}

	return &dto.ModuleResponse{
		ID:        module.ID,
		CourseID:  module.CourseID,
		Name:      module.Name,
		Number:    module.Number,
		CreatedAt: module.CreatedAt,
		UpdatedAt: module.UpdatedAt,
	}, nil
}

func (s *ModuleService) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateModuleRequest) (*dto.ModuleResponse, error) {
	var module models.Module
	if err := s.db.First(&module, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("module not found")
		}
		return nil, err
	}

	if req.Name != nil {
		module.Name = *req.Name
	}
	if req.Number != nil {
		module.Number = *req.Number
	}

	if err := s.db.Save(&module).Error; err != nil {
		return nil, err
	}

	return &dto.ModuleResponse{
		ID:        module.ID,
		CourseID:  module.CourseID,
		Name:      module.Name,
		Number:    module.Number,
		CreatedAt: module.CreatedAt,
		UpdatedAt: module.UpdatedAt,
	}, nil
}

func (s *ModuleService) GetByID(ctx context.Context, id uuid.UUID) (*dto.ModuleDetailResponse, error) {
	var module models.Module
	if err := s.db.Preload("Lessons").First(&module, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("module not found")
		}
		return nil, err
	}

	lessons := make([]dto.LessonResponse, len(module.Lessons))
	for i, lesson := range module.Lessons {
		lessons[i] = dto.LessonResponse{
			ID:           lesson.ID,
			ModuleID:     lesson.ModuleID,
			VideoID:      lesson.VideoID,
			Name:         lesson.Name,
			VideoDetails: lesson.VideoDetails,
			CourseOrder:  lesson.CourseOrder,
			CreatedAt:    lesson.CreatedAt,
			UpdatedAt:    lesson.UpdatedAt,
		}
	}

	return &dto.ModuleDetailResponse{
		ID:        module.ID,
		CourseID:  module.CourseID,
		Name:      module.Name,
		Number:    module.Number,
		Lessons:   lessons,
		CreatedAt: module.CreatedAt,
		UpdatedAt: module.UpdatedAt,
	}, nil
}

func (s *ModuleService) GetByCourse(ctx context.Context, courseID uuid.UUID) ([]dto.ModuleResponse, error) {
	var course models.Course
	if err := s.db.First(&course, "id = ?", courseID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("course not found")
		}
		return nil, err
	}

	var modules []models.Module
	if err := s.db.Where("course_id = ?", courseID).Order("number ASC").Find(&modules).Error; err != nil {
		return nil, err
	}

	responses := make([]dto.ModuleResponse, len(modules))
	for i, module := range modules {
		responses[i] = dto.ModuleResponse{
			ID:        module.ID,
			CourseID:  module.CourseID,
			Name:      module.Name,
			Number:    module.Number,
			CreatedAt: module.CreatedAt,
			UpdatedAt: module.UpdatedAt,
		}
	}

	return responses, nil
}

func (s *ModuleService) SoftDelete(ctx context.Context, id uuid.UUID) error {
	var module models.Module
	if err := s.db.First(&module, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("module not found")
		}
		return err
	}

	return s.db.Delete(&module).Error
}

func (s *ModuleService) Restore(ctx context.Context, id uuid.UUID) (*dto.ModuleResponse, error) {
	var module models.Module
	if err := s.db.Unscoped().First(&module, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("module not found")
		}
		return nil, err
	}

	if module.DeletedAt.Time.IsZero() {
		return nil, errors.New("module is not deleted")
	}

	if err := s.db.Model(&module).Unscoped().Update("deleted_at", nil).Error; err != nil {
		return nil, err
	}

	return &dto.ModuleResponse{
		ID:        module.ID,
		CourseID:  module.CourseID,
		Name:      module.Name,
		Number:    module.Number,
		CreatedAt: module.CreatedAt,
		UpdatedAt: module.UpdatedAt,
	}, nil
}

func (s *ModuleService) PermanentDelete(ctx context.Context, id uuid.UUID) error {
	var module models.Module
	if err := s.db.Unscoped().First(&module, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("module not found")
		}
		return err
	}

	return s.db.Unscoped().Delete(&module).Error
}