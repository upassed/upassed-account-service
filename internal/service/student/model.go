package student

import (
	"github.com/google/uuid"
	"github.com/upassed/upassed-account-service/internal/service/group"
)

type Student struct {
	ID               uuid.UUID
	FirstName        string
	LastName         string
	MiddleName       string
	EducationalEmail string
	Username         string
	Group            group.Group
}

type StudentCreateResponse struct {
	CreatedStudentID uuid.UUID
}
