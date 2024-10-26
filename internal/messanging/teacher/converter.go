package teacher

import (
	"encoding/json"
	"github.com/google/uuid"
	event "github.com/upassed/upassed-account-service/internal/messanging/model"
	business "github.com/upassed/upassed-account-service/internal/service/model"
)

func ConvertToTeacherCreateRequest(messageBody []byte) (*event.TeacherCreateRequest, error) {
	var request event.TeacherCreateRequest
	if err := json.Unmarshal(messageBody, &request); err != nil {
		return &event.TeacherCreateRequest{}, err
	}

	return &request, nil
}

func ConvertToTeacher(request *event.TeacherCreateRequest) *business.Teacher {
	return &business.Teacher{
		ID:          uuid.New(),
		FirstName:   request.FirstName,
		LastName:    request.LastName,
		MiddleName:  request.MiddleName,
		ReportEmail: request.ReportEmail,
		Username:    request.Username,
	}
}

func ConvertToTeacherCreateResponse(response *business.TeacherCreateResponse) *event.TeacherCreateResponse {
	return &event.TeacherCreateResponse{
		CreatedTeacherID: response.CreatedTeacherID.String(),
	}
}
