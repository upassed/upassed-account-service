package student

import (
	"github.com/upassed/upassed-account-service/internal/repository/group"
	repository "github.com/upassed/upassed-account-service/internal/repository/student"
	serviceGroup "github.com/upassed/upassed-account-service/internal/service/group"
)

func ConvertToRepositoryStudent(serviceStudent Student) repository.Student {
	return repository.Student{
		ID:               serviceStudent.ID,
		FirstName:        serviceStudent.FirstName,
		LastName:         serviceStudent.LastName,
		MiddleName:       serviceStudent.MiddleName,
		EducationalEmail: serviceStudent.EducationalEmail,
		Username:         serviceStudent.Username,
		Group: group.Group{
			ID:                 serviceStudent.Group.ID,
			SpecializationCode: serviceStudent.Group.SpecializationCode,
			GroupNumber:        serviceStudent.Group.GroupNumber,
		},
	}
}

func ConvertToServiceStudent(repoStudent repository.Student) Student {
	return Student{
		ID:               repoStudent.ID,
		FirstName:        repoStudent.FirstName,
		LastName:         repoStudent.LastName,
		MiddleName:       repoStudent.MiddleName,
		EducationalEmail: repoStudent.EducationalEmail,
		Username:         repoStudent.Username,
		Group: serviceGroup.Group{
			ID:                 repoStudent.Group.ID,
			SpecializationCode: repoStudent.Group.SpecializationCode,
			GroupNumber:        repoStudent.Group.GroupNumber,
		},
	}
}
