package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateQuizRequest struct {
	CourseID       uuid.UUID `json:"course_id" binding:"required"`
	Name           string    `json:"name" binding:"required,max=200"`
	Number         int       `json:"number" binding:"required,min=1"`
	MinPassScore   int       `json:"min_pass_score" binding:"required,min=0,max=100"`
	IsPassRequired bool      `json:"is_pass_required"`
}

type UpdateQuizRequest struct {
	Name           string `json:"name" binding:"required,max=200"`
	Number         int    `json:"number" binding:"required,min=1"`
	MinPassScore   int    `json:"min_pass_score" binding:"required,min=0,max=100"`
	IsPassRequired bool   `json:"is_pass_required"`
}

type QuizResponse struct {
	ID             uuid.UUID  `json:"id"`
	CourseID       uuid.UUID  `json:"course_id"`
	Name           string     `json:"name"`
	Number         int        `json:"number"`
	MinPassScore   int        `json:"min_pass_score"`
	IsPassRequired bool       `json:"is_pass_required"`
	CreatedBy      *uuid.UUID `json:"created_by"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `gorm:"index" json:"deleted_at,omitempty"`
	TotalQuestions int        `json:"total_questions"`
}

type QuizDetailResponse struct {
	ID             uuid.UUID  `json:"id"`
	CourseID       uuid.UUID  `json:"course_id"`
	Name           string     `json:"name"`
	Number         int        `json:"number"`
	MinPassScore   int        `json:"min_pass_score"`
	IsPassRequired bool       `json:"is_pass_required"`
	CreatedBy      *uuid.UUID `json:"created_by"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	Questions      int        `json:"questions"`
}
