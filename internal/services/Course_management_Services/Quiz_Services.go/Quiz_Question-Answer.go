package quizservicesgo

import (
	"course_online_backend/internal/dto"
	"course_online_backend/internal/models"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuestionService struct {
	db *gorm.DB
}

func NewQuestionService(db *gorm.DB) *QuestionService {
	return &QuestionService{db: db}
}
func (s *QuestionService) CreateQuestion(quizID uuid.UUID, req dto.CreateQuestionRequest) (*dto.QuestionResponse, error) {
	var quiz models.Quiz
	if err := s.db.First(&quiz, quizID).Error; err != nil {
		return nil, errors.New("quiz not found")
	}

	correctCount := 0
	for _, ans := range req.Answers {
		if ans.IsCorrect {
			correctCount++
		}
	}

	if correctCount != 1 {
		return nil, errors.New("exactly one answer must be marked as correct")
	}

	tx := s.db.Begin()

	question := &models.QuizQuestion{
		QuizID:        quizID,
		QuestionTitle: req.QuestionTitle,
		Number:        req.Number,
		Point:         req.Point,
	}

	if err := tx.Create(question).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, ans := range req.Answers {
		answer := models.QuizAnswer{
			QuestionID: question.ID,
			AnswerText: ans.AnswerText,
			IsCorrect:  ans.IsCorrect,
		}

		if err := tx.Create(&answer).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	tx.Commit()

	s.db.Preload("Answers").First(question, question.ID)

	return s.mapToQuestionResponse(question), nil
}


func (s *QuestionService) BulkCreateQuestions(quizID uuid.UUID, req dto.BulkCreateQuestionsRequest) ([]dto.QuestionResponse, error) {
	var quiz models.Quiz
	if err := s.db.First(&quiz, quizID).Error; err != nil {
		return nil, errors.New("quiz not found")
	}

	tx := s.db.Begin()

	responses := []dto.QuestionResponse{}

	for _, qReq := range req.Questions {
		correctCount := 0
		for _, ans := range qReq.Answers {
			if ans.IsCorrect {
				correctCount++
			}
		}

		if correctCount != 1 {
			tx.Rollback()
			return nil, errors.New("exactly one answer must be marked as correct for each question")
		}

		question := &models.QuizQuestion{
			QuizID:        quizID,
			QuestionTitle: qReq.QuestionTitle,
			Number:        qReq.Number,
			Point:         qReq.Point,
		}

		if err := tx.Create(question).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		answers := []models.QuizAnswer{}
		for _, ans := range qReq.Answers {
			answer := models.QuizAnswer{
				QuestionID: question.ID,
				AnswerText: ans.AnswerText,
				IsCorrect:  ans.IsCorrect,
			}

			if err := tx.Create(&answer).Error; err != nil {
				tx.Rollback()
				return nil, err
			}

			answers = append(answers, answer)
		}

		question.Answers = answers
		responses = append(responses, *s.mapToQuestionResponse(question))
	}

	tx.Commit()

	return responses, nil
}

func (s *QuestionService) GetQuestionsByQuiz(quizID uuid.UUID) ([]dto.QuestionResponse, error) {
	var questions []models.QuizQuestion

	if err := s.db.Preload("Answers").Where("quiz_id = ?", quizID).Order("number ASC").Find(&questions).Error; err != nil {
		return nil, err
	}

	responses := make([]dto.QuestionResponse, len(questions))
	for i, question := range questions {
		responses[i] = *s.mapToQuestionResponse(&question)
	}

	return responses, nil
}

func (s *QuestionService) GetQuestionByID(questionID uuid.UUID) (*dto.QuestionDetailResponse, error) {
	var question models.QuizQuestion

	if err := s.db.Preload("Answers").First(&question, questionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("question not found")
		}
		return nil, err
	}

	return s.mapToQuestionDetailResponse(&question), nil
}

func (s *QuestionService) UpdateQuestion(questionID uuid.UUID, req dto.UpdateQuestionRequest) (*dto.QuestionResponse, error) {
	var question models.QuizQuestion

	if err := s.db.First(&question, questionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("question not found")
		}
		return nil, err
	}

	correctCount := 0
	for _, ans := range req.Answers {
		if ans.IsCorrect {
			correctCount++
		}
	}

	if correctCount != 1 {
		return nil, errors.New("exactly one answer must be marked as correct")
	}

	tx := s.db.Begin()

	question.QuestionTitle = req.QuestionTitle
	question.Number = req.Number
	question.Point = req.Point
	question.UpdatedAt = time.Now()

	if err := tx.Save(&question).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Where("question_id = ?", questionID).Delete(&models.QuizAnswer{}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, ans := range req.Answers {
		answer := models.QuizAnswer{
			QuestionID: questionID,
			AnswerText: ans.AnswerText,
			IsCorrect:  ans.IsCorrect,
		}

		if err := tx.Create(&answer).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

	}

	tx.Commit()

	s.db.Preload("Answers").First(&question, questionID)

	return s.mapToQuestionResponse(&question), nil
}

func (s *QuestionService) SoftDeleteQuestion(questionID uuid.UUID) error {
	var question models.QuizQuestion

	if err := s.db.First(&question, questionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("question not found")
		}
		return err
	}

	if err := s.db.Delete(&question).Error; err != nil {
		return err
	}

	return nil
}

