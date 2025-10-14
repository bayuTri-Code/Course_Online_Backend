
package models

import (
	"github.com/google/uuid"
	"time"
)

type Enrollment struct {
	ID                 uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CourseID           uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex:idx_user_course,composite:user_course" json:"course_id"`
	UserID             uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex:idx_user_course,composite:user_course" json:"user_id"`
	EnrollmentDatetime time.Time  `gorm:"not null" json:"enrollment_datetime"`
	CompletedDatetime  *time.Time `json:"completed_datetime"`

	Course Course `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE" json:"course"`
	User   User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user"`
}

type UserLesson struct {
	ID                uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID            uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_user_lesson,composite:user_lesson" json:"user_id"`
	LessonID          uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_user_lesson,composite:user_lesson" json:"lesson_id"`
	CompletedDatetime time.Time `gorm:"not null" json:"completed_datetime"`

	User   User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user"`
	Lesson Lesson `gorm:"foreignKey:LessonID;constraint:OnDelete:CASCADE" json:"lesson"`
}
