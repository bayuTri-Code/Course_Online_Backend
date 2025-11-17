package CourseManagementServices

import (
	"course_online_backend/internal/dto"
	"course_online_backend/internal/models"

	"gorm.io/gorm"
)

type CourseTypeService struct {
	db *gorm.DB
}

func NewCourseTypeService(db *gorm.DB) *CourseTypeService {
	return &CourseTypeService{db: db}
}

func (s *CourseTypeService) Create(req *dto.CreateCourseTypeRequest) (*dto.CourseTypeResponse, error) {
	courseType := &models.CourseType{
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   req.CreatedAt,
		UpdatedAt:   req.UpdatedAt,
	}

	if err := s.db.Create(courseType).Error; err != nil {
		return nil, err
	}

	return &dto.CourseTypeResponse{
		ID:          courseType.ID,
		Name:        courseType.Name,
		Description: courseType.Description,
		CreatedAt:   courseType.CreatedAt,
		UpdatedAt:   courseType.UpdatedAt,
	}, nil
}

func (s *CourseTypeService) GetAll() ([]dto.CourseTypeResponse, error) {
	var courseTypes []models.CourseType
	if err := s.db.Order("name ASC").Find(&courseTypes).Error; err != nil {
		return nil, err
	}

	responses := make([]dto.CourseTypeResponse, len(courseTypes))
	for i, ct := range courseTypes {
		responses[i] = dto.CourseTypeResponse{
			ID:          ct.ID,
			Name:        ct.Name,
			Description: ct.Description,
			CreatedAt:   ct.CreatedAt,
			UpdatedAt:   ct.UpdatedAt,
		}
	}

	return responses, nil
}