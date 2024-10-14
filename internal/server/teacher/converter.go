package teacher

import (
	"github.com/google/uuid"
	"github.com/upassed/upassed-account-service/internal/service/teacher"
	"github.com/upassed/upassed-account-service/pkg/client"
)

func ConvertToTeacher(request *client.TeacherCreateRequest) teacher.Teacher {
	return teacher.Teacher{
		ID:          uuid.New(),
		FirstName:   request.GetFirstName(),
		LastName:    request.GetLastName(),
		MiddleName:  request.GetMiddleName(),
		ReportEmail: request.GetReportEmail(),
		Username:    request.GetUsername(),
	}
}

func ConvertToTeacherCreateResponse(response teacher.TeacherCreateResponse) *client.TeacherCreateResponse {
	return &client.TeacherCreateResponse{
		CreatedTeacherId: response.CreatedTeacherID.String(),
	}
}

func ConvertToFindByIDResponse(teacher teacher.Teacher) *client.TeacherFindByIDResponse {
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
