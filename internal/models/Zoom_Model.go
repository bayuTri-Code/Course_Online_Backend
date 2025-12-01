package models

import (
	"time"

	"github.com/google/uuid"
)

type Zoom struct {
	ID          uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Title       string     `gorm:"size:255;not null" json:"title"` 
	Link        string     `gorm:"size:500;not null" json:"link"`
	Description string     `gorm:"type:text" json:"description"`
	ScheduledAt *time.Time `gorm:"type:timestamp" json:"scheduled_at"`
	Duration    int        `gorm:"type:int" json:"duration"` 
	CourseID    uuid.UUID  `gorm:"type:uuid;not null;index" json:"course_id"` 
	CreatedAt   time.Time  `gorm:"type:timestamp;default:current_timestamp" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"type:timestamp;default:current_timestamp" json:"updated_at"`
	
	Course Course `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE" json:"course"`
}