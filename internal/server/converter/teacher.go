package converter

import (
	"github.com/google/uuid"
	business "github.com/upassed/upassed-account-service/internal/service/model"
	"github.com/upassed/upassed-account-service/pkg/client"
)

func ConvertTeacherCreateRequest(response *client.TeacherCreateRequest) business.Teacher {
	return business.Teacher{
		ID:          uuid.New(),
		FirstName:   response.GetFirstName(),
		LastName:    response.GetLastName(),
		MiddleName:  response.GetMiddleName(),
		ReportEmail: response.GetReportEmail(),
		Username:    response.GetUsername(),
	}
}

func ConvertTeacherCreateResponse(response business.TeacherCreateResponse) client.TeacherCreateResponse {
	return client.TeacherCreateResponse{
		CreatedTeacherId: response.CreatedTeacherID.String(),
	}
}
