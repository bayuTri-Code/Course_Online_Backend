package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Enrollment struct {
	ID                 uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CourseID           uuid.UUID      `gorm:"type:uuid;not null;index;uniqueIndex:idx_user_course" json:"course_id"`
	UserID             uuid.UUID      `gorm:"type:uuid;not null;index;uniqueIndex:idx_user_course" json:"user_id"`
	EnrollmentDatetime time.Time      `gorm:"not null" json:"enrollment_datetime"`
	CompletedDatetime  *time.Time     `json:"completed_datetime,omitempty"`
	Status             string         `gorm:"type:varchar(20);not null;default:'active';index" json:"status"`
	StatusPayment      string         `gorm:"type:varchar(20);not null;default:'free';index" json:"status_payment"`
	ExpiredDate        *time.Time     `json:"expired_date,omitempty"`
	Progress           float64        `gorm:"type:decimal(5,2);default:0;check:progress >= 0 AND progress <= 100" json:"progress"`
	CreatedAt          time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`

	Course Course `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE" json:"course,omitempty"`
	User   User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
}

func (Enrollment) TableName() string {
	return "enrollments"
}

type UserLesson struct {
	ID                uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID            uuid.UUID      `gorm:"type:uuid;not null;index;uniqueIndex:idx_user_lesson" json:"user_id"`
	LessonID          uuid.UUID      `gorm:"type:uuid;not null;index;uniqueIndex:idx_user_lesson" json:"lesson_id"`
	IsCompleted       bool           `gorm:"default:false;index" json:"is_completed"`
	CompletedDatetime *time.Time     `json:"completed_datetime,omitempty"`
	LastAccessedAt    time.Time      `gorm:"not null;index" json:"last_accessed_at"`
	TimeSpent         int            `gorm:"default:0" json:"time_spent"`
	AttemptCount      int            `gorm:"default:0" json:"attempt_count"`
	LastPosition      *int           `json:"last_position,omitempty"`
	CreatedAt         time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`

	User   User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Lesson Lesson `gorm:"foreignKey:LessonID;constraint:OnDelete:CASCADE" json:"lesson,omitempty"`
}

func (UserLesson) TableName() string {
	return "user_lessons"
}
