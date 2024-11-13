package teacher

import (
	business "github.com/upassed/upassed-account-service/internal/service/model"
	"github.com/upassed/upassed-account-service/pkg/client"
)

func ConvertToFindByUsernameResponse(teacher *business.Teacher) *client.TeacherFindByUsernameResponse {
	return &client.TeacherFindByUsernameResponse{
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
