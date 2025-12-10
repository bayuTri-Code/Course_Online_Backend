package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Quiz struct {
	ID             uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CourseID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"course_id"`
	Name           string         `gorm:"size:200;not null" json:"name"`
	Number         int            `gorm:"not null" json:"number"`
	MinPassScore   int            `gorm:"not null" json:"min_pass_score"`
	IsPassRequired bool           `gorm:"default:false" json:"is_pass_required"`
	CreatedBy      *uuid.UUID     `gorm:"type:uuid;index" json:"created_by"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty" swaggertype:"string" swaggerignore:"true"`
	CreatedAt      time.Time      `gorm:"type:timestamp;default:current_timestamp" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"type:timestamp;default:current_timestamp" json:"updated_at"`

	Course    Course            `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE" json:"course"`
	Creator   *User             `gorm:"foreignKey:CreatedBy;constraint:OnDelete:SET NULL" json:"creator"`
	Questions []QuizQuestion    `gorm:"foreignKey:QuizID;constraint:OnDelete:CASCADE" json:"questions"`
	Attempts  []UserQuizAttempt `gorm:"foreignKey:QuizID;constraint:OnDelete:CASCADE" json:"attempts"`
}

type QuizQuestion struct {
	ID            uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	QuizID        uuid.UUID      `gorm:"type:uuid;not null;index" json:"quiz_id"`
	QuestionTitle string         `gorm:"type:text;not null" json:"question_title"`
	Number        int            `gorm:"not null" json:"number"`
	Point         int            `gorm:"not null" json:"point"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty" swaggertype:"string" swaggerignore:"true" format:"date-time"`

	Quiz    Quiz         `gorm:"foreignKey:QuizID;constraint:OnDelete:CASCADE" json:"quiz"`
	Answers []QuizAnswer `gorm:"foreignKey:QuestionID;constraint:OnDelete:CASCADE" json:"answers"`
}

type QuizAnswer struct {
	ID         uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	QuestionID uuid.UUID      `gorm:"type:uuid;not null;index" json:"question_id"`
	AnswerText string         `gorm:"type:text;not null" json:"answer_text"`
	IsCorrect  bool           `gorm:"default:false" json:"is_correct"`
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty" swaggertype:"string" swaggerignore:"true" format:"date-time"`

	Question QuizQuestion `gorm:"foreignKey:QuestionID;constraint:OnDelete:CASCADE" json:"question"`
}

type UserQuizAttempt struct {
	ID              uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID          uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	QuizID          uuid.UUID      `gorm:"type:uuid;not null;index" json:"quiz_id"`
	AttemptDatetime time.Time      `gorm:"not null" json:"attempt_datetime"`
	ScoreAchieved   int            `gorm:"not null" json:"score_achieved"`
	TotalQuestions  int            `gorm:"not null" json:"total_questions"`
	CorrectAnswers  int            `gorm:"not null" json:"correct_answers"`
	IsPassed        bool           `gorm:"default:false" json:"is_passed"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty" swaggertype:"string" swaggerignore:"true" format:"date-time"`

	User    User             `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user"`
	Quiz    Quiz             `gorm:"foreignKey:QuizID;constraint:OnDelete:CASCADE" json:"quiz"`
	Answers []UserQuizAnswer `gorm:"foreignKey:AttemptID;constraint:OnDelete:CASCADE" json:"answers"`
}

type UserQuizAnswer struct {
	ID         uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	AttemptID  uuid.UUID      `gorm:"type:uuid;not null;index" json:"attempt_id"`
	QuestionID uuid.UUID      `gorm:"type:uuid;not null;index" json:"question_id"`
	AnswerID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"answer_id"`
	IsCorrect  bool           `gorm:"default:false" json:"is_correct"`
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty" swaggertype:"string" swaggerignore:"true" format:"date-time"`

	Attempt  UserQuizAttempt `gorm:"foreignKey:AttemptID;constraint:OnDelete:CASCADE" json:"attempt"`
	Question QuizQuestion    `gorm:"foreignKey:QuestionID;constraint:OnDelete:CASCADE" json:"question"`
	Answer   QuizAnswer      `gorm:"foreignKey:AnswerID;constraint:OnDelete:CASCADE" json:"answer"`
}
