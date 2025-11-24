package auth

import (
	"context"
	"course_online_backend/internal/models"
	"fmt"
	"time"
	"gorm.io/gorm"

	firebase "firebase.google.com/go/v4"
)

type AuthService struct {
	Db          *gorm.DB
	FirebaseApp *firebase.App
}

func NewAuthService(db *gorm.DB, app *firebase.App) *AuthService {
	return &AuthService{
		Db:          db,
		FirebaseApp: app,
	}
}
func (s *AuthService) LoginWithFirebaseToken(tokenString string) (*models.User, error) {
	authClient, err := s.FirebaseApp.Auth(context.Background())
	if err != nil {
		return nil, err
	}

	token, err := authClient.VerifyIDToken(context.Background(), tokenString)
	if err != nil {
		return nil, err
	}

	uid := token.UID
	email := token.Claims["email"].(string)
	name := ""
	if n, ok := token.Claims["name"].(string); ok {
		name = n
	}

	var user models.User
	err = s.Db.Where("firebase_uid = ?", uid).Preload("Roles").First(&user).Error
	if err == nil {
		user.LastLogin = time.Now()
		if saveErr := s.Db.Save(&user).Error; saveErr != nil {
			return nil, fmt.Errorf("failed to update last login: %v", saveErr)
		}
		return &user, nil
	}

	if err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to fetch user: %v", err)
	}

	var defaultRole models.Role
	if err := s.Db.Where("name = ?", "student").First(&defaultRole).Error; err != nil {
		return nil, fmt.Errorf("default role not found: %v", err)
	}

	user = models.User{
		FirebaseUID:  uid,
		EmailAddress: email,
		Username:     name,
		IsActive:     true,
		CreatedAt:    time.Now(),
		LastLogin:    time.Now(),
	}

	if err := s.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		if err := tx.Model(&user).Association("Roles").Append(&defaultRole); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to register user: %v", err)
	}

	return &user, nil
}
