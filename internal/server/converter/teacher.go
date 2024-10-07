package converter

import (
	business "github.com/upassed/upassed-account-service/internal/service/model"
	"github.com/upassed/upassed-account-service/pkg/client"
)

func ConvertTeacherCreateRequest(response *client.TeacherCreateRequest) business.TeacherCreateRequest {
	return business.TeacherCreateRequest{
		FirstName:   response.GetFirstName(),
		LastName:    response.GetLastName(),
		MiddleName:  response.GetMiddleName(),
		ReportEmail: response.GetReportEmail(),
		Username:    response.GetUsername(),
	}
}

func TestConvertTeacherCreateResponse(response business.TeacherCreateResponse) client.TeacherCreateResponse {
	return client.TeacherCreateResponse{
		CreatedTeacherId: response.CreatedTeacherID.String(),
	}
}
