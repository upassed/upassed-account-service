package converter

import (
	"github.com/google/uuid"
	business "github.com/upassed/upassed-account-service/internal/service/model"
	"github.com/upassed/upassed-account-service/pkg/client"
)

func ConvertTeacherCreateRequest(request *client.TeacherCreateRequest) business.Teacher {
	return business.Teacher{
		ID:          uuid.New(),
		FirstName:   request.GetFirstName(),
		LastName:    request.GetLastName(),
		MiddleName:  request.GetMiddleName(),
		ReportEmail: request.GetReportEmail(),
		Username:    request.GetUsername(),
	}
}

func ConvertTeacherCreateResponse(response business.TeacherCreateResponse) *client.TeacherCreateResponse {
	return &client.TeacherCreateResponse{
		CreatedTeacherId: response.CreatedTeacherID.String(),
	}
}

func ConvertTeacher(teacher business.Teacher) *client.TeacherFindByIDResponse {
	return &client.TeacherFindByIDResponse{
		Teacher: &client.TeacherDTO{
			Id:          teacher.ID.String(),
			FirstName:   teacher.FirstName,
			LastName:    teacher.LastName,
			MiddleName:  teacher.MiddleName,
			ReportEmail: teacher.ReportEmail,
			Username:    teacher.Username,
		},
	}
}
