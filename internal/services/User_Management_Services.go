package services

import (
	"course_online_backend/internal/models"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{DB: db}
}

func (s *UserService) preloadAll(tx *gorm.DB) *gorm.DB {
	return tx.
		Preload("Biodata").
		Preload("Roles").
		Preload("Activities").
		Preload("Activities.User").
		Preload("Activities.User.Biodata")
	// Preload("CreatedCourses").
	// Preload("CreatedQuizzes").
	// Preload("Enrollments").
	// Preload("Enrollments.Course").
	// Preload("QuizAttempts").
	// Preload("CompletedLessons")
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := s.preloadAll(s.DB).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserService) GetUserByID(id string) (*models.User, error) {
	var user models.User
	err := s.preloadAll(s.DB).First(&user, "id = ?", id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}


func (s *UserService) UpdateUser(id string, updatedUser models.User) (*models.User, error) {
	var user models.User

	if err := s.preloadAll(s.DB).First(&user, "id = ?", id).Error; err != nil {
		return nil, errors.New("user not found")
	}

	if updatedUser.Username != "" {
		user.Username = updatedUser.Username
	}
	if updatedUser.EmailAddress != "" {
		user.EmailAddress = updatedUser.EmailAddress
	}

	if updatedUser.IsActive != user.IsActive {
		user.IsActive = updatedUser.IsActive
	}

	if err := s.DB.Save(&user).Error; err != nil {
		return nil, err
	}

	if len(updatedUser.Roles) > 0 {
		if err := s.DB.Model(&user).Association("Roles").Replace(updatedUser.Roles); err != nil {
			return nil, err
		}
	}

	if updatedUser.Biodata.ID != uuid.Nil {
		if err := s.DB.Model(&user.Biodata).Updates(updatedUser.Biodata).Error; err != nil {
			return nil, err
		}
	}

	if err := s.preloadAll(s.DB).First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}


func (s *UserService) DeleteUser(id string) error {
	if err := s.DB.Delete(&models.User{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}
