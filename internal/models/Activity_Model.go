
package models

import (
	"github.com/google/uuid"
	"time"
)

type Activity struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ActivityName string    `gorm:"size:100;not null" json:"activity_name"`
	UserID       uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	When         time.Time `gorm:"not null" json:"when"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user"`
}