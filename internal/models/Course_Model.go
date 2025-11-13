package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type CourseType struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Courses []Course `gorm:"foreignKey:CourseTypeID;constraint:OnDelete:SET NULL" json:"courses,omitempty"`
}

type Course struct {
	ID                uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name              string         `gorm:"size:200;not null" json:"name"`
	Description       string         `gorm:"type:text" json:"description"`
	Thumbnail         string         `gorm:"type:text" json:"thumbnail"`
	Price             float64        `gorm:"type:decimal(10,2);check:price >= 0" json:"price"`
	IsProgressLimited bool           `gorm:"default:false" json:"is_progress_limited"`
	CourseTypeID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"course_type_id"`
	CreatedBy         *uuid.UUID     `gorm:"type:uuid;index" json:"created_by"`
	CreatedAt         time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty" swaggertype:"string" format:"date-time"`

	CourseType  *CourseType  `gorm:"foreignKey:CourseTypeID;constraint:OnDelete:RESTRICT" json:"course_type,omitempty"`
	Creator     *User        `gorm:"foreignKey:CreatedBy;constraint:OnDelete:SET NULL" json:"creator,omitempty"`
	Modules     []Module     `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE" json:"modules,omitempty"`
	Quizzes     []Quiz       `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE" json:"quizzes,omitempty"`
	Enrollments []Enrollment `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE" json:"enrollments,omitempty"`
	Zoom        *Zoom        `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE" json:"zoom,omitempty"`
}

type Module struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CourseID  uuid.UUID      `gorm:"type:uuid;not null;index" json:"course_id"`
	Name      string         `gorm:"size:200;not null" json:"name"`
	Number    int            `gorm:"not null" json:"number"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Course  Course   `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE" json:"course,omitempty"`
	Lessons []Lesson `gorm:"foreignKey:ModuleID;constraint:OnDelete:CASCADE" json:"lessons,omitempty"`
}

type Lesson struct {
	ID           uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ModuleID     uuid.UUID      `gorm:"type:uuid;not null;index" json:"module_id"`
	VideoID      string         `gorm:"size:100" json:"video_id"`
	Name         string         `gorm:"size:200;not null" json:"name"`
	VideoDetails string         `gorm:"type:text" json:"video_details"`
	CourseOrder  int            `json:"course_order"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Module          Module       `gorm:"foreignKey:ModuleID;constraint:OnDelete:CASCADE" json:"module,omitempty"`
	UserCompletions []UserLesson `gorm:"foreignKey:LessonID;constraint:OnDelete:CASCADE" json:"user_completions,omitempty"`
}
