
package models

import "github.com/google/uuid"

type Zoom struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Link     string    `gorm:"size:500;not null" json:"link"`
	CourseID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"course_id"`

	Course Course `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE" json:"course"`
}
