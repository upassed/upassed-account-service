package student

import (
	"github.com/google/uuid"
	"github.com/upassed/upassed-account-service/internal/service/group"
	"github.com/upassed/upassed-account-service/internal/service/student"
	"github.com/upassed/upassed-account-service/pkg/client"
)

func ConvertToStudent(request *client.StudentCreateRequest) student.Student {
	return student.Student{
		ID:               uuid.New(),
		FirstName:        request.GetFirstName(),
		LastName:         request.GetLastName(),
		MiddleName:       request.GetMiddleName(),
		EducationalEmail: request.GetEducationalEmail(),
		Username:         request.GetUsername(),
		Group: group.Group{
			ID: uuid.MustParse(request.GetGroupId()),
		},
	}
}

func ConvertToStudentCreateResponse(response student.StudentCreateResponse) *client.StudentCreateResponse {
	return &client.StudentCreateResponse{
		CreatedStudentId: response.CreatedStudentID.String(),
	}
}

func ConvertToFindByIDResponse(student student.Student) *client.StudentFindByIDResponse {
	return &client.StudentFindByIDResponse{
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
