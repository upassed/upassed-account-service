package student

import (
	"github.com/google/uuid"
	"github.com/upassed/upassed-account-service/internal/repository/group"
)

type Student struct {
	ID               uuid.UUID
	FirstName        string
	LastName         string
	MiddleName       string
	EducationalEmail string
	Username         string
	GroupID          uuid.UUID
	Group            group.Group `gorm:"foreignKey:GroupID;references:ID"`
}

func (Student) TableName() string {
	return "student"
}
