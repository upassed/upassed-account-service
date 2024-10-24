package student

import (
	"encoding/json"
	"github.com/google/uuid"
	event "github.com/upassed/upassed-account-service/internal/messanging/model"
	business "github.com/upassed/upassed-account-service/internal/service/model"
)

func ConvertToStudentCreateRequest(messageBody []byte) (event.StudentCreateRequest, error) {
	var request event.StudentCreateRequest
	if err := json.Unmarshal(messageBody, &request); err != nil {
		return event.StudentCreateRequest{}, err
	}

	return request, nil
}

func ConvertToStudent(request event.StudentCreateRequest) business.Student {
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

func ConvertToStudentCreateResponse(response business.StudentCreateResponse) event.StudentCreateResponse {
	return event.StudentCreateResponse{
		CreatedStudentID: response.CreatedStudentID.String(),
	}
}
