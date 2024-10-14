package teacher

import "github.com/google/uuid"

type Teacher struct {
	ID          uuid.UUID
	FirstName   string
	LastName    string
	MiddleName  string
	ReportEmail string
	Username    string
}

type TeacherCreateResponse struct {
	CreatedTeacherID uuid.UUID
}
