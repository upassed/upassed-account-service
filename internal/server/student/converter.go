package student

import (
	business "github.com/upassed/upassed-account-service/internal/service/model"
	"github.com/upassed/upassed-account-service/pkg/client"
)

func ConvertToFindByUsernameResponse(student *business.Student) *client.StudentFindByUsernameResponse {
	return &client.StudentFindByUsernameResponse{
		Student: &client.StudentDTO{
			Id:               student.ID.String(),
			FirstName:        student.FirstName,
			LastName:         student.LastName,
			MiddleName:       student.MiddleName,
			EducationalEmail: student.EducationalEmail,
			Username:         student.Username,
			Group: &client.GroupDTO{
				Id:                 student.Group.ID.String(),
				SpecializationCode: student.Group.SpecializationCode,
				GroupNumber:        student.Group.GroupNumber,
			},
		},
	}
}
