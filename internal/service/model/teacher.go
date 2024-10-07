package business

import "github.com/google/uuid"

type TeacherCreateRequest struct {
	FirstName   string
	LastName    string
	MiddleName  string
	ReportEmail string
	Username    string
}

type TeacherCreateResponse struct {
	CreatedTeacherID uuid.UUID
}
