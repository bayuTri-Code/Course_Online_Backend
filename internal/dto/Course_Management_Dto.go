package dto

import (
	"github.com/google/uuid"
	"time"
)

type CreateCourseRequest struct {
	Name              string    `json:"name" binding:"required,min=3,max=200"`
	Description       string    `json:"description" binding:"required"`
	Price             float64   `json:"price" binding:"required,min=0"`
	IsProgressLimited bool      `json:"is_progress_limited"`
	CourseTypeID      string `form:"course_type_id" json:"course_type_id" binding:"required"`
}

type UpdateCourseRequest struct {
	Name              string    `json:"name" binding:"omitempty,min=3,max=200"`
	Description       string    `json:"description"`
	Price             float64   `json:"price" binding:"omitempty,min=0"`
	IsProgressLimited *bool     `json:"is_progress_limited"`
	CourseTypeID      uuid.UUID `json:"course_type_id"`
}

type CourseQueryParams struct {
	Page           int       `form:"page" binding:"omitempty,min=1"`
	Limit          int       `form:"limit" binding:"omitempty,min=1,max=100"`
	Search         string    `form:"search"`
	CourseTypeID   uuid.UUID `form:"course_type_id"`
	CreatedBy      uuid.UUID `form:"created_by"`
	MinPrice       float64   `form:"min_price" binding:"omitempty,min=0"`
	MaxPrice       float64   `form:"max_price" binding:"omitempty,min=0"`
	SortBy         string    `form:"sort_by" binding:"omitempty,oneof=name price created_at"`
	SortOrder      string    `form:"sort_order" binding:"omitempty,oneof=asc desc"`
	IncludeDeleted bool      `form:"include_deleted"`
}

type CourseResponse struct {
	ID                uuid.UUID           `json:"id"`
	Name              string              `json:"name"`
	Description       string              `json:"description"`
	Thumbnail         string              `json:"thumbnail"`
	Price             float64             `json:"price"`
	IsProgressLimited bool                `json:"is_progress_limited"`
	CourseTypeID      uuid.UUID           `json:"course_type_id"`
	CreatedBy         *uuid.UUID          `json:"created_by"`
	CreatedAt         time.Time           `json:"created_at"`
	UpdatedAt         time.Time           `json:"updated_at"`
	DeletedAt         *time.Time          `json:"deleted_at,omitempty"`
	CourseType        *CourseTypeResponse `json:"course_type,omitempty"`
	Creator           *CreatorResponse    `json:"creator,omitempty"`
	ModulesCount      int                 `json:"modules_count,omitempty"`
	LessonsCount      int                 `json:"lessons_count,omitempty"`
	EnrollmentsCount  int                 `json:"enrollments_count,omitempty"`
}

type CourseDetailResponse struct {
	ID                uuid.UUID           `json:"id"`
	Name              string              `json:"name"`
	Description       string              `json:"description"`
	Thumbnail         string              `json:"thumbnail"`
	Price             float64             `json:"price"`
	IsProgressLimited bool                `json:"is_progress_limited"`
	CourseTypeID      uuid.UUID           `json:"course_type_id"`
	CreatedBy         *uuid.UUID          `json:"created_by"`
	CreatedAt         time.Time           `json:"created_at"`
	UpdatedAt         time.Time           `json:"updated_at"`
	DeletedAt         *time.Time          `json:"deleted_at,omitempty"`
	CourseType        *CourseTypeResponse `json:"course_type,omitempty"`
	Creator           *CreatorResponse    `json:"creator,omitempty"`
	Modules           []ModuleResponse    `json:"modules,omitempty"`
	EnrollmentsCount  int                 `json:"enrollments_count"`
}

type CreateCourseTypeRequest struct {
	Name        string    `json:"name" binding:"required,min=3,max=100"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CourseTypeResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreatorResponse struct {
	ID       uuid.UUID `json:"id"`
	FullName string    `json:"full_name"`
	Email    string    `json:"email"`
}

type ModuleResponse struct {
	ID           uuid.UUID        `json:"id"`
	Name         string           `json:"name"`
	Number       int              `json:"number"`
	LessonsCount int              `json:"lessons_count"`
	Lessons      []LessonResponse `json:"lessons,omitempty"`
}

type LessonResponse struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	VideoID      string    `json:"video_id"`
	VideoDetails string    `json:"video_details"`
	CourseOrder  int       `json:"course_order"`
}

type PaginationResponse struct {
	Total       int64       `json:"total"`
	Page        int         `json:"page"`
	Limit       int         `json:"limit"`
	TotalPages  int         `json:"total_pages"`
	HasNext     bool        `json:"has_next"`
	HasPrevious bool        `json:"has_previous"`
	Data        interface{} `json:"data"`
}
