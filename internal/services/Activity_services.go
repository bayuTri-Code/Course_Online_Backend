package services

import (
	"course_online_backend/internal/dto"
	"course_online_backend/internal/models"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ActivityService struct {
	DB *gorm.DB
}

func NewActivityService(db *gorm.DB) *ActivityService {
	return &ActivityService{DB: db}
}

// 🔄 Fungsi helper untuk Preload semua relasi utama
func (s *ActivityService) preloadActivityAll(db *gorm.DB) *gorm.DB {
	return db.
		Preload("User").
		Preload("User.Roles").
		Preload("User.Biodata")
}

// 🧠 LogActivity - Mencatat aktivitas baru (tanpa duplikasi)
func (s *ActivityService) LogActivity(userID uuid.UUID, activityName string) error {
	// Cek apakah aktivitas sama pernah dicatat dalam 1 menit terakhir untuk user yang sama
	var count int64
	oneMinuteAgo := time.Now().Add(-1 * time.Minute)
	if err := s.DB.Model(&models.Activity{}).
		Where("user_id = ? AND activity_name = ? AND \"when\" > ?", userID, activityName, oneMinuteAgo).
		Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		// Jangan log lagi kalau baru saja dilakukan
		return nil
	}

	activity := models.Activity{
		ID:           uuid.New(),
		UserID:       userID,
		ActivityName: activityName,
		When:         time.Now(),
	}

	return s.DB.Create(&activity).Error
}

func (s *ActivityService) GetUserByFirebaseUID(firebaseUID string) (*models.User, error) {
	var user models.User
	if err := s.DB.Where("firebase_uid = ?", firebaseUID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *ActivityService) GetAllActivity(userID uuid.UUID) ([]dto.ActivityResponse, error) {
	var activities []models.Activity
	if err := s.preloadActivityAll(s.DB).
		Where("user_id = ?", userID).
		Order("\"when\" DESC").
		Find(&activities).Error; err != nil {
		return nil, err
	}

	var response []dto.ActivityResponse
	for _, act := range activities {
		var roleNames []string
		for _, role := range act.User.Roles {
			roleNames = append(roleNames, role.Name)
		}

		bio := act.User.Biodata
		biodataResponse := &dto.BiodataResponseForAct{
			Name:           bio.Name,
			Age:            bio.Age,
			School:         bio.School,
			ProfilePicture: bio.ProfilePicture,
		}

		response = append(response, dto.ActivityResponse{
			ID:           act.ID.String(),
			ActivityName: act.ActivityName,
			When:         act.When,
			User: dto.UserSimpleDTO{
				ID:       act.User.ID.String(),
				Username: act.User.Username,
				Email:    act.User.EmailAddress,
				Roles:    roleNames,
				Biodata:  biodataResponse,
			},
		})
	}

	return response, nil
}

func (s *ActivityService) GetActivityByUserID(userID uuid.UUID) ([]models.Activity, error) {
	var activities []models.Activity
	if err := s.preloadActivityAll(s.DB).
		Where("user_id = ?", userID).
		Order("\"when\" DESC").
		Find(&activities).Error; err != nil {
		return nil, err
	}
	return activities, nil
}

func (s *ActivityService) GetRecentActivities(limit int) ([]models.Activity, error) {
	var activities []models.Activity
	if err := s.preloadActivityAll(s.DB).
		Order("\"when\" DESC").
		Limit(limit).
		Find(&activities).Error; err != nil {
		return nil, err
	}
	return activities, nil
}

func (s *ActivityService) SearchActivity(keyword string) ([]models.Activity, error) {
	var activities []models.Activity
	if err := s.preloadActivityAll(s.DB).
		Where("activity_name ILIKE ?", "%"+keyword+"%").
		Order("\"when\" DESC").
		Find(&activities).Error; err != nil {
		return nil, err
	}
	return activities, nil
}

func (s *ActivityService) GetActivitySummary(userID uuid.UUID) (map[string]interface{}, error) {
	var count int64
	if err := s.DB.Model(&models.Activity{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return nil, err
	}

	var lastActivity models.Activity
	s.DB.Where("user_id = ?", userID).Order("\"when\" DESC").First(&lastActivity)

	return map[string]interface{}{
		"user_id":          userID,
		"total_activities": count,
		"last_activity":    lastActivity,
	}, nil
}

func (s *ActivityService) GetActivityByRole(roleName string) ([]models.Activity, error) {
	var activities []models.Activity
	err := s.preloadActivityAll(s.DB).
		Joins("JOIN users ON users.id = activities.user_id").
		Joins("JOIN user_roles ur ON ur.user_id = users.id").
		Joins("JOIN roles r ON r.id = ur.role_id").
		Where("r.name = ?", roleName).
		Find(&activities).Error

	if err != nil {
		return nil, errors.New("failed to fetch activity by role: " + err.Error())
	}

	return activities, nil
}
