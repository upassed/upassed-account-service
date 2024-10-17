package group

import (
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	business "github.com/upassed/upassed-account-service/internal/service/model"
	"github.com/upassed/upassed-account-service/internal/service/student"
)

func ConvertToServiceStudents(studentsToConvert []domain.Student) []business.Student {
	result := make([]business.Student, 0, len(studentsToConvert))
	for _, studentToConvert := range studentsToConvert {
		result = append(result, student.ConvertToServiceStudent(studentToConvert))
	}

	return result
}
