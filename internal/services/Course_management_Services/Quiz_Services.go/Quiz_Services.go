package quizservicesgo

import (
	"course_online_backend/internal/dto"
	"course_online_backend/internal/models"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuizService struct {
	db *gorm.DB
}

func NewQuizService(db *gorm.DB) *QuizService {
	return &QuizService{db: db}
}

func (s *QuizService) CreateQuiz(req dto.CreateQuizRequest, userID uuid.UUID) (*dto.QuizResponse, error) {

	var course models.Course
	if err := s.db.First(&course, "id = ?", req.CourseID).Error; err != nil {
		return nil, errors.New("course not found")
	}

	quiz := &models.Quiz{
		CourseID:       req.CourseID,
		Name:           req.Name,
		Number:         req.Number,
		MinPassScore:   req.MinPassScore,
		IsPassRequired: req.IsPassRequired,
		ThumbnailURL:   req.ThumbnailURL,
		CreatedBy:      &userID,
	}

	if err := s.db.Create(quiz).Error; err != nil {
		return nil, err
	}

	return &dto.QuizResponse{
		ID:             quiz.ID,
		CourseID:       quiz.CourseID,
		Name:           quiz.Name,
		Number:         quiz.Number,
		MinPassScore:   quiz.MinPassScore,
		IsPassRequired: quiz.IsPassRequired,
		ThumbnailURL:   quiz.ThumbnailURL,
		CreatedBy:      quiz.CreatedBy,
		CreatedAt:      quiz.CreatedAt,
		UpdatedAt:      quiz.UpdatedAt,
		TotalQuestions: 0,
	}, nil
}

func (s *QuizService) GetQuizzesByCourse(courseID uuid.UUID) ([]dto.QuizResponse, error) {
	var quizzes []models.Quiz

	if err := s.db.Preload("Questions").Where("course_id = ?", courseID).Find(&quizzes).Error; err != nil {
		return nil, err
	}

	responses := make([]dto.QuizResponse, len(quizzes))
	for i, quiz := range quizzes {
		responses[i] = dto.QuizResponse{
			ID:             quiz.ID,
			CourseID:       quiz.CourseID,
			Name:           quiz.Name,
			Number:         quiz.Number,
			MinPassScore:   quiz.MinPassScore,
			IsPassRequired: quiz.IsPassRequired,
			ThumbnailURL:   quiz.ThumbnailURL,
			CreatedBy:      quiz.CreatedBy,
			CreatedAt:      quiz.CreatedAt,
			UpdatedAt:      quiz.UpdatedAt,
			TotalQuestions: len(quiz.Questions),
		}
	}

	return responses, nil
}

func (s *QuizService) GetQuizByID(quizID uuid.UUID) (*dto.QuizDetailResponse, error) {
	var quiz models.Quiz

	if err := s.db.Preload("Questions").First(&quiz, quizID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("quiz not found")
		}
		return nil, err
	}

	return &dto.QuizDetailResponse{
		ID:             quiz.ID,
		CourseID:       quiz.CourseID,
		Name:           quiz.Name,
		Number:         quiz.Number,
		MinPassScore:   quiz.MinPassScore,
		IsPassRequired: quiz.IsPassRequired,
		ThumbnailURL:   quiz.ThumbnailURL,
		CreatedBy:      quiz.CreatedBy,
		CreatedAt:      quiz.CreatedAt,
		UpdatedAt:      quiz.UpdatedAt,
		Questions:      len(quiz.Questions),
	}, nil
}
func (s *QuizService) UpdateQuiz(quizID uuid.UUID, req dto.UpdateQuizForm) (*dto.QuizResponse, error) {
	var quiz models.Quiz

	if err := s.db.First(&quiz, quizID).Error; err != nil {
		return nil, errors.New("quiz not found")
	}

	if req.Name != "" {
		quiz.Name = req.Name
	}
	if req.Number != 0 {
		quiz.Number = req.Number
	}

	quiz.MinPassScore = req.MinPassScore
	quiz.IsPassRequired = req.IsPassRequired

	if req.ThumbnailURL != "" {
		quiz.ThumbnailURL = req.ThumbnailURL
	}

	quiz.UpdatedAt = time.Now()

	if err := s.db.Save(&quiz).Error; err != nil {
		return nil, err
	}

	return &dto.QuizResponse{
		ID:             quiz.ID,
		CourseID:       quiz.CourseID,
		Name:           quiz.Name,
		Number:         quiz.Number,
		MinPassScore:   quiz.MinPassScore,
		IsPassRequired: quiz.IsPassRequired,
		ThumbnailURL:   quiz.ThumbnailURL,
		CreatedBy:      quiz.CreatedBy,
		CreatedAt:      quiz.CreatedAt,
		UpdatedAt:      quiz.UpdatedAt,
	}, nil
}

func (s *QuizService) SoftDeleteQuiz(quizID uuid.UUID) error {
	var quiz models.Quiz

	if err := s.db.First(&quiz, "id = ?", quizID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("quiz not found")
		}
		return err
	}

	return s.db.Delete(&quiz).Error
}

func (s *QuizService) PermanentDeleteQuiz(quizID uuid.UUID) error {
	var quiz models.Quiz

	if err := s.db.Unscoped().First(&quiz, quizID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("quiz not found")
		}
		return err
	}

	if err := s.db.Unscoped().Delete(&quiz).Error; err != nil {
		return err
	}

	return nil
}

func (s *QuizService) GetDeletedQuizzes(courseID uuid.UUID) ([]dto.QuizResponse, error) {
	var quizzes []models.Quiz

	if err := s.db.Unscoped().Preload("Questions").Where("course_id = ? AND deleted_at IS NOT NULL", courseID).Find(&quizzes).Error; err != nil {
		return nil, err
	}

	responses := make([]dto.QuizResponse, len(quizzes))
	for i, quiz := range quizzes {
		responses[i] = dto.QuizResponse{
			ID:             quiz.ID,
			CourseID:       quiz.CourseID,
			Name:           quiz.Name,
			Number:         quiz.Number,
			MinPassScore:   quiz.MinPassScore,
			IsPassRequired: quiz.IsPassRequired,
			ThumbnailURL:   quiz.ThumbnailURL,
			CreatedBy:      quiz.CreatedBy,
			CreatedAt:      quiz.CreatedAt,
			UpdatedAt:      quiz.UpdatedAt,
			DeletedAt: func() *time.Time {
				if quiz.DeletedAt.Valid {
					return &quiz.DeletedAt.Time
				}
				return nil
			}(),

			TotalQuestions: len(quiz.Questions),
		}
	}

	return responses, nil
}

func (s *QuizService) RestoreQuiz(quizID uuid.UUID) error {
	var quiz models.Quiz

	if err := s.db.Unscoped().First(&quiz, quizID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("quiz not found")
		}
		return err
	}

	if !quiz.DeletedAt.Valid {
		return errors.New("quiz is not deleted")
	}

	quiz.DeletedAt.Valid = false
	quiz.DeletedAt.Time = time.Time{}
	quiz.UpdatedAt = time.Now()

	return s.db.Unscoped().Save(&quiz).Error
}
