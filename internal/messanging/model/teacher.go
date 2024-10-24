package event

type TeacherCreateRequest struct {
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	MiddleName  string `json:"middle_name,omitempty"`
	ReportEmail string `json:"report_email,omitempty"`
	Username    string `json:"username,omitempty"`
}

type TeacherCreateResponse struct {
	CreatedTeacherID string `json:"created_teacher_id,omitempty"`
}
