
package models

import (
	"github.com/google/uuid"
	"time"
)

type CourseType struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`

	Courses []Course `gorm:"foreignKey:CourseTypeID;constraint:OnDelete:SET NULL" json:"courses"`
}

type Course struct {
	ID                uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name              string     `gorm:"size:200;not null" json:"name"`
	Description       string     `gorm:"type:text" json:"description"`
	Price             float64    `gorm:"type:decimal(10,2)" json:"price"`
	IsProgressLimited bool       `gorm:"default:false" json:"is_progress_limited"`
	CourseTypeID      *uuid.UUID `gorm:"type:uuid;index" json:"course_type_id"`
	CreatedBy         *uuid.UUID `gorm:"type:uuid;index" json:"created_by"`
	CreatedAt         time.Time  `gorm:"autoCreateTime" json:"created_at"`

	CourseType  *CourseType  `gorm:"foreignKey:CourseTypeID;constraint:OnDelete:SET NULL" json:"course_type"`
	Creator     *User        `gorm:"foreignKey:CreatedBy;constraint:OnDelete:SET NULL" json:"creator"`
	Modules     []Module     `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE" json:"modules"`
	Quizzes     []Quiz       `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE" json:"quizzes"`
	Enrollments []Enrollment `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE" json:"enrollments"`
	Zoom        *Zoom        `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE" json:"zoom"`
}

type Module struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CourseID uuid.UUID `gorm:"type:uuid;not null;index" json:"course_id"`
	Name     string    `gorm:"size:200;not null" json:"name"`
	Number   int       `gorm:"not null" json:"number"`

	Course  Course   `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE" json:"course"`
	Lessons []Lesson `gorm:"foreignKey:ModuleID;constraint:OnDelete:CASCADE" json:"lessons"`
}

type Lesson struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ModuleID     uuid.UUID `gorm:"type:uuid;not null;index" json:"module_id"`
	VideoID      string    `gorm:"size:100" json:"video_id"`
	Name         string    `gorm:"size:200;not null" json:"name"`
	VideoDetails string    `gorm:"type:text" json:"video_details"`
	CourseOrder  int       `json:"course_order"`

	Module          Module       `gorm:"foreignKey:ModuleID;constraint:OnDelete:CASCADE" json:"module"`
	UserCompletions []UserLesson `gorm:"foreignKey:LessonID;constraint:OnDelete:CASCADE" json:"user_completions"`
}
