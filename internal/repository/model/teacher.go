package domain

import "github.com/google/uuid"

type Teacher struct {
	ID          uuid.UUID
	FirstName   string
	LastName    string
	MiddleName  string
	ReportEmail string
	Username    string
}

// TableName overrides the table name used by Teacher to `teacher`
func (Teacher) TableName() string {
	return "teacher"
}
