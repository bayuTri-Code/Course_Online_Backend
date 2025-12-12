package dto

import (
	"time"
	"github.com/google/uuid"
)

type QuizInCourseResponse struct {
	ID              uuid.UUID  `json:"id"`
	CourseID        uuid.UUID  `json:"course_id"`
	Name            string     `json:"name"`
	Number          int        `json:"number"`
	MinPassScore    int        `json:"min_pass_score"`
	IsPassRequired  bool       `json:"is_pass_required"`
	TotalQuestions  int        `json:"total_questions"`
	TotalPoints     int        `json:"total_points"`
	MyBestScore     *int       `json:"my_best_score"`
	MyAttemptsCount int        `json:"my_attempts_count"`
	IsCompleted     bool       `json:"is_completed"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type StartQuizResponse struct {
	QuizID         uuid.UUID                    `json:"quiz_id"`
	QuizName       string                       `json:"quiz_name"`
	Description    string                       `json:"description,omitempty"`
	MinPassScore   int                          `json:"min_pass_score"`
	IsPassRequired bool                         `json:"is_pass_required"`
	TotalQuestions int                          `json:"total_questions"`
	TotalPoints    int                          `json:"total_points"`
	Questions      []QuizQuestionForStudentResponse `json:"questions"`
}

type QuizQuestionForStudentResponse struct {
	ID            uuid.UUID                      `json:"id"`
	QuestionTitle string                         `json:"question_title"`
	Number        int                            `json:"number"`
	Point         int                            `json:"point"`
	Answers       []QuizAnswerForStudentResponse `json:"answers"`
}

type QuizAnswerForStudentResponse struct {
	ID         uuid.UUID `json:"id"`
	AnswerText string    `json:"answer_text"`
}

type SubmitQuizRequest struct {
	Answers []StudentAnswerInput `json:"answers" binding:"required,min=1,dive"`
}

type StudentAnswerInput struct {
	QuestionID uuid.UUID `json:"question_id" binding:"required"`
	AnswerID   uuid.UUID `json:"answer_id" binding:"required"`
}

type SubmitQuizResponse struct {
	AttemptID       uuid.UUID              `json:"attempt_id"`
	QuizID          uuid.UUID              `json:"quiz_id"`
	QuizName        string                 `json:"quiz_name"`
	AttemptDatetime time.Time              `json:"attempt_datetime"`
	ScoreAchieved   int                    `json:"score_achieved"`
	TotalQuestions  int                    `json:"total_questions"`
	CorrectAnswers  int                    `json:"correct_answers"`
	IsPassed        bool                   `json:"is_passed"`
	MinPassScore    int                    `json:"min_pass_score"`
	ResultSummary   QuizResultSummary      `json:"result_summary"`
	Answers         []AnswerResultDetail   `json:"answers"`
}

type QuizResultSummary struct {
	TotalPoints   int     `json:"total_points"`
	PointsEarned  int     `json:"points_earned"`
	Percentage    float64 `json:"percentage"`
}

type AnswerResultDetail struct {
	QuestionID         uuid.UUID `json:"question_id"`
	QuestionTitle      string    `json:"question_title"`
	QuestionNumber     int       `json:"question_number"`
	QuestionPoint      int       `json:"question_point"`
	YourAnswerID       uuid.UUID `json:"your_answer_id"`
	YourAnswerText     string    `json:"your_answer_text"`
	IsCorrect          bool      `json:"is_correct"`
	CorrectAnswerID    uuid.UUID `json:"correct_answer_id"`
	CorrectAnswerText  string    `json:"correct_answer_text"`
	PointEarned        int       `json:"point_earned"`
}

type QuizAttemptsHistoryResponse struct {
	QuizID         uuid.UUID              `json:"quiz_id"`
	QuizName       string                 `json:"quiz_name"`
	MinPassScore   int                    `json:"min_pass_score"`
	TotalAttempts  int                    `json:"total_attempts"`
	BestScore      *int                   `json:"best_score"`
	LatestScore    *int                   `json:"latest_score"`
	IsCompleted    bool                   `json:"is_completed"`
	Attempts       []AttemptSummary       `json:"attempts"`
}

type AttemptSummary struct {
	AttemptID       uuid.UUID `json:"attempt_id"`
	AttemptNumber   int       `json:"attempt_number"`
	AttemptDatetime time.Time `json:"attempt_datetime"`
	ScoreAchieved   int       `json:"score_achieved"`
	CorrectAnswers  int       `json:"correct_answers"`
	TotalQuestions  int       `json:"total_questions"`
	IsPassed        bool      `json:"is_passed"`
}

type AttemptDetailResponse struct {
	AttemptID       uuid.UUID            `json:"attempt_id"`
	QuizID          uuid.UUID            `json:"quiz_id"`
	QuizName        string               `json:"quiz_name"`
	CourseID        uuid.UUID            `json:"course_id"`
	CourseName      string               `json:"course_name"`
	AttemptDatetime time.Time            `json:"attempt_datetime"`
	ScoreAchieved   int                  `json:"score_achieved"`
	TotalQuestions  int                  `json:"total_questions"`
	CorrectAnswers  int                  `json:"correct_answers"`
	IsPassed        bool                 `json:"is_passed"`
	MinPassScore    int                  `json:"min_pass_score"`
	ResultSummary   QuizResultSummary    `json:"result_summary"`
	Answers         []AnswerResultDetail `json:"answers"`
}

type AllMyAttemptsResponse struct {
	Total      int64                     `json:"total"`
	Page       int                       `json:"page"`
	Limit      int                       `json:"limit"`
	TotalPages int                       `json:"total_pages"`
	Statistics AttemptsStatistics        `json:"statistics"`
	Attempts   []MyAttemptItem           `json:"attempts"`
}

type AttemptsStatistics struct {
	TotalQuizzesAttempted int     `json:"total_quizzes_attempted"`
	TotalAttempts         int     `json:"total_attempts"`
	PassedAttempts        int     `json:"passed_attempts"`
	FailedAttempts        int     `json:"failed_attempts"`
	AverageScore          float64 `json:"average_score"`
}

type MyAttemptItem struct {
	AttemptID       uuid.UUID `json:"attempt_id"`
	QuizID          uuid.UUID `json:"quiz_id"`
	QuizName        string    `json:"quiz_name"`
	QuizNumber      int       `json:"quiz_number"`
	CourseID        uuid.UUID `json:"course_id"`
	CourseName      string    `json:"course_name"`
	AttemptDatetime time.Time `json:"attempt_datetime"`
	ScoreAchieved   int       `json:"score_achieved"`
	TotalQuestions  int       `json:"total_questions"`
	CorrectAnswers  int       `json:"correct_answers"`
	IsPassed        bool      `json:"is_passed"`
	MinPassScore    int       `json:"min_pass_score"`
}

type AllMyAttemptsQueryParams struct {
	Page      int        `form:"page" binding:"omitempty,min=1"`
	Limit     int        `form:"limit" binding:"omitempty,min=1,max=100"`
	CourseID  *uuid.UUID `form:"course_id"`
	IsPassed  *bool      `form:"is_passed"`
	SortBy    string     `form:"sort_by" binding:"omitempty,oneof=attempt_datetime score_achieved"`
	SortOrder string     `form:"sort_order" binding:"omitempty,oneof=asc desc"`
}