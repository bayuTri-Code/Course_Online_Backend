package services

import (
	"course_online_backend/internal/dto"
	"course_online_backend/internal/models"
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

func (s *ActivityService) LogActivity(userID uuid.UUID, activityName string) error {
	activity := models.Activity{
		ID:           uuid.New(),
		UserID:       userID,
		ActivityName: activityName,
		When:         time.Now(),
	}

	return s.DB.Create(&activity).Error
}

func (s *ActivityService) GetActivity(userID uuid.UUID) ([]dto.ActivityResponse, error) {
	var activities []models.Activity
	if err := s.DB.
		Preload("User").
		Preload("User.Roles").
		Preload("User.Biodata").
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

func (s *ActivityService) GetUserByFirebaseUID(firebaseUID string) (*models.User, error) {
	var user models.User
	if err := s.DB.Where("firebase_uid = ?", firebaseUID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
