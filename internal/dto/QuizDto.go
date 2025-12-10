package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateQuizRequest struct {
	CourseID       uuid.UUID `json:"course_id" binding:"required"`
	Name           string    `json:"name" binding:"required,max=200"`
	Number         int       `json:"number" binding:"required,min=1"`
	MinPassScore   int       `json:"min_pass_score" binding:"required,min=0,max=100"`
	IsPassRequired bool      `json:"is_pass_required"`
}

type UpdateQuizRequest struct {
	Name           string `json:"name" binding:"required,max=200"`
	Number         int    `json:"number" binding:"required,min=1"`
	MinPassScore   int    `json:"min_pass_score" binding:"required,min=0,max=100"`
	IsPassRequired bool   `json:"is_pass_required"`
}

type QuizResponse struct {
	ID             uuid.UUID  `json:"id"`
	CourseID       uuid.UUID  `json:"course_id"`
	Name           string     `json:"name"`
	Number         int        `json:"number"`
	MinPassScore   int        `json:"min_pass_score"`
	IsPassRequired bool       `json:"is_pass_required"`
	CreatedBy      *uuid.UUID `json:"created_by"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `gorm:"index" json:"deleted_at,omitempty"`
	TotalQuestions int        `json:"total_questions"`
}

type QuizDetailResponse struct {
	ID             uuid.UUID  `json:"id"`
	CourseID       uuid.UUID  `json:"course_id"`
	Name           string     `json:"name"`
	Number         int        `json:"number"`
	MinPassScore   int        `json:"min_pass_score"`
	IsPassRequired bool       `json:"is_pass_required"`
	CreatedBy      *uuid.UUID `json:"created_by"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	Questions      int        `json:"questions"`
}


//quiz question and answer DTOs

type CreateQuestionRequest struct {
	QuestionTitle string               `json:"question_title" binding:"required"`
	Number        int                  `json:"number" binding:"required,min=1"`
	Point         int                  `json:"point" binding:"required,min=1"`
	Answers       []CreateAnswerRequest `json:"answers" binding:"required,min=2,dive"`
}

type CreateAnswerRequest struct {
	AnswerText string `json:"answer_text" binding:"required"`
	IsCorrect  bool   `json:"is_correct"`
}

type BulkCreateQuestionsRequest struct {
	Questions []CreateQuestionRequest `json:"questions" binding:"required,min=1,dive"`
}

type UpdateQuestionRequest struct {
	QuestionTitle string               `json:"question_title" binding:"required"`
	Number        int                  `json:"number" binding:"required,min=1"`
	Point         int                  `json:"point" binding:"required,min=1"`
	Answers       []UpdateAnswerRequest `json:"answers" binding:"required,min=2,dive"`
}

type UpdateAnswerRequest struct {
	ID         *uuid.UUID `json:"id"`
	AnswerText string     `json:"answer_text" binding:"required"`
	IsCorrect  bool       `json:"is_correct"`
}

type QuestionResponse struct {
	ID            uuid.UUID        `json:"id"`
	QuizID        uuid.UUID        `json:"quiz_id"`
	QuestionTitle string           `json:"question_title"`
	Number        int              `json:"number"`
	Point         int              `json:"point"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
	DeletedAt     *time.Time       `json:"deleted_at,omitempty"`
	Answers       []AnswerResponse `json:"answers"`
}

type AnswerResponse struct {
	ID         uuid.UUID  `json:"id"`
	QuestionID uuid.UUID  `json:"question_id"`
	AnswerText string     `json:"answer_text"`
	IsCorrect  bool       `json:"is_correct"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
}

type QuestionDetailResponse struct {
	ID            uuid.UUID        `json:"id"`
	QuizID        uuid.UUID        `json:"quiz_id"`
	QuestionTitle string           `json:"question_title"`
	Number        int              `json:"number"`
	Point         int              `json:"point"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
	Answers       []AnswerResponse `json:"answers"`
}