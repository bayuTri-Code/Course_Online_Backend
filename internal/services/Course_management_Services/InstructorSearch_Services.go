package CourseManagementServices

import (
	"context"
	"course_online_backend/internal/dto"
	"course_online_backend/internal/models"
	"errors"
	"gorm.io/gorm"
	"github.com/google/uuid"
)

type InstructorService struct {
	db *gorm.DB
}

func NewInstructorService(db *gorm.DB) *InstructorService {
	return &InstructorService{db: db}
}


func (s *InstructorService) SearchInstructors(ctx context.Context, query string, limit int) ([]dto.InstructorResponse, error) {
	if limit < 1 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	var users []models.User
	
	err := s.db.
		Preload("Biodata").
		Joins("JOIN user_roles ON user_roles.user_id = users.id").
		Joins("JOIN roles ON roles.id = user_roles.role_id").
		Where("roles.name = ?", "instructor").
		Where("users.is_active = ?", true).
		Where("users.username ILIKE ? OR users.email_address ILIKE ?", "%"+query+"%", "%"+query+"%").
		Limit(limit).
		Find(&users).Error

	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return []dto.InstructorResponse{}, nil
	}

	instructors := make([]dto.InstructorResponse, len(users))
	for i, user := range users {
		avatar := ""
		if user.Biodata.ProfilePicture != "" {
			avatar = user.Biodata.ProfilePicture
		}
		
		instructors[i] = dto.InstructorResponse{
			ID:       user.ID,
			FullName: user.Username,
			Email:    user.EmailAddress,
			Avatar:   avatar,
		}
	}

	return instructors, nil
}

func (s *InstructorService) ValidateInstructor(ctx context.Context, instructorID uuid.UUID) (*models.User, error) {
	var user models.User
	
	err := s.db.
		Preload("Roles").
		Preload("Biodata").
		First(&user, "id = ? AND is_active = ?", instructorID, true).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("instructor not found or inactive")
		}
		return nil, err
	}

	// Check if user has instructor role
	if !user.HasRole("instructor") {
		return nil, errors.New("user is not an instructor")
	}

	return &user, nil
}

func (s *InstructorService) GetInstructorByID(ctx context.Context, instructorID uuid.UUID) (*dto.InstructorResponse, error) {
	user, err := s.ValidateInstructor(ctx, instructorID)
	if err != nil {
		return nil, err
	}

	avatar := ""
	if user.Biodata.ProfilePicture != "" {
		avatar = user.Biodata.ProfilePicture
	}

	return &dto.InstructorResponse{
		ID:       user.ID,
		FullName: user.Username,
		Email:    user.EmailAddress,
		Avatar:   avatar,
	}, nil
}