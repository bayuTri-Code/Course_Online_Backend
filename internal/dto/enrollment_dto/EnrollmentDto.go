package enrollmentdto

import (
	"time"

	"github.com/google/uuid"
)

type EnrollFreeCourseRequest struct {
	CourseID string `json:"course_id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000"`
}

type UpdateEnrollmentStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=active dropped completed" example:"completed"`
}

type CheckEnrollmentResponse struct {
	IsEnrolled bool      `json:"is_enrolled" example:"true"`
	CourseID   uuid.UUID `json:"course_id" example:"123e4567-e89b-12d3-a456-426614174000"`
}

// response
type EnrollmentResponse struct {
	ID                 uuid.UUID  `json:"id"`
	CourseID           uuid.UUID  `json:"course_id"`
	UserID             uuid.UUID  `json:"user_id"`
	EnrollmentDatetime time.Time  `json:"enrollment_datetime"`
	CompletedDatetime  *time.Time `json:"completed_datetime,omitempty"`
	Status             string     `json:"status" example:"active"`
	StatusPayment      string     `json:"status_payment" example:"free"`
	Progress           float64    `json:"progress" example:"0"`
	ExpiredDate        *time.Time `json:"expired_date,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	Course             CourseInfo `json:"course"`
}

type EnrollmentListResponse struct {
	ID                 uuid.UUID  `json:"id"`
	EnrollmentDatetime time.Time  `json:"enrollment_datetime"`
	CompletedDatetime  *time.Time `json:"completed_datetime,omitempty"`
	Status             string     `json:"status"`
	StatusPayment      string     `json:"status_payment"`
	Progress           float64    `json:"progress"`
	ExpiredDate        *time.Time `json:"expired_date,omitempty"`
	Course             CourseInfo `json:"course"`
}

type EnrollmentDetailResponse struct {
	ID                 uuid.UUID  `json:"id"`
	EnrollmentDatetime time.Time  `json:"enrollment_datetime"`
	CompletedDatetime  *time.Time `json:"completed_datetime,omitempty"`
	Status             string     `json:"status"`
	StatusPayment      string     `json:"status_payment"`
	Progress           float64    `json:"progress"`
	ExpiredDate        *time.Time `json:"expired_date,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	Course             CourseInfo `json:"course"`
	User               *UserInfo  `json:"user,omitempty"`
}

type CourseInfo struct {
	ID          uuid.UUID       `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Thumbnail   string          `json:"thumbnail"`
	Price       float64         `json:"price"`
	Status      string          `json:"status"`
	CourseType  *CourseTypeInfo `json:"course_type,omitempty"`
	CreatedBy   *CreatorInfo    `json:"created_by,omitempty"`
}

type CourseTypeInfo struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
}

type CreatorInfo struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email,omitempty"`
}

type UserInfo struct {
	ID       uuid.UUID    `json:"id"`
	Username string       `json:"username"`
	Email    string       `json:"email"`
	IsActive bool         `json:"is_active"`
	Roles    []string     `json:"roles,omitempty"`
	Biodata  *BiodataInfo `json:"biodata,omitempty"`
}

type BiodataInfo struct {
	Name           string `json:"name"`
	Age            int    `json:"age,omitempty"`
	School         string `json:"school,omitempty"`
	ProfilePicture string `json:"profile_picture,omitempty"`
}
