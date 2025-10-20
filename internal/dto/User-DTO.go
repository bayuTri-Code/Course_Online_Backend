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