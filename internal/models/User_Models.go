package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	FirebaseUID  string    `gorm:"unique;not null" json:"firebase_uid"`
	Username     string    `gorm:"size:100" json:"username"`
	EmailAddress string    `gorm:"unique;not null" json:"email_address"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	LastLogin    time.Time `json:"last_login"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`

	Roles            []Role            `gorm:"many2many:user_roles;joinForeignKey:UserID;joinReferences:RoleID" json:"roles"`
	Biodata          Biodata           `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"biodata"`
	Activities       []Activity        `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"activities"`
	CreatedCourses   []Course          `gorm:"foreignKey:CreatedBy;constraint:OnDelete:SET NULL" json:"created_courses"`
	CreatedQuizzes   []Quiz            `gorm:"foreignKey:CreatedBy;constraint:OnDelete:SET NULL" json:"created_quizzes"`
	Enrollments      []Enrollment      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"enrollments"`
	QuizAttempts     []UserQuizAttempt `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"quiz_attempts"`
	CompletedLessons []UserLesson      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"completed_lessons"`
}

type Role struct {
	ID   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name string    `gorm:"size:50;unique;not null" json:"name"`

	Users []User `gorm:"many2many:user_roles;joinForeignKey:RoleID;joinReferences:UserID" json:"-"`
}

type UserRole struct {
	ID     uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_user_role,composite:user_role" json:"user_id"`
	RoleID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_user_role,composite:user_role" json:"role_id"`
	User   User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Role   Role      `gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE" json:"-"`
}

type Biodata struct {
	ID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID         uuid.UUID `json:"user_id" gorm:"type:uuid;not null;unique"` 
	Name           string    `json:"name"`
	Age            int       `json:"age"`
	School         string    `json:"school"`
	ProfilePicture string    `json:"profile_picture"`
}
