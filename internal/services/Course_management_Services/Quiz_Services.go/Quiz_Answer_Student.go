package quizservicesgo

import (
	"course_online_backend/internal/dto"
	"course_online_backend/internal/models"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StudentQuizService struct {
	db *gorm.DB
}

func NewStudentQuizService(db *gorm.DB) *StudentQuizService {
	return &StudentQuizService{db: db}
}

func (s *StudentQuizService) GetQuizzesInCourse(courseID uuid.UUID, userID uuid.UUID) ([]dto.QuizInCourseResponse, error) {
	var quizzes []models.Quiz

	if err := s.db.Preload("Questions").Where("course_id = ?", courseID).Order("number ASC").Find(&quizzes).Error; err != nil {
		return nil, err
	}

	responses := make([]dto.QuizInCourseResponse, len(quizzes))
	for i, quiz := range quizzes {
		totalPoints := 0
		for _, q := range quiz.Questions {
			totalPoints += q.Point
		}

		var bestScore *int
		var attemptsCount int64
		var isCompleted bool

		s.db.Model(&models.UserQuizAttempt{}).
			Where("user_id = ? AND quiz_id = ?", userID, quiz.ID).
			Count(&attemptsCount)

		var bestAttempt models.UserQuizAttempt
		if err := s.db.Where("user_id = ? AND quiz_id = ?", userID, quiz.ID).
			Order("score_achieved DESC").
			First(&bestAttempt).Error; err == nil {
			bestScore = &bestAttempt.ScoreAchieved
			isCompleted = bestAttempt.IsPassed
		}

		responses[i] = dto.QuizInCourseResponse{
			ID:              quiz.ID,
			CourseID:        quiz.CourseID,
			Name:            quiz.Name,
			Number:          quiz.Number,
			MinPassScore:    quiz.MinPassScore,
			IsPassRequired:  quiz.IsPassRequired,
			TotalQuestions:  len(quiz.Questions),
			TotalPoints:     totalPoints,
			MyBestScore:     bestScore,
			MyAttemptsCount: int(attemptsCount),
			IsCompleted:     isCompleted,
			CreatedAt:       quiz.CreatedAt,
			UpdatedAt:       quiz.UpdatedAt,
		}
	}

	return responses, nil
}

func (s *StudentQuizService) StartQuiz(quizID uuid.UUID, userID uuid.UUID) (*dto.StartQuizResponse, error) {
	var quiz models.Quiz

	if err := s.db.Preload("Questions", func(db *gorm.DB) *gorm.DB {
		return db.Order("number ASC")
	}).Preload("Questions.Answers").First(&quiz, quizID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("quiz not found")
		}
		return nil, err
	}

	totalPoints := 0
	questions := make([]dto.QuizQuestionForStudentResponse, len(quiz.Questions))

	for i, q := range quiz.Questions {
		totalPoints += q.Point

		answers := make([]dto.QuizAnswerForStudentResponse, len(q.Answers))
		for j, a := range q.Answers {
			answers[j] = dto.QuizAnswerForStudentResponse{
				ID:         a.ID,
				AnswerText: a.AnswerText,
			}
		}

		questions[i] = dto.QuizQuestionForStudentResponse{
			ID:            q.ID,
			QuestionTitle: q.QuestionTitle,
			Number:        q.Number,
			Point:         q.Point,
			Answers:       answers,
		}
	}

	return &dto.StartQuizResponse{
		QuizID:         quiz.ID,
		QuizName:       quiz.Name,
		MinPassScore:   quiz.MinPassScore,
		IsPassRequired: quiz.IsPassRequired,
		TotalQuestions: len(quiz.Questions),
		TotalPoints:    totalPoints,
		Questions:      questions,
	}, nil
}

