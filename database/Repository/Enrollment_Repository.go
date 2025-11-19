package repository

import (
	"course_online_backend/internal/models"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EnrollmentRepository interface {
	Create(enrollment *models.Enrollment) error
	GetByID(id uuid.UUID) (*models.Enrollment, error)
	GetByUserAndCourse(userID, courseID uuid.UUID) (*models.Enrollment, error)
	CheckExists(userID, courseID uuid.UUID) (bool, error)
	GetUserEnrollments(userID uuid.UUID) ([]models.Enrollment, error)
	GetCourseEnrollments(courseID uuid.UUID) ([]models.Enrollment, error)
	Update(enrollment *models.Enrollment) error
	Delete(id uuid.UUID) error
	CountByCourse(courseID uuid.UUID) (int64, error)
	UpdateProgress(enrollmentID uuid.UUID, progress float64) error
	UpdateStatus(enrollmentID uuid.UUID, status string) error
}

type enrollmentRepository struct {
	db *gorm.DB
}

func NewEnrollmentRepository(db *gorm.DB) EnrollmentRepository {
	return &enrollmentRepository{db: db}
}

func (r *enrollmentRepository) Create(enrollment *models.Enrollment) error {
	if enrollment.ID == uuid.Nil {
		enrollment.ID = uuid.New()
	}
	return r.db.Create(enrollment).Error
}

func (r *enrollmentRepository) GetByID(id uuid.UUID) (*models.Enrollment, error) {
	var enrollment models.Enrollment
	err := r.db.
		Preload("Course").
		Preload("User").
		First(&enrollment, "id = ?", id).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("enrollment not found")
		}
		return nil, err
	}
	return &enrollment, nil
}

func (r *enrollmentRepository) GetByUserAndCourse(userID, courseID uuid.UUID) (*models.Enrollment, error) {
	var enrollment models.Enrollment
	err := r.db.
		Preload("Course").
		Preload("User").
		Where("user_id = ? AND course_id = ?", userID, courseID).
		First(&enrollment).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil 
		}
		return nil, err
	}
	return &enrollment, nil
}

func (r *enrollmentRepository) CheckExists(userID, courseID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&models.Enrollment{}).
		Where("user_id = ? AND course_id = ?", userID, courseID).
		Count(&count).Error
	
	return count > 0, err
}

func (r *enrollmentRepository) GetUserEnrollments(userID uuid.UUID) ([]models.Enrollment, error) {
	var enrollments []models.Enrollment
	err := r.db.
		Preload("Course").
		Preload("Course.Instructor").
		Where("user_id = ?", userID).
		Order("enrollment_datetime DESC").
		Find(&enrollments).Error
	
	return enrollments, err
}

func (r *enrollmentRepository) GetCourseEnrollments(courseID uuid.UUID) ([]models.Enrollment, error) {
	var enrollments []models.Enrollment
	err := r.db.
		Preload("User").
		Preload("User.Biodata").
		Preload("User.Roles").
		Preload("Course").
		Preload("Course.CourseType").
		Preload("Course.Creator").
		Where("course_id = ?", courseID).
		Order("enrollment_datetime DESC").
		Find(&enrollments).Error
	
	return enrollments, err
}

func (r *enrollmentRepository) Update(enrollment *models.Enrollment) error {
	return r.db.Save(enrollment).Error
}

func (r *enrollmentRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Enrollment{}, "id = ?", id).Error
}

func (r *enrollmentRepository) CountByCourse(courseID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&models.Enrollment{}).
		Where("course_id = ? AND status = ?", courseID, "active").
		Count(&count).Error
	
	return count, err
}

func (r *enrollmentRepository) UpdateProgress(enrollmentID uuid.UUID, progress float64) error {
	updates := map[string]interface{}{
		"progress":   progress,
		"updated_at": time.Now(),
	}
	
	if progress >= 100 {
		now := time.Now()
		updates["completed_datetime"] = &now
		updates["status"] = "completed"
	}
	
	return r.db.Model(&models.Enrollment{}).
		Where("id = ?", enrollmentID).
		Updates(updates).Error
}

func (r *enrollmentRepository) UpdateStatus(enrollmentID uuid.UUID, status string) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}
	
	if status == "completed" {
		now := time.Now()
		updates["completed_datetime"] = &now
		updates["progress"] = 100.0
	}
	
	return r.db.Model(&models.Enrollment{}).
		Where("id = ?", enrollmentID).
		Updates(updates).Error
}