package constants

const (
	CourseStatusDraft     = "draft"
	CourseStatusPublished = "published"
	CourseStatusPaused    = "paused"
	CourseStatusArchived  = "archived"
)

var ValidCourseStatuses = []string{
	CourseStatusDraft,
	CourseStatusPublished,
	CourseStatusPaused,
	CourseStatusArchived,
}