func (s *StudentQuizService) SubmitQuiz(quizID uuid.UUID, userID uuid.UUID, req dto.SubmitQuizRequest) (*dto.SubmitQuizResponse, error) {
	var quiz models.Quiz

	if err := s.db.Preload("Questions.Answers").First(&quiz, quizID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("quiz not found")
		}
		return nil, err
	}

	if len(req.Answers) != len(quiz.Questions) {
		return nil, errors.New("you must answer all questions")
	}

	questionMap := make(map[uuid.UUID]models.QuizQuestion)
	for _, q := range quiz.Questions {
		questionMap[q.ID] = q
	}

	answerMap := make(map[uuid.UUID]models.QuizAnswer)
	for _, q := range quiz.Questions {
		for _, a := range q.Answers {
			answerMap[a.ID] = a
		}
	}

	answeredQuestions := make(map[uuid.UUID]bool)
	for _, ans := range req.Answers {
		if answeredQuestions[ans.QuestionID] {
			return nil, errors.New("duplicate answer for question")
		}
		answeredQuestions[ans.QuestionID] = true

		if _, exists := questionMap[ans.QuestionID]; !exists {
			return nil, errors.New("invalid question_id")
		}

		answer, exists := answerMap[ans.AnswerID]
		if !exists {
			return nil, errors.New("invalid answer_id")
		}

		if answer.QuestionID != ans.QuestionID {
			return nil, errors.New("answer does not belong to the question")
		}
	}

	correctAnswers := 0
	totalPoints := 0
	pointsEarned := 0
	answerDetails := []dto.AnswerResultDetail{}

	for _, studentAns := range req.Answers {
		question := questionMap[studentAns.QuestionID]
		studentAnswer := answerMap[studentAns.AnswerID]

		totalPoints += question.Point

		var correctAnswer models.QuizAnswer
		for _, a := range question.Answers {
			if a.IsCorrect {
				correctAnswer = a
				break
			}
		}

		isCorrect := studentAnswer.IsCorrect
		pointEarned := 0
		if isCorrect {
			correctAnswers++
			pointEarned = question.Point
			pointsEarned += question.Point
		}

		answerDetails = append(answerDetails, dto.AnswerResultDetail{
			QuestionID:        question.ID,
			QuestionTitle:     question.QuestionTitle,
			QuestionNumber:    question.Number,
			QuestionPoint:     question.Point,
			YourAnswerID:      studentAnswer.ID,
			YourAnswerText:    studentAnswer.AnswerText,
			IsCorrect:         isCorrect,
			CorrectAnswerID:   correctAnswer.ID,
			CorrectAnswerText: correctAnswer.AnswerText,
			PointEarned:       pointEarned,
		})
	}

	scoreAchieved := 0
	if len(quiz.Questions) > 0 {
		scoreAchieved = (correctAnswers * 100) / len(quiz.Questions)
	}

	isPassed := scoreAchieved >= quiz.MinPassScore

	tx := s.db.Begin()

	attempt := &models.UserQuizAttempt{
		UserID:          userID,
		QuizID:          quizID,
		AttemptDatetime: time.Now(),
		ScoreAchieved:   scoreAchieved,
		TotalQuestions:  len(quiz.Questions),
		CorrectAnswers:  correctAnswers,
		IsPassed:        isPassed,
	}

	if err := tx.Create(attempt).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, studentAns := range req.Answers {
		studentAnswer := answerMap[studentAns.AnswerID]

		userAnswer := &models.UserQuizAnswer{
			AttemptID:  attempt.ID,
			QuestionID: studentAns.QuestionID,
			AnswerID:   studentAns.AnswerID,
			IsCorrect:  studentAnswer.IsCorrect,
		}

		if err := tx.Create(userAnswer).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	tx.Commit()

	percentage := 0.0
	if totalPoints > 0 {
		percentage = (float64(pointsEarned) / float64(totalPoints)) * 100
	}

	return &dto.SubmitQuizResponse{
		AttemptID:       attempt.ID,
		QuizID:          quiz.ID,
		QuizName:        quiz.Name,
		AttemptDatetime: attempt.AttemptDatetime,
		ScoreAchieved:   scoreAchieved,
		TotalQuestions:  len(quiz.Questions),
		CorrectAnswers:  correctAnswers,
		IsPassed:        isPassed,
		MinPassScore:    quiz.MinPassScore,
		ResultSummary: dto.QuizResultSummary{
			TotalPoints:  totalPoints,
			PointsEarned: pointsEarned,
			Percentage:   percentage,
		},
		Answers: answerDetails,
	}, nil
}

