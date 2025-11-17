package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateModuleRequest struct {
	CourseID uuid.UUID `json:"-" swaggerignore:"true"`
	Name     string    `json:"name" binding:"required,min=3,max=200"`
	Number   int       `json:"number" binding:"required,min=1"`
}

type UpdateModuleRequest struct {
	Name   *string `json:"name" binding:"omitempty,min=3,max=200"`
	Number *int    `json:"number" binding:"omitempty,min=1"`
}

type CreateModuleResponse struct {
	ID           uuid.UUID `json:"id"`
	CourseID     uuid.UUID `json:"course_id"`
	Name         string    `json:"name"`
	Number       int       `json:"number"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	LessonsCount int       `json:"lessons_count"`
}

type ModuleResponse struct {
	ID           uuid.UUID        `json:"id"`
	CourseID     uuid.UUID        `json:"course_id"`
	Name         string           `json:"name"`
	Number       int              `json:"number"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
	LessonsCount int              `json:"lessons_count"`
	Lessons      []LessonResponse `json:"lessons,omitempty"`
}

type ModuleDetailResponse struct {
	ID        uuid.UUID        `json:"id"`
	CourseID  uuid.UUID        `json:"course_id"`
	Name      string           `json:"name"`
	Number    int              `json:"number"`
	Lessons   []LessonResponse `json:"lessons"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

type LessonResponse struct {
	ID           uuid.UUID `json:"id"`
	ModuleID     uuid.UUID `json:"module_id"`
	VideoID      string    `json:"video_id"`
	Name         string    `json:"name"`
	VideoDetails string    `json:"video_details"`
	CourseOrder  int       `json:"course_order"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
