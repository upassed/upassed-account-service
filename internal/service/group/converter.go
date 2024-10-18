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

func ConvertToServiceGroup(groupToConvert domain.Group) business.Group {
	return business.Group{
		ID:                 groupToConvert.ID,
		SpecializationCode: groupToConvert.SpecializationCode,
		GroupNumber:        groupToConvert.GroupNumber,
	}
}

func ConvertToGroupFilter(filterToConvert business.GroupFilter) domain.GroupFilter {
	return domain.GroupFilter{
		SpecializationCode: filterToConvert.SpecializationCode,
		GroupNumber:        filterToConvert.GroupNumber,
	}
}

func ConvertToServiceGroups(groupsToConvert []domain.Group) []business.Group {
	convertedGroups := make([]business.Group, 0, len(groupsToConvert))
	for idx := range groupsToConvert {
		convertedGroups = append(convertedGroups, ConvertToServiceGroup(groupsToConvert[idx]))
	}

	return convertedGroups
}
