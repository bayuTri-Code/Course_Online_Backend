package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	FirebaseUID  string    `gorm:"unique;not null" json:"firebase_uid"`
	Username     string    `gorm:"size:100" json:"username"`
	EmailAddress string    `gorm:"unique;not null" json:"email_address"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	LastLogin    time.Time `json:"last_login"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`

	Roles            []Role            `gorm:"many2many:user_roles;joinForeignKey:UserID;joinReferences:RoleID" json:"roles"`
	Biodata          Biodata           `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"biodata"`
	Activities       []Activity        `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"activities"`
	CreatedCourses   []Course          `gorm:"foreignKey:CreatedBy;constraint:OnDelete:SET NULL" json:"created_courses"`
	CreatedQuizzes   []Quiz            `gorm:"foreignKey:CreatedBy;constraint:OnDelete:SET NULL" json:"created_quizzes"`
	Enrollments      []Enrollment      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"enrollments"`
	QuizAttempts     []UserQuizAttempt `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"quiz_attempts"`
	CompletedLessons []UserLesson      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"completed_lessons"`
}

type Role struct {
	ID   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name string    `gorm:"size:50;unique;not null" json:"name"`

	Users []User `gorm:"many2many:user_roles;joinForeignKey:RoleID;joinReferences:UserID" json:"-"`
}

type UserRole struct {
	ID     uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_user_role,composite:user_role" json:"user_id"`
	RoleID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_user_role,composite:user_role" json:"role_id"`
	User   User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Role   Role      `gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE" json:"-"`
}

type Biodata struct {
	Id             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID         uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
	Name           string    `gorm:"size:100;not null" json:"name"`
	ProfilePicture string    `gorm:"size:255" json:"profile_picture"`
	Age            int       `json:"age"`
	School         string    `gorm:"size:100" json:"school"`
}

type Activity struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ActivityName string    `gorm:"size:100;not null" json:"activity_name"`
	UserID       uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	When         time.Time `gorm:"not null" json:"when"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user"`
}

type CourseType struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`

	Courses []Course `gorm:"foreignKey:CourseTypeID;constraint:OnDelete:SET NULL" json:"courses"`
}

type Course struct {
	ID                uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name              string     `gorm:"size:200;not null" json:"name"`
	Description       string     `gorm:"type:text" json:"description"`
	Price             float64    `gorm:"type:decimal(10,2)" json:"price"`
	IsProgressLimited bool       `gorm:"default:false" json:"is_progress_limited"`
	CourseTypeID      *uuid.UUID `gorm:"type:uuid;index" json:"course_type_id"`
	CreatedBy         *uuid.UUID `gorm:"type:uuid;index" json:"created_by"`
	CreatedAt         time.Time  `gorm:"autoCreateTime" json:"created_at"`

	CourseType  *CourseType  `gorm:"foreignKey:CourseTypeID;constraint:OnDelete:SET NULL" json:"course_type"`
	Creator     *User        `gorm:"foreignKey:CreatedBy;constraint:OnDelete:SET NULL" json:"creator"`
	Modules     []Module     `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE" json:"modules"`
	Quizzes     []Quiz       `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE" json:"quizzes"`
	Enrollments []Enrollment `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE" json:"enrollments"`
	Zoom        *Zoom        `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE" json:"zoom"`
}

type Module struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CourseID uuid.UUID `gorm:"type:uuid;not null;index" json:"course_id"`
	Name     string    `gorm:"size:200;not null" json:"name"`
	Number   int       `gorm:"not null" json:"number"`

	Course  Course   `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE" json:"course"`
	Lessons []Lesson `gorm:"foreignKey:ModuleID;constraint:OnDelete:CASCADE" json:"lessons"`
}

