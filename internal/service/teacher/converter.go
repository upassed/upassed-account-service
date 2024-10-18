package teacher

import (
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	business "github.com/upassed/upassed-account-service/internal/service/model"
)

func ConvertToRepositoryTeacher(teacher business.Teacher) domain.Teacher {
	return domain.Teacher{
		ID:          teacher.ID,
		FirstName:   teacher.FirstName,
		LastName:    teacher.LastName,
		MiddleName:  teacher.MiddleName,
		ReportEmail: teacher.ReportEmail,
		Username:    teacher.Username,
	}
}

func ConvertToServiceTeacher(teacher domain.Teacher) business.Teacher {
	return business.Teacher{
		ID:          teacher.ID,
		FirstName:   teacher.FirstName,
		LastName:    teacher.LastName,
		MiddleName:  teacher.MiddleName,
		ReportEmail: teacher.ReportEmail,
		Username:    teacher.Username,
	}
}
