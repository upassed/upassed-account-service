package student

import (
	"github.com/google/uuid"
	event "github.com/upassed/upassed-account-service/internal/messanging/model"
	business "github.com/upassed/upassed-account-service/internal/service/model"
)

func ConvertToStudent(request *event.StudentCreateRequest) business.Student {
	return business.Student{
		ID:               uuid.New(),
		FirstName:        request.FirstName,
		LastName:         request.LastName,
		MiddleName:       request.MiddleName,
		EducationalEmail: request.EducationalEmail,
		Username:         request.Username,
		Group: business.Group{
			ID: uuid.MustParse(request.GroupId),
		},
	}
}

func ConvertToStudentCreateResponse(response business.StudentCreateResponse) *event.StudentCreateResponse {
	return &event.StudentCreateResponse{
		CreatedStudentID: response.CreatedStudentID.String(),
	}
}
