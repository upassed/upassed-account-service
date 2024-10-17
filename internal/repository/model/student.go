package domain

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
	GroupID          uuid.UUID
	Group            Group `gorm:"foreignKey:GroupID;references:ID"`
}

func (Student) TableName() string {
	return "student"
}
