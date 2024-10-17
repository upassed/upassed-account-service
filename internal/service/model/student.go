package business

import (
	"github.com/google/uuid"
)

type Student struct {
	ID               uuid.UUID
	FirstName        string
	LastName         string
	MiddleName       string
	EducationalEmail string
	Username         string
	Group            Group
}

type StudentCreateResponse struct {
	CreatedStudentID uuid.UUID
}
