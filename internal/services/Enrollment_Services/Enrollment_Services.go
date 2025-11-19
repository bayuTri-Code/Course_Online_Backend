package enrollmentservices

import (
	repository "course_online_backend/database/Repository"
	"course_online_backend/internal/models"
	constants "course_online_backend/internal/models/Constants"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EnrollmentService interface {
	EnrollFreeCourse(userID, courseID uuid.UUID) (*models.Enrollment, error)
	CheckEnrollment(userID, courseID uuid.UUID) (bool, error)
	GetMyEnrollments(userID uuid.UUID) ([]models.Enrollment, error)
	GetEnrollmentDetail(userID, enrollmentID uuid.UUID) (*models.Enrollment, error)
	UnenrollCourse(userID, enrollmentID uuid.UUID) error
	UpdateEnrollmentStatus(userID, enrollmentID uuid.UUID, status string) error
	GetCourseStudents(courseID uuid.UUID, userRoles []string) ([]models.Enrollment, error)
}

type enrollmentService struct {
	enrollmentRepo repository.EnrollmentRepository
	db             *gorm.DB
}

func NewEnrollmentService(enrollmentRepo repository.EnrollmentRepository, db *gorm.DB) EnrollmentService {
	return &enrollmentService{
		enrollmentRepo: enrollmentRepo,
		db:             db,
	}
}

func (s *enrollmentService) EnrollFreeCourse(userID, courseID uuid.UUID) (*models.Enrollment, error) {
	if userID == uuid.Nil || courseID == uuid.Nil {
		return nil, errors.New("invalid user ID or course ID")
	}

	var course models.Course
	if err := s.db.First(&course, "id = ?", courseID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("course not found")
		}
		return nil, err
	}

	if course.Status != constants.CourseStatusPublished {
		return nil, errors.New("course is not available for enrollment")
	}

	if course.Price > 0 {
		return nil, errors.New("this course is not free, please use payment endpoint")
	}

	if course.EnrollmentDeadline != nil && time.Now().After(*course.EnrollmentDeadline) {
		return nil, errors.New("enrollment deadline has passed")
	}

	if course.MaxStudents > 0 {
		count, err := s.enrollmentRepo.CountByCourse(courseID)
		if err != nil {
			return nil, err
		}
		if count >= int64(course.MaxStudents) {
			return nil, errors.New("course is full, max capacity reached")
		}
	}

	exists, err := s.enrollmentRepo.CheckExists(userID, courseID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("you are already enrolled in this course")
	}

	now := time.Now()
	enrollment := &models.Enrollment{
		ID:                 uuid.New(),
		UserID:             userID,
		CourseID:           courseID,
		EnrollmentDatetime: now,
		Status:             "active",
		StatusPayment:      "free",
		Progress:           0,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	if course.Duration > 0 {
		expiredDate := now.AddDate(0, 0, course.Duration)
		enrollment.ExpiredDate = &expiredDate
	}

	if err := s.enrollmentRepo.Create(enrollment); err != nil {
		return nil, err
	}

	result, err := s.enrollmentRepo.GetByID(enrollment.ID)
	if err != nil {
		return enrollment, nil
	}

	return result, nil
}

func (s *enrollmentService) CheckEnrollment(userID, courseID uuid.UUID) (bool, error) {
	if userID == uuid.Nil || courseID == uuid.Nil {
		return false, errors.New("invalid user ID or course ID")
	}
	return s.enrollmentRepo.CheckExists(userID, courseID)
}

func (s *enrollmentService) GetMyEnrollments(userID uuid.UUID) ([]models.Enrollment, error) {
	if userID == uuid.Nil {
		return nil, errors.New("invalid user ID")
	}
	return s.enrollmentRepo.GetUserEnrollments(userID)
}

func (s *enrollmentService) GetEnrollmentDetail(userID, enrollmentID uuid.UUID) (*models.Enrollment, error) {
	if userID == uuid.Nil || enrollmentID == uuid.Nil {
		return nil, errors.New("invalid user ID or enrollment ID")
	}

	enrollment, err := s.enrollmentRepo.GetByID(enrollmentID)
	if err != nil {
		return nil, err
	}

	if enrollment.UserID != userID {
		return nil, errors.New("unauthorized to access this enrollment")
	}

	return enrollment, nil
}

func (s *enrollmentService) UnenrollCourse(userID, enrollmentID uuid.UUID) error {
	if userID == uuid.Nil || enrollmentID == uuid.Nil {
		return errors.New("invalid user ID or enrollment ID")
	}

	enrollment, err := s.enrollmentRepo.GetByID(enrollmentID)
	if err != nil {
		return err
	}

	if enrollment.UserID != userID {
		return errors.New("unauthorized to unenroll this course")
	}

	if enrollment.Status == "completed" {
		return errors.New("cannot unenroll from completed course")
	}

	return s.enrollmentRepo.Delete(enrollmentID)
}

func (s *enrollmentService) UpdateEnrollmentStatus(userID, enrollmentID uuid.UUID, status string) error {
	if userID == uuid.Nil || enrollmentID == uuid.Nil {
		return errors.New("invalid user ID or enrollment ID")
	}

	validStatuses := map[string]bool{
		"active":    true,
		"dropped":   true,
		"completed": true,
	}
	if !validStatuses[status] {
		return errors.New("invalid status, must be: active, dropped, or completed")
	}

	enrollment, err := s.enrollmentRepo.GetByID(enrollmentID)
	if err != nil {
		return err
	}

	if enrollment.UserID != userID {
		return errors.New("unauthorized to update this enrollment")
	}

	return s.enrollmentRepo.UpdateStatus(enrollmentID, status)
}

func (s *enrollmentService) GetCourseStudents(courseID uuid.UUID, userRoles []string) ([]models.Enrollment, error) {
	if courseID == uuid.Nil {
		return nil, errors.New("invalid course ID")
	}

	hasAccess := false
	allowedRoles := []string{"instructor", "admin", "super_admin"}
	for _, role := range userRoles {
		for _, allowed := range allowedRoles {
			if role == allowed {
				hasAccess = true
				break
			}
		}
	}

	if !hasAccess {
		return nil, errors.New("only instructors and admins can view course students")
	}

	return s.enrollmentRepo.GetCourseEnrollments(courseID)
}


