package otpemail

import (
	"context"
	"course_online_backend/internal/models"
	"errors"
	"fmt"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
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
		return s.createNewUser(ctx, authClient, email)
	}

	if err != nil {
		return nil, "", fmt.Errorf("database error: %v", err)
	}

	return s.loginExistingUser(ctx, authClient, &user)
}

func (s *OTPAuthService) createNewUser(ctx context.Context, authClient *auth.Client, email string) (*models.User, string, error) {
	firebaseUser, err := authClient.CreateUser(ctx, (&auth.UserToCreate{}).
		Email(email).
		EmailVerified(true).
		Disabled(false))

	if err != nil {
		if auth.IsEmailAlreadyExists(err) {
			existingFirebaseUser, getErr := authClient.GetUserByEmail(ctx, email)
			if getErr != nil {
				return nil, "", fmt.Errorf("user exists in Firebase but cannot retrieve: %v", getErr)
			}
			firebaseUser = existingFirebaseUser
		} else {
			return nil, "", fmt.Errorf("failed to create Firebase user: %v", err)
		}
	}

	user := models.User{
		ID:           uuid.New(),
		FirebaseUID:  firebaseUser.UID,
		EmailAddress: email,
		Username:     email,
		IsActive:     true,
		CreatedAt:    time.Now(),
		LastLogin:    time.Now(),
	}

	var defaultRole models.Role
	if err := s.DB.Where("name = ?", "student").First(&defaultRole).Error; err != nil {
		_ = authClient.DeleteUser(ctx, firebaseUser.UID)
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
		_ = authClient.DeleteUser(ctx, firebaseUser.UID)
		return nil, "", fmt.Errorf("failed to create user in database: %v", err)
	}

	if err := s.DB.Preload("Roles").First(&user, user.ID).Error; err != nil {
		return nil, "", fmt.Errorf("failed to load user roles: %v", err)
	}

	customToken, err := authClient.CustomToken(ctx, user.FirebaseUID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate firebase custom token: %v", err)
	}

	return &user, customToken, nil
}

func (s *OTPAuthService) loginExistingUser(ctx context.Context, authClient *auth.Client, user *models.User) (*models.User, string, error) {
	_, err := authClient.GetUser(ctx, user.FirebaseUID)
	if err != nil {
		if auth.IsUserNotFound(err) {
			recreatedUser, createErr := authClient.CreateUser(ctx, (&auth.UserToCreate{}).
				UID(user.FirebaseUID).
				Email(user.EmailAddress).
				EmailVerified(true).
				Disabled(false))

			if createErr != nil {
				if auth.IsUIDAlreadyExists(createErr) {
					newUser, newErr := authClient.CreateUser(ctx, (&auth.UserToCreate{}).
						Email(user.EmailAddress).
						EmailVerified(true).
						Disabled(false))

					if newErr != nil {
						return nil, "", fmt.Errorf("failed to recreate Firebase user: %v", newErr)
					}

					user.FirebaseUID = newUser.UID
					if err := s.DB.Save(user).Error; err != nil {
						return nil, "", fmt.Errorf("failed to update user Firebase UID: %v", err)
					}
				} else {
					return nil, "", fmt.Errorf("failed to recreate Firebase user: %v", createErr)
				}
			} else {
				user.FirebaseUID = recreatedUser.UID
			}
		} else {
			return nil, "", fmt.Errorf("failed to verify Firebase user: %v", err)
		}
	}

	user.LastLogin = time.Now()
	if err := s.DB.Save(user).Error; err != nil {
		return nil, "", fmt.Errorf("failed to update last login: %v", err)
	}

	customToken, err := authClient.CustomToken(ctx, user.FirebaseUID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate firebase custom token: %v", err)
	}

	return user, customToken, nil
}
