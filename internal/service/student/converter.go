package student

import (
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	business "github.com/upassed/upassed-account-service/internal/service/model"
)

func ConvertToRepositoryStudent(serviceStudent *business.Student) *domain.Student {
	return &domain.Student{
		ID:               serviceStudent.ID,
		FirstName:        serviceStudent.FirstName,
		LastName:         serviceStudent.LastName,
		MiddleName:       serviceStudent.MiddleName,
		EducationalEmail: serviceStudent.EducationalEmail,
		Username:         serviceStudent.Username,
		Group: domain.Group{
			ID:                 serviceStudent.Group.ID,
			SpecializationCode: serviceStudent.Group.SpecializationCode,
			GroupNumber:        serviceStudent.Group.GroupNumber,
		},
	}
}

func ConvertToServiceStudent(repoStudent *domain.Student) *business.Student {
	return &business.Student{
		ID:               repoStudent.ID,
		FirstName:        repoStudent.FirstName,
		LastName:         repoStudent.LastName,
		MiddleName:       repoStudent.MiddleName,
		EducationalEmail: repoStudent.EducationalEmail,
		Username:         repoStudent.Username,
		Group: business.Group{
			ID:                 repoStudent.Group.ID,
			SpecializationCode: repoStudent.Group.SpecializationCode,
			GroupNumber:        repoStudent.Group.GroupNumber,
		},
	}
}
