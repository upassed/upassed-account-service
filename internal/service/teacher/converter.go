package teacher

import (
	repository "github.com/upassed/upassed-account-service/internal/repository/teacher"
)

func ConvertToRepositoryTeacher(teacher Teacher) repository.Teacher {
	return repository.Teacher{
		ID:          teacher.ID,
		FirstName:   teacher.FirstName,
		LastName:    teacher.LastName,
		MiddleName:  teacher.MiddleName,
		ReportEmail: teacher.ReportEmail,
		Username:    teacher.Username,
	}
}

func ConvertToServiceTeacher(teacher repository.Teacher) Teacher {
	return Teacher{
		ID:          teacher.ID,
		FirstName:   teacher.FirstName,
		LastName:    teacher.LastName,
		MiddleName:  teacher.MiddleName,
		ReportEmail: teacher.ReportEmail,
		Username:    teacher.Username,
	}
}