func (s *StudentQuizService) GetQuizAttemptsHistory(quizID uuid.UUID, userID uuid.UUID) (*dto.QuizAttemptsHistoryResponse, error) {
	var quiz models.Quiz

	if err := s.db.First(&quiz, quizID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("quiz not found")
		}
		return nil, err
	}

	var attempts []models.UserQuizAttempt

	if err := s.db.Where("user_id = ? AND quiz_id = ?", userID, quizID).
		Order("attempt_datetime DESC").
		Find(&attempts).Error; err != nil {
		return nil, err
	}

	var bestScore *int
	var latestScore *int
	isCompleted := false

	if len(attempts) > 0 {
		latestScore = &attempts[0].ScoreAchieved

		best := 0
		for _, att := range attempts {
			if att.ScoreAchieved > best {
				best = att.ScoreAchieved
			}
			if att.IsPassed {
				isCompleted = true
			}
		}
		bestScore = &best
	}

	attemptSummaries := make([]dto.AttemptSummary, len(attempts))
	for i, att := range attempts {
		attemptSummaries[i] = dto.AttemptSummary{
			AttemptID:       att.ID,
			AttemptNumber:   len(attempts) - i,
			AttemptDatetime: att.AttemptDatetime,
			ScoreAchieved:   att.ScoreAchieved,
			CorrectAnswers:  att.CorrectAnswers,
			TotalQuestions:  att.TotalQuestions,
			IsPassed:        att.IsPassed,
		}
	}

	return &dto.QuizAttemptsHistoryResponse{
		QuizID:        quiz.ID,
		QuizName:      quiz.Name,
		MinPassScore:  quiz.MinPassScore,
		TotalAttempts: len(attempts),
		BestScore:     bestScore,
		LatestScore:   latestScore,
		IsCompleted:   isCompleted,
		Attempts:      attemptSummaries,
	}, nil
}

func (s *StudentQuizService) GetAttemptDetail(attemptID uuid.UUID, userID uuid.UUID) (*dto.AttemptDetailResponse, error) {
	var attempt models.UserQuizAttempt

	if err := s.db.Preload("Quiz").Preload("Quiz.Course").
		Preload("Answers").Preload("Answers.Question").
		Preload("Answers.Answer").First(&attempt, attemptID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("attempt not found")
		}
		return nil, err
	}

	if attempt.UserID != userID {
		return nil, errors.New("access denied: this attempt belongs to another user")
	}

	var quiz models.Quiz
	if err := s.db.Preload("Questions.Answers").First(&quiz, attempt.QuizID).Error; err != nil {
		return nil, err
	}

	questionMap := make(map[uuid.UUID]models.QuizQuestion)
	for _, q := range quiz.Questions {
		questionMap[q.ID] = q
	}

	totalPoints := 0
	pointsEarned := 0
	answerDetails := []dto.AnswerResultDetail{}

	for _, userAns := range attempt.Answers {
		question := questionMap[userAns.QuestionID]
		totalPoints += question.Point

		var correctAnswer models.QuizAnswer
		for _, a := range question.Answers {
			if a.IsCorrect {
				correctAnswer = a
				break
			}
		}

		pointEarned := 0
		if userAns.IsCorrect {
			pointEarned = question.Point
			pointsEarned += question.Point
		}

		answerDetails = append(answerDetails, dto.AnswerResultDetail{
			QuestionID:        question.ID,
			QuestionTitle:     question.QuestionTitle,
			QuestionNumber:    question.Number,
			QuestionPoint:     question.Point,
			YourAnswerID:      userAns.AnswerID,
			YourAnswerText:    userAns.Answer.AnswerText,
			IsCorrect:         userAns.IsCorrect,
			CorrectAnswerID:   correctAnswer.ID,
			CorrectAnswerText: correctAnswer.AnswerText,
			PointEarned:       pointEarned,
		})
	}

	percentage := 0.0
	if totalPoints > 0 {
		percentage = (float64(pointsEarned) / float64(totalPoints)) * 100
	}

	return &dto.AttemptDetailResponse{
		AttemptID:       attempt.ID,
		QuizID:          attempt.QuizID,
		QuizName:        attempt.Quiz.Name,
		CourseID:        attempt.Quiz.CourseID,
		CourseName:      attempt.Quiz.Course.Name,
		AttemptDatetime: attempt.AttemptDatetime,
		ScoreAchieved:   attempt.ScoreAchieved,
		TotalQuestions:  attempt.TotalQuestions,
		CorrectAnswers:  attempt.CorrectAnswers,
		IsPassed:        attempt.IsPassed,
		MinPassScore:    attempt.Quiz.MinPassScore,
		ResultSummary: dto.QuizResultSummary{
			TotalPoints:  totalPoints,
			PointsEarned: pointsEarned,
			Percentage:   percentage,
		},
		Answers: answerDetails,
	}, nil
}

