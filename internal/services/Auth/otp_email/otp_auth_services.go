package otpemail

import (
	"context"
	"course_online_backend/internal/models"
	"errors"
	"fmt"
	"time"

	firebase "firebase.google.com/go/v4"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OTPAuthService struct {
	DB          *gorm.DB
	FirebaseApp *firebase.App
	OTPService  *OTPService
}

func NewOTPAuthService(db *gorm.DB, app *firebase.App, otpService *OTPService) *OTPAuthService {
	return &OTPAuthService{
		DB:          db,
		FirebaseApp: app,
		OTPService:  otpService,
	}
}

func (s *OTPAuthService) VerifyOTPAndLogin(ctx context.Context, email string, otp string) (*models.User, string, error) {
	err := s.OTPService.VerifyOTP(ctx, email, otp)
	if err != nil {
		return nil, "", errors.New("invalid or expired OTP")
	}

	authClient, err := s.FirebaseApp.Auth(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("firebase error: %v", err)
	}
	var user models.User
	err = s.DB.Where("email_address = ?", email).Preload("Roles").First(&user).Error

	if err == gorm.ErrRecordNotFound {

		firebaseUID := uuid.New().String()

		user = models.User{
			ID:           uuid.New(),
			FirebaseUID:  firebaseUID,
			EmailAddress: email,
			Username:     email,
			// Roles: user.Roles,
			IsActive:     true,
			CreatedAt:    time.Now(),
			LastLogin:    time.Now(),
		}

		var defaultRole models.Role
		if err := s.DB.Where("name = ?", "student").First(&defaultRole).Error; err != nil {
			return nil, "", fmt.Errorf("default role not found: %v", err)
		}

		if err := s.DB.Transaction(func(tx *gorm.DB) error {

			if err := tx.Create(&user).Error; err != nil {
				return err
			}

			if err := tx.Model(&user).Association("Roles").Append(&defaultRole); err != nil {
				return err
			}

			return nil
		}); err != nil {
			return nil, "", fmt.Errorf("failed to create user: %v", err)
		}
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, "", err
	}

	user.LastLogin = time.Now()
	_ = s.DB.Save(&user)

	customToken, err := authClient.CustomToken(ctx, user.FirebaseUID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate firebase custom token: %v", err)
	}

	return &user, customToken, nil
}
