package util

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	business "github.com/upassed/upassed-account-service/internal/service/model"
)

func RandomDomainStudent() domain.Student {
	return domain.Student{
		ID:               uuid.New(),
		FirstName:        gofakeit.FirstName(),
		LastName:         gofakeit.LastName(),
		MiddleName:       gofakeit.MiddleName(),
		EducationalEmail: gofakeit.Email(),
		Username:         gofakeit.Username(),
		Group: domain.Group{
			ID:                 uuid.New(),
			SpecializationCode: gofakeit.WeekDay(),
			GroupNumber:        gofakeit.WeekDay(),
		},
	}
}

func RandomBusinessStudent() business.Student {
	return business.Student{
		ID:               uuid.New(),
		FirstName:        gofakeit.FirstName(),
		LastName:         gofakeit.LastName(),
		MiddleName:       gofakeit.MiddleName(),
		EducationalEmail: gofakeit.Email(),
		Username:         gofakeit.Username(),
		Group: business.Group{
			ID:                 uuid.New(),
			SpecializationCode: gofakeit.WeekDay(),
			GroupNumber:        gofakeit.WeekDay(),
		},
	}
}

func RandomDomainTeacher() domain.Teacher {
	return domain.Teacher{
		ID:          uuid.New(),
		FirstName:   gofakeit.FirstName(),
		LastName:    gofakeit.LastName(),
		MiddleName:  gofakeit.MiddleName(),
		ReportEmail: gofakeit.Email(),
		Username:    gofakeit.Username(),
	}
}

func RandomBusinessTeacher() business.Teacher {
	return business.Teacher{
		ID:          uuid.New(),
		FirstName:   gofakeit.FirstName(),
		LastName:    gofakeit.LastName(),
		MiddleName:  gofakeit.MiddleName(),
		ReportEmail: gofakeit.Email(),
		Username:    gofakeit.Username(),
	}
}

func RandomDomainGroup() domain.Group {
	return domain.Group{
		ID:                 uuid.New(),
		SpecializationCode: gofakeit.WeekDay(),
		GroupNumber:        gofakeit.WeekDay(),
	}
}

func RandomBusinessGroup() business.Group {
	return business.Group{
		ID:                 uuid.New(),
		SpecializationCode: gofakeit.WeekDay(),
		GroupNumber:        gofakeit.WeekDay(),
	}
}
