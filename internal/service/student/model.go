package student

import "github.com/google/uuid"

type Student struct {
	ID               uuid.UUID
	FirstName        string
	LastName         string
	MiddleName       string
	EducationalEmail string
	Username         string
	Group            Group
}

type Group struct {
	ID                 uuid.UUID
	SpecializationCode string
	GroupNumber        string
}

type StudentCreateResponse struct {
	CreatedStudentID uuid.UUID
}
