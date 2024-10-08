package converter

import (
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	business "github.com/upassed/upassed-account-service/internal/service/model"
)

func ConvertTeacher(teacher business.Teacher) domain.Teacher {
	return domain.Teacher{
		ID:          teacher.ID,
		FirstName:   teacher.FirstName,
		LastName:    teacher.LastName,
		MiddleName:  teacher.MiddleName,
		ReportEmail: teacher.ReportEmail,
		Username:    teacher.Username,
	}
}