type Lesson struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ModuleID     uuid.UUID `gorm:"type:uuid;not null;index" json:"module_id"`
	VideoID      string    `gorm:"size:100" json:"video_id"`
	Name         string    `gorm:"size:200;not null" json:"name"`
	VideoDetails string    `gorm:"type:text" json:"video_details"`
	CourseOrder  int       `json:"course_order"`

	Module          Module       `gorm:"foreignKey:ModuleID;constraint:OnDelete:CASCADE" json:"module"`
	UserCompletions []UserLesson `gorm:"foreignKey:LessonID;constraint:OnDelete:CASCADE" json:"user_completions"`
}

type Enrollment struct {
	ID                 uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CourseID           uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex:idx_user_course,composite:user_course" json:"course_id"`
	UserID             uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex:idx_user_course,composite:user_course" json:"user_id"`
	EnrollmentDatetime time.Time  `gorm:"not null" json:"enrollment_datetime"`
	CompletedDatetime  *time.Time `json:"completed_datetime"`

	Course Course `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE" json:"course"`
	User   User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user"`
}

type Quiz struct {
	ID             uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CourseID       uuid.UUID  `gorm:"type:uuid;not null;index" json:"course_id"`
	Name           string     `gorm:"size:200;not null" json:"name"`
	Number         int        `gorm:"not null" json:"number"`
	MinPassScore   int        `gorm:"not null" json:"min_pass_score"`
	IsPassRequired bool       `gorm:"default:false" json:"is_pass_required"`
	CreatedBy      *uuid.UUID `gorm:"type:uuid;index" json:"created_by"`

	Course    Course            `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE" json:"course"`
	Creator   *User             `gorm:"foreignKey:CreatedBy;constraint:OnDelete:SET NULL" json:"creator"`
	Questions []QuizQuestion    `gorm:"foreignKey:QuizID;constraint:OnDelete:CASCADE" json:"questions"`
	Attempts  []UserQuizAttempt `gorm:"foreignKey:QuizID;constraint:OnDelete:CASCADE" json:"attempts"`
}

type QuizQuestion struct {
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	QuizID        uuid.UUID `gorm:"type:uuid;not null;index" json:"quiz_id"`
	QuestionTitle string    `gorm:"type:text;not null" json:"question_title"`

	Quiz    Quiz         `gorm:"foreignKey:QuizID;constraint:OnDelete:CASCADE" json:"quiz"`
	Answers []QuizAnswer `gorm:"foreignKey:QuestionID;constraint:OnDelete:CASCADE" json:"answers"`
}

type QuizAnswer struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	QuestionID uuid.UUID `gorm:"type:uuid;not null;index" json:"question_id"`
	AnswerText string    `gorm:"type:text;not null" json:"answer_text"`
	IsCorrect  bool      `gorm:"default:false" json:"is_correct"`

	Question QuizQuestion `gorm:"foreignKey:QuestionID;constraint:OnDelete:CASCADE" json:"question"`
}

type UserQuizAttempt struct {
	ID              uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID          uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	QuizID          uuid.UUID `gorm:"type:uuid;not null;index" json:"quiz_id"`
	AttemptDatetime time.Time `gorm:"not null" json:"attempt_datetime"`
	ScoreAchieved   int       `gorm:"not null" json:"score_achieved"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user"`
	Quiz Quiz `gorm:"foreignKey:QuizID;constraint:OnDelete:CASCADE" json:"quiz"`
}

type UserLesson struct {
	ID                uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID            uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_user_lesson,composite:user_lesson" json:"user_id"`
	LessonID          uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_user_lesson,composite:user_lesson" json:"lesson_id"`
	CompletedDatetime time.Time `gorm:"not null" json:"completed_datetime"`

	User   User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user"`
	Lesson Lesson `gorm:"foreignKey:LessonID;constraint:OnDelete:CASCADE" json:"lesson"`
}

type Zoom struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Link     string    `gorm:"size:500;not null" json:"link"`
	CourseID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"course_id"`

	Course Course `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE" json:"course"`
}
