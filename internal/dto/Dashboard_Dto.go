package dto

type AdminDashboardResponse struct {
	TotalUsers   int64  `json:"total_users"`
	TotalCourses int64  `json:"total_courses"`
}

type SuperAdminDashboardResponse struct {
	TotalInstructors int64  `json:"total_instructors"`
	TotalStudents    int64  `json:"total_students"`
	SystemStatus     string `json:"system_status"`
}

type InstructorDashboardResponse struct {
	MyCourses int64  `json:"my_courses"`
	Message   string `json:"message"`
}

type StudentDashboardResponse struct {
	EnrolledCourses int64  `json:"enrolled_courses"`
	Message         string `json:"message"`
}
