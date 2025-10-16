package dto

import (
	"course_online_backend/internal/models"
	"time"
)

type FirebaseLoginRequest struct {
	Token string `json:"token"`
}
type CreateBiodataRequest struct {
	Name   string `json:"name" binding:"required"`
	Age    int    `json:"age" binding:"required"`
	School string `json:"school" binding:"required"`
}

type UpdateBiodataRequest struct {
	Name   string `json:"name"`
	Age    int    `json:"age"`
	School string `json:"school"`
}

type BiodataResponse struct {
	ID             string `json:"id"`
	UserID         string `json:"user_id"`
	Name           string `json:"name"`
	Age            int    `json:"age"`
	School         string `json:"school"`
	ProfilePicture string `json:"profile_picture"`
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
