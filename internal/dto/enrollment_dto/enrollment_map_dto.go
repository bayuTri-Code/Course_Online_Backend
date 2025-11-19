package enrollmentdto

import (
	"course_online_backend/internal/models"

	"github.com/google/uuid"
)

// ==========================================
// ENROLLMENT MAPPERS
// ==========================================

// ToEnrollmentResponse - Convert model to response (for create/enroll)
func ToEnrollmentResponse(enrollment *models.Enrollment) *EnrollmentResponse {
	return &EnrollmentResponse{
		ID:                 enrollment.ID,
		CourseID:           enrollment.CourseID,
		UserID:             enrollment.UserID,
		EnrollmentDatetime: enrollment.EnrollmentDatetime,
		CompletedDatetime:  enrollment.CompletedDatetime,
		Status:             enrollment.Status,
		StatusPayment:      enrollment.StatusPayment,
		Progress:           enrollment.Progress,
		ExpiredDate:        enrollment.ExpiredDate,
		CreatedAt:          enrollment.CreatedAt,
		UpdatedAt:          enrollment.UpdatedAt,
		Course:             toCourseInfo(&enrollment.Course),
	}
}

// ToEnrollmentListResponse - Convert to list response (for GetMyEnrollments)
func ToEnrollmentListResponse(enrollment *models.Enrollment) *EnrollmentListResponse {
	return &EnrollmentListResponse{
		ID:                 enrollment.ID,
		EnrollmentDatetime: enrollment.EnrollmentDatetime,
		CompletedDatetime:  enrollment.CompletedDatetime,
		Status:             enrollment.Status,
		StatusPayment:      enrollment.StatusPayment,
		Progress:           enrollment.Progress,
		ExpiredDate:        enrollment.ExpiredDate,
		Course:             toCourseInfo(&enrollment.Course),
	}
}

// ToEnrollmentDetailResponse - Convert to detail response (for GetEnrollmentDetail)
func ToEnrollmentDetailResponse(enrollment *models.Enrollment, includeUser bool) *EnrollmentDetailResponse {
	response := &EnrollmentDetailResponse{
		ID:                 enrollment.ID,
		EnrollmentDatetime: enrollment.EnrollmentDatetime,
		CompletedDatetime:  enrollment.CompletedDatetime,
		Status:             enrollment.Status,
		StatusPayment:      enrollment.StatusPayment,
		Progress:           enrollment.Progress,
		ExpiredDate:        enrollment.ExpiredDate,
		CreatedAt:          enrollment.CreatedAt,
		UpdatedAt:          enrollment.UpdatedAt,
		Course:             toCourseInfo(&enrollment.Course),
	}

	// Include user info only if requested (for admin/instructor view)
	if includeUser {
		response.User = toUserInfo(&enrollment.User)
	}

	return response
}

// ToEnrollmentListResponses - Convert slice of enrollments
func ToEnrollmentListResponses(enrollments []models.Enrollment) []*EnrollmentListResponse {
	responses := make([]*EnrollmentListResponse, len(enrollments))
	for i, enrollment := range enrollments {
		responses[i] = ToEnrollmentListResponse(&enrollment)
	}
	return responses
}

// ToEnrollmentDetailResponses - Convert slice with user info
func ToEnrollmentDetailResponses(enrollments []models.Enrollment, includeUser bool) []*EnrollmentDetailResponse {
	responses := make([]*EnrollmentDetailResponse, len(enrollments))
	for i, enrollment := range enrollments {
		responses[i] = ToEnrollmentDetailResponse(&enrollment, includeUser)
	}
	return responses
}

// ==========================================
// HELPER MAPPERS
// ==========================================

// toCourseInfo - Transform course data
func toCourseInfo(course *models.Course) CourseInfo {
	info := CourseInfo{
		ID:          course.ID,
		Name:        course.Name,
		Description: course.Description,
		Thumbnail:   course.Thumbnail,
		Price:       course.Price,
		Status:      course.Status,
	}

	// Add course type if exists
	if course.CourseType != nil {
		info.CourseType = &CourseTypeInfo{
			ID:          course.CourseType.ID,
			Name:        course.CourseType.Name,
			Description: course.CourseType.Description,
		}
	}

	// Add instructor/creator if exists
	if course.Creator != nil {
		info.CreatedBy = &CreatorInfo{
			ID:       course.Creator.ID,
			Username: course.Creator.Username,
			Email:    course.Creator.EmailAddress,
		}
	}

	return info
}

// toUserInfo - Transform user data (NO SENSITIVE DATA)
func toUserInfo(user *models.User) *UserInfo {
	info := &UserInfo{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.EmailAddress,
		IsActive: user.IsActive,
	}

	// Extract role names
	if len(user.Roles) > 0 {
		roles := make([]string, len(user.Roles))
		for i, role := range user.Roles {
			roles[i] = role.Name
		}
		info.Roles = roles
	}

	zeroUUID := uuid.UUID{}
	if user.Biodata.ID != zeroUUID {
		info.Biodata = &BiodataInfo{
			Name:           user.Biodata.Name,
			Age:            user.Biodata.Age,
			School:         user.Biodata.School,
			ProfilePicture: user.Biodata.ProfilePicture,
		}
	}

	return info
}

// ToCheckEnrollmentResponse - Simple check response
func ToCheckEnrollmentResponse(isEnrolled bool, courseID uuid.UUID) *CheckEnrollmentResponse {
	return &CheckEnrollmentResponse{
		IsEnrolled: isEnrolled,
		CourseID:   courseID,
	}
}