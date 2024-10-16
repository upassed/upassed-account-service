package student

import (
	"github.com/upassed/upassed-account-service/internal/repository/group"
	"github.com/upassed/upassed-account-service/internal/repository/student"
)

func ConvertToRepositoryStudent(serviceStudent Student) student.Student {
	return student.Student{
		ID:               serviceStudent.ID,
		FirstName:        serviceStudent.FirstName,
		LastName:         serviceStudent.LastName,
		MiddleName:       serviceStudent.MiddleName,
		EducationalEmail: serviceStudent.EducationalEmail,
		Username:         serviceStudent.Username,
		Group: group.Group{
			ID: serviceStudent.Group.ID,
		},
	}
}
