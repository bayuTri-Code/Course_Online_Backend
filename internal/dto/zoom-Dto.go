package dto

import (
	"github.com/google/uuid"
	"time"
)

type CreateZoomRequest struct {
	Title       string     `json:"title" binding:"required" example:"Pertemuan 1 - Introduction"`
	Link        string     `json:"link" binding:"required,url" example:"https://zoom.us/j/123456789"`
	Description string     `json:"description" example:"Pengenalan materi dasar"`
	ScheduledAt *time.Time `json:"scheduled_at" example:"2024-01-15T10:00:00Z"`
	Duration    int        `json:"duration" example:"90"`
	CourseID    uuid.UUID  `json:"course_id" binding:"required"`
}

type UpdateZoomRequest struct {
	Title       string     `json:"title" binding:"required"`
	Link        string     `json:"link" binding:"required,url"`
	Description string     `json:"description"`
	ScheduledAt *time.Time `json:"scheduled_at"`
	Duration    int        `json:"duration"`
}

type ZoomResponse struct {
	ID          uuid.UUID           `json:"id"`
	Title       string              `json:"title"`
	Link        string              `json:"link"`
	Description string              `json:"description"`
	ScheduledAt *time.Time          `json:"scheduled_at"`
	Duration    int                 `json:"duration"`
	CourseID    uuid.UUID           `json:"course_id"`
	Course      *ZoomCourseResponse `json:"course,omitempty"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}

type ZoomCourseResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Thumbnail   string    `json:"thumbnail"`
}

type ZoomListResponse struct {
	Total int64          `json:"total"`
	Zooms []ZoomResponse `json:"zooms"`
}