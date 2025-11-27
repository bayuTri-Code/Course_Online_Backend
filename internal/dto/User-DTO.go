package dto

import (
	"course_online_backend/internal/models"
	"time"
)

// login dto
type FirebaseLoginRequest struct {
	Token string `json:"token"`
}

type UserResponseDTO struct {
	ID         string      `json:"id"`
	FirebaseID string      `json:"firebaseId"`
	Email      string      `json:"email"`
	Username   string      `json:"username"`
	Roles      models.Role `json:"roles"`
	IsActive   bool        `json:"isActive"`
	CreatedAt  time.Time   `json:"createdAt"`
	LastLogin  time.Time   `json:"lastLogin"`
}

type UserLoginResponse struct {
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Data    UserResponseDTO `json:"data"`
}

//otp dto
type SendOTPResponse struct {
	Status  bool   `json:"status" example:"true"`
	Message string `json:"message" example:"OTP sent successfully"`
	Data    struct {
		Email string `json:"email" example:"user@example.com"`
	} `json:"data"`
}

type ErrorResponse struct {
	Status  bool   `json:"status" example:"false"`
	Message string `json:"message" example:"Invalid email format"`
}

type VerifyOTPResponse struct {
	Status  bool   `json:"status" example:"true"`
	Message string `json:"message" example:"Login success"`
	Data    struct {
		CustomToken string      `json:"customToken" example:"xxxxx.yyyyy.zzzzz"`
		User        UserProfile `json:"user"`
	} `json:"data"`
}

type UserProfile struct {
	ID          string      `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	FirebaseUID string      `json:"firebaseId" example:"firebase-uid-123"`
	Email       string      `json:"email" example:"user@example.com"`
	Username    string      `json:"username" example:"bayutri"`
	Roles       interface{} `json:"roles"`
	IsActive    bool        `json:"isActive" example:"true"`
	CreatedAt   string      `json:"createdAt" example:"2025-01-01T12:00:00Z"`
	LastLogin   string      `json:"lastLogin" example:"2025-01-20T11:00:00Z"`
}


// biodata dto
type CreateBiodataRequest struct {
	Name   string `form:"name" binding:"required"`
	Age    int    `form:"age" binding:"required"`
	School string `form:"school" binding:"required"`
}

type UpdateBiodataRequest struct {
	Name   string `form:"name"`
	Age    int    `form:"age"`
	School string `form:"school"`
}

type BiodataResponse struct {
	ID             string `json:"id"`
	UserID         string `json:"user_id"`
	Name           string `json:"name"`
	Age            int    `json:"age"`
	School         string `json:"school"`
	ProfilePicture string `json:"profile_picture"`
}

type BiodataResponseForAct struct {
	Name           string `json:"name"`
	Age            int    `json:"age"`
	School         string `json:"school"`
	ProfilePicture string `json:"profile_picture"`
}

type BaseResponseBiodata struct {
	Status  bool            `json:"status" example:"true"`
	Message string          `json:"message"`
	Data    BiodataResponse `json:"data"`
}



//user Management
type UpdateUserRequest struct {
	Username     string `json:"name"`
	EmailAddress string `json:"email"`
	IsActive     *bool   `json:"is_active"`
	RoleIDs      []string `json:"role_ids" example:"role_id"`
}

type BaseResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type BaseResponseDelete struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
}


