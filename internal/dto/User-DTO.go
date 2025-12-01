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

type CreateBiodataRequest struct {
	Name        string `form:"name" json:"name" binding:"required"`
	FirstName   string `form:"firstName" json:"firstName"`
	LastName    string `form:"lastName" json:"lastName"`
	Description string `form:"description" json:"description"`
	Contact     string `form:"contact" json:"contact"`
	Age         int    `form:"age" json:"age" binding:"required"`
	School      string `form:"school" json:"school" binding:"required"`
}

type UpdateBiodataRequest struct {
	Name        string `form:"name" json:"name" binding:"required"`
	FirstName   string `form:"firstName" json:"firstName"`
	LastName    string `form:"lastName" json:"lastName"`
	Description string `form:"description" json:"description"`
	Contact     string `form:"contact" json:"contact"`
	Age         int    `form:"age" json:"age" binding:"required"`
	School      string `form:"school" json:"school" binding:"required"`
}

type BiodataResponse struct {
	ID             string `json:"id"`
	UserID         string `json:"user_id"`
	Name           string `json:"name"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	Description    string `json:"description"`
	Contact        string `json:"contact"`
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

type UpdateUserRequest struct {
	Username     string   `json:"name"`
	EmailAddress string   `json:"email"`
	IsActive     *bool    `json:"is_active"`
	RoleIDs      []string `json:"role_ids" example:"role_id"`
}

type BaseResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type BaseResponseDelete struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

type SendOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifyOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required,min=4,max=8"`
}

type SendOTPData struct {
	Email string `json:"email"`
}

type SendOTPResponse struct {
	Status  bool        `json:"status" example:"true"`
	Message string      `json:"message" example:"OTP sent successfully"`
	Data    SendOTPData `json:"data"`
}

type UserRoleResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type OTPUserResponse struct {
	ID         string             `json:"id"`
	FirebaseID string             `json:"firebaseId"`
	Email      string             `json:"email"`
	Username   string             `json:"username"`
	IsActive   bool               `json:"isActive"`
	CreatedAt  time.Time          `json:"createdAt"`
	LastLogin  time.Time          `json:"lastLogin"`
	Roles      []UserRoleResponse `json:"roles"`
}

type VerifyOTPData struct {
	CustomToken string          `json:"customToken"`
	User        OTPUserResponse `json:"user"`
}

type VerifyOTPResponse struct {
	Status  bool          `json:"status" example:"true"`
	Message string        `json:"message" example:"Login success"`
	Data    VerifyOTPData `json:"data"`
}