package event

type StudentCreateRequest struct {
	FirstName        string `json:"first_name,omitempty"`
	LastName         string `json:"last_name,omitempty"`
	MiddleName       string `json:"middle_name,omitempty"`
	EducationalEmail string `json:"educational_email,omitempty"`
	Username         string `json:"username,omitempty"`
	GroupId          string `json:"group_id,omitempty"`
}

type StudentCreateResponse struct {
	CreatedStudentID string `json:"created_student_id,omitempty"`
}
