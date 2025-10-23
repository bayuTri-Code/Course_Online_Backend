package services

import (
    "course_online_backend/internal/models"
    "gorm.io/gorm"
)

type DashboardService struct {
    DB *gorm.DB
}

func NewDashboardService(db *gorm.DB) *DashboardService {
    return &DashboardService{DB: db}
}

func (s *DashboardService) GetAdminDashboard() (map[string]interface{}, error) {
    var userCount int64
    var courseCount int64
    s.DB.Model(&models.User{}).Count(&userCount)
    s.DB.Model(&models.Course{}).Count(&courseCount)

    return map[string]interface{}{
        "total_users":   userCount,
        "total_courses": courseCount,
    }, nil
}

func (s *DashboardService) GetSuperAdminDashboard() (map[string]interface{}, error) {
    var instructorCount, studentCount int64
    s.DB.Model(&models.User{}).Where("role = ?", "instructor").Count(&instructorCount)
    s.DB.Model(&models.User{}).Where("role = ?", "student").Count(&studentCount)

    return map[string]interface{}{
        "total_instructors": instructorCount,
        "total_students":    studentCount,
        "system_status":     "All systems operational",
    }, nil
}

func (s *DashboardService) GetStudentDashboard(studentID uint) (map[string]interface{}, error) {
    var enrolledCount int64
    s.DB.Model(&models.Enrollment{}).Where("student_id = ?", studentID).Count(&enrolledCount)

    return map[string]interface{}{
        "enrolled_courses": enrolledCount,
        "message":          "Keep learning and growing!",
    }, nil
}

func (s *DashboardService) GetInstructorDashboard(instructorID uint) (map[string]interface{}, error) {
    var courseCount int64
    s.DB.Model(&models.Course{}).Where("instructor_id = ?", instructorID).Count(&courseCount)

    return map[string]interface{}{
        "my_courses": courseCount,
        "message":    "Welcome back, Instructor!",
    }, nil
}