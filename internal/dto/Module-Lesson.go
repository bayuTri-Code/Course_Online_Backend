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

type CreateLessonRequest struct {
	ModuleID     uuid.UUID `json:"-" swaggerignore:"true"`
	VideoID      string    `json:"video_id" binding:"required,min=1,max=100"`
	Name         string    `json:"name" binding:"required,min=3,max=200"`
	VideoDetails string    `json:"video_details" binding:"omitempty"`
	CourseOrder  int       `json:"course_order" binding:"required,min=1"`
}

type UpdateLessonRequest struct {
	VideoID      *string `json:"video_id" binding:"omitempty,min=1,max=100"`
	Name         *string `json:"name" binding:"omitempty,min=3,max=200"`
	VideoDetails *string `json:"video_details" binding:"omitempty"`
	CourseOrder  *int    `json:"course_order" binding:"omitempty,min=1"`
}

type BulkCreateLessonRequest struct {
	Lessons []CreateLessonItem `json:"lessons" binding:"required,dive"`
}

type CreateLessonItem struct {
	VideoID      string `json:"video_id" binding:"required,min=1,max=100"`
	Name         string `json:"name" binding:"required,min=3,max=200"`
	VideoDetails string `json:"video_details" binding:"omitempty"`
	CourseOrder  int    `json:"course_order" binding:"required,min=1"`
}

type ReorderLessonsRequest struct {
	LessonOrders []LessonOrderItem `json:"lesson_orders" binding:"required,dive"`
}

type LessonOrderItem struct {
	LessonID    uuid.UUID `json:"lesson_id" binding:"required"`
	CourseOrder int       `json:"course_order" binding:"required,min=1"`
}

// ========== RESPONSE DTOs ==========

type LessonResponse struct {
	ID           uuid.UUID  `json:"id"`
	ModuleID     uuid.UUID  `json:"module_id"`
	VideoID      string     `json:"video_id"`
	Name         string     `json:"name"`
	VideoDetails string     `json:"video_details"`
	CourseOrder  int        `json:"course_order"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

type LessonDetailResponse struct {
	ID           uuid.UUID       `json:"id"`
	ModuleID     uuid.UUID       `json:"module_id"`
	VideoID      string          `json:"video_id"`
	Name         string          `json:"name"`
	VideoDetails string          `json:"video_details"`
	CourseOrder  int             `json:"course_order"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	DeletedAt    *time.Time      `json:"deleted_at,omitempty"`
	Module       *ModuleBasicInfo `json:"module,omitempty"`
}

type ModuleBasicInfo struct {
	ID       uuid.UUID `json:"id"`
	CourseID uuid.UUID `json:"course_id"`
	Name     string    `json:"name"`
	Number   int       `json:"number"`
}

type LessonListResponse struct {
	Lessons    []LessonResponse `json:"lessons"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
	TotalPages int              `json:"total_pages"`
}