func (s *StudentQuizService) GetAllMyAttempts(userID uuid.UUID, params dto.AllMyAttemptsQueryParams) (*dto.AllMyAttemptsResponse, error) {
	if params.Page == 0 {
		params.Page = 1
	}
	if params.Limit == 0 {
		params.Limit = 10
	}
	if params.SortBy == "" {
		params.SortBy = "attempt_datetime"
	}
	if params.SortOrder == "" {
		params.SortOrder = "desc"
	}

	query := s.db.Model(&models.UserQuizAttempt{}).Where("user_id = ?", userID)

	if params.CourseID != nil {
		query = query.Joins("JOIN quizzes ON quizzes.id = user_quiz_attempts.quiz_id").
			Where("quizzes.course_id = ?", *params.CourseID)
	}

	if params.IsPassed != nil {
		query = query.Where("is_passed = ?", *params.IsPassed)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	var attempts []models.UserQuizAttempt
	offset := (params.Page - 1) * params.Limit

	orderClause := params.SortBy + " " + params.SortOrder
	if err := query.Preload("Quiz").Preload("Quiz.Course").
		Order(orderClause).
		Limit(params.Limit).
		Offset(offset).
		Find(&attempts).Error; err != nil {
		return nil, err
	}

	var allAttempts []models.UserQuizAttempt
	s.db.Where("user_id = ?", userID).Find(&allAttempts)

	uniqueQuizzes := make(map[uuid.UUID]bool)
	passedCount := 0
	failedCount := 0
	totalScore := 0

	for _, att := range allAttempts {
		uniqueQuizzes[att.QuizID] = true
		if att.IsPassed {
			passedCount++
		} else {
			failedCount++
		}
		totalScore += att.ScoreAchieved
	}

	avgScore := 0.0
	if len(allAttempts) > 0 {
		avgScore = float64(totalScore) / float64(len(allAttempts))
	}

	statistics := dto.AttemptsStatistics{
		TotalQuizzesAttempted: len(uniqueQuizzes),
		TotalAttempts:         len(allAttempts),
		PassedAttempts:        passedCount,
		FailedAttempts:        failedCount,
		AverageScore:          avgScore,
	}

	attemptItems := make([]dto.MyAttemptItem, len(attempts))
	for i, att := range attempts {
		attemptItems[i] = dto.MyAttemptItem{
			AttemptID:       att.ID,
			QuizID:          att.QuizID,
			QuizName:        att.Quiz.Name,
			QuizNumber:      att.Quiz.Number,
			CourseID:        att.Quiz.CourseID,
			CourseName:      att.Quiz.Course.Name,
			AttemptDatetime: att.AttemptDatetime,
			ScoreAchieved:   att.ScoreAchieved,
			TotalQuestions:  att.TotalQuestions,
			CorrectAnswers:  att.CorrectAnswers,
			IsPassed:        att.IsPassed,
			MinPassScore:    att.Quiz.MinPassScore,
		}
	}

	totalPages := int(total) / params.Limit
	if int(total)%params.Limit != 0 {
		totalPages++
	}

	return &dto.AllMyAttemptsResponse{
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
		Statistics: statistics,
		Attempts:   attemptItems,
	}, nil
}