func (s *QuestionService) PermanentDeleteQuestion(questionID uuid.UUID) error {
	var question models.QuizQuestion

	if err := s.db.Unscoped().First(&question, questionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("question not found")
		}
		return err
	}

	if err := s.db.Unscoped().Delete(&question).Error; err != nil {
		return err
	}

	return nil
}

func (s *QuestionService) GetDeletedQuestions(quizID uuid.UUID) ([]dto.QuestionResponse, error) {
	var questions []models.QuizQuestion

	if err := s.db.Unscoped().Preload("Answers").Where("quiz_id = ? AND deleted_at IS NOT NULL", quizID).Order("number ASC").Find(&questions).Error; err != nil {
		return nil, err
	}

	responses := make([]dto.QuestionResponse, len(questions))
	for i, question := range questions {
		responses[i] = *s.mapToQuestionResponse(&question)
	}

	return responses, nil
}

func (s *QuestionService) RestoreQuestion(questionID uuid.UUID) error {
	var question models.QuizQuestion

	if err := s.db.Unscoped().First(&question, questionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("question not found")
		}
		return err
	}

	if question.DeletedAt.Time.IsZero() {
		return errors.New("question is not deleted")
	}

	question.UpdatedAt = time.Now()

	if err := s.db.Unscoped().Model(&question).Update("deleted_at", nil).Error; err != nil {
		return err
	}

	return nil
}

func (s *QuestionService) mapToQuestionResponse(question *models.QuizQuestion) *dto.QuestionResponse {
	answers := make([]dto.AnswerResponse, len(question.Answers))
	for i, ans := range question.Answers {
		var deletedAt *time.Time
		if !ans.DeletedAt.Time.IsZero() {
			deletedAt = &ans.DeletedAt.Time
		}

		answers[i] = dto.AnswerResponse{
			ID:         ans.ID,
			QuestionID: ans.QuestionID,
			AnswerText: ans.AnswerText,
			IsCorrect:  ans.IsCorrect,
			CreatedAt:  ans.CreatedAt,
			UpdatedAt:  ans.UpdatedAt,
			DeletedAt:  deletedAt,
		}
	}

	var deletedAt *time.Time
	if !question.DeletedAt.Time.IsZero() {
		deletedAt = &question.DeletedAt.Time
	}

	return &dto.QuestionResponse{
		ID:            question.ID,
		QuizID:        question.QuizID,
		QuestionTitle: question.QuestionTitle,
		Number:        question.Number,
		Point:         question.Point,
		CreatedAt:     question.CreatedAt,
		UpdatedAt:     question.UpdatedAt,
		DeletedAt:     deletedAt,
		Answers:       answers,
	}
}

func (s *QuestionService) mapToQuestionDetailResponse(question *models.QuizQuestion) *dto.QuestionDetailResponse {
	answers := make([]dto.AnswerResponse, len(question.Answers))
	for i, ans := range question.Answers {
		var deletedAt *time.Time
		if !ans.DeletedAt.Time.IsZero() {
			deletedAt = &ans.DeletedAt.Time
		}

		answers[i] = dto.AnswerResponse{
			ID:         ans.ID,
			QuestionID: ans.QuestionID,
			AnswerText: ans.AnswerText,
			IsCorrect:  ans.IsCorrect,
			CreatedAt:  ans.CreatedAt,
			UpdatedAt:  ans.UpdatedAt,
			DeletedAt:  deletedAt,
		}
	}

	return &dto.QuestionDetailResponse{
		ID:            question.ID,
		QuizID:        question.QuizID,
		QuestionTitle: question.QuestionTitle,
		Number:        question.Number,
		Point:         question.Point,
		CreatedAt:     question.CreatedAt,
		UpdatedAt:     question.UpdatedAt,
		Answers:       answers,
	}
}
