package enrollmentservices

import (
	"course_online_backend/internal/models"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EnrollmentService struct {
	db *gorm.DB
}

func NewEnrollmentService(db *gorm.DB) *EnrollmentService {
	return &EnrollmentService{db: db}
}

func (s *EnrollmentService) EnrollCourse(userID uuid.UUID, courseID uuid.UUID) (*models.Enrollment, error) {
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

	if !course.IsPublished() {
		return nil, errors.New("course is not available for enrollment")
	}

	if !course.IsFree() {
		return nil, errors.New("this course is not free, please use payment endpoint")
	}

	if course.EnrollmentDeadline != nil && time.Now().After(*course.EnrollmentDeadline) {
		return nil, errors.New("enrollment deadline has passed")
	}

	if course.MaxStudents > 0 {
		var count int64
		if err := s.db.Model(&models.Enrollment{}).
			Where("course_id = ? AND status = ?", courseID, "active").
			Count(&count).Error; err != nil {
			return nil, err
		}
		if count >= int64(course.MaxStudents) {
			return nil, errors.New("course is full, max capacity reached")
		}
	}

	var existCount int64
	if err := s.db.Model(&models.Enrollment{}).
		Where("user_id = ? AND course_id = ?", userID, courseID).
		Count(&existCount).Error; err != nil {
		return nil, err
	}
	if existCount > 0 {
		return nil, errors.New("you are already enrolled in this course")
	}

	now := time.Now()
	enroll := &models.Enrollment{
		ID:                 uuid.New(),
		UserID:             userID,
		CourseID:           courseID,
		EnrollmentDatetime: now,
		Status:             "active",
		StatusPayment:      "free",
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	if course.Duration > 0 {
		exp := now.AddDate(0, 0, course.Duration)
		enroll.ExpiredDate = &exp
	}

	if err := s.db.Create(enroll).Error; err != nil {
		return nil, err
	}

	var result models.Enrollment
	if err := s.db.
		Preload("Course").
		Preload("Course.CourseType").
		Preload("Course.Creator").
		Preload("Course.Instructor").
		Preload("User").
		Preload("User.Biodata").
		Preload("User.Roles").
		First(&result, "id = ?", enroll.ID).Error; err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *EnrollmentService) CheckEnrollment(userID uuid.UUID, courseID uuid.UUID) (bool, error) {
	var count int64
	if err := s.db.Model(&models.Enrollment{}).
		Where("user_id = ? AND course_id = ?", userID, courseID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *EnrollmentService) GetMyEnrollments(userID uuid.UUID) ([]models.Enrollment, error) {
	var enrollments []models.Enrollment
	if err := s.db.
		Preload("Course").
		Preload("Course.CourseType").
		Preload("Course.Instructor").
		Where("user_id = ?", userID).
		Order("enrollment_datetime DESC").
		Find(&enrollments).Error; err != nil {
		return nil, err
	}

	return enrollments, nil
}

func (s *EnrollmentService) GetEnrollmentDetail(
	userID uuid.UUID,
	enrollmentID uuid.UUID,
	roles []string,
) (*models.Enrollment, error) {

	var e models.Enrollment
	if err := s.db.
		Preload("Course").
		Preload("Course.CourseType").
		Preload("Course.Creator").
		Preload("Course.Instructor").
		Preload("User").
		Preload("User.Biodata").
		Preload("User.Roles").
		First(&e, "id = ?", enrollmentID).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("enrollment not found")
		}
		return nil, err
	}

	isAdmin := false
	for _, r := range roles {
		if r == "admin" || r == "super_admin" {
			isAdmin = true
			break
		}
	}

	if !isAdmin && e.UserID != userID {
		return nil, errors.New("unauthorized to access this enrollment")
	}

	return &e, nil
}

func (s *EnrollmentService) UnenrollCourse(userID uuid.UUID, enrollmentID uuid.UUID) error {
	var e models.Enrollment
	if err := s.db.First(&e, "id = ?", enrollmentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("enrollment not found")
		}
		return err
	}

	if e.UserID != userID {
		return errors.New("unauthorized to unenroll this course")
	}

	if e.Status == "completed" {
		return errors.New("cannot unenroll from completed course")
	}

	return s.db.Delete(&models.Enrollment{}, "id = ?", enrollmentID).Error
}

func (s *EnrollmentService) UpdateEnrollmentStatus(
	userID uuid.UUID,
	enrollmentID uuid.UUID,
	status string,
	roles []string,
) error {

	valid := map[string]bool{"active": true, "dropped": true, "completed": true}
	if !valid[status] {
		return errors.New("invalid status, must be: active, dropped, or completed")
	}

	isAdmin := false
	for _, r := range roles {
		if r == "admin" || r == "super_admin" {
			isAdmin = true
			break
		}
	}

	if !isAdmin {
		return errors.New("only admin or super_admin can update enrollment status")
	}

	var e models.Enrollment
	if err := s.db.First(&e, "id = ?", enrollmentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("enrollment not found")
		}
		return err
	}

	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}

	if status == "completed" {
		now := time.Now()
		updates["completed_datetime"] = &now
		updates["progress"] = 100.0
	}

	return s.db.Model(&models.Enrollment{}).
		Where("id = ?", enrollmentID).
		Updates(updates).Error
}

func (s *EnrollmentService) GetCourseStudents(courseID uuid.UUID, userRoles []string) ([]models.Enrollment, error) {
	allowed := map[string]bool{"admin": true, "super_admin": true, "instructor": true}
	has := false
	for _, r := range userRoles {
		if allowed[r] {
			has = true
			break
		}
	}
	if !has {
		return nil, errors.New("only instructors and admins can view course students")
	}

	var enrollments []models.Enrollment
	if err := s.db.
		Preload("User").
		Preload("User.Biodata").
		Preload("User.Roles").
		Preload("Course").
		Preload("Course.CourseType").
		Preload("Course.Creator").
		Where("course_id = ?", courseID).
		Order("enrollment_datetime DESC").
		Find(&enrollments).Error; err != nil {
		return nil, err
	}

	return enrollments, nil
}
