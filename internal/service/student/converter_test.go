package student_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	business "github.com/upassed/upassed-account-service/internal/service/model"
	"github.com/upassed/upassed-account-service/internal/service/student"
)

func TestConvertToRepositoryStudent(t *testing.T) {
	serviceStudent := business.Student{
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

	convertedStudent := student.ConvertToRepositoryStudent(serviceStudent)
	require.NotNil(t, convertedStudent)

	assert.Equal(t, serviceStudent.ID, convertedStudent.ID)
	assert.Equal(t, serviceStudent.FirstName, convertedStudent.FirstName)
	assert.Equal(t, serviceStudent.LastName, convertedStudent.LastName)
	assert.Equal(t, serviceStudent.MiddleName, convertedStudent.MiddleName)
	assert.Equal(t, serviceStudent.EducationalEmail, convertedStudent.EducationalEmail)
	assert.Equal(t, serviceStudent.Username, convertedStudent.Username)
	assert.Equal(t, serviceStudent.Group.ID, convertedStudent.Group.ID)
	assert.Equal(t, serviceStudent.Group.SpecializationCode, convertedStudent.Group.SpecializationCode)
	assert.Equal(t, serviceStudent.Group.GroupNumber, convertedStudent.Group.GroupNumber)
}

func TestConvertToServiceStudent(t *testing.T) {
	repositoryrStudent := domain.Student{
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

	serviceStudent := student.ConvertToServiceStudent(repositoryrStudent)
	require.NotNil(t, serviceStudent)

	assert.Equal(t, repositoryrStudent.ID, serviceStudent.ID)
	assert.Equal(t, repositoryrStudent.FirstName, serviceStudent.FirstName)
	assert.Equal(t, repositoryrStudent.LastName, serviceStudent.LastName)
	assert.Equal(t, repositoryrStudent.MiddleName, serviceStudent.MiddleName)
	assert.Equal(t, repositoryrStudent.EducationalEmail, serviceStudent.EducationalEmail)
	assert.Equal(t, repositoryrStudent.Username, serviceStudent.Username)
	assert.Equal(t, repositoryrStudent.Group.ID, serviceStudent.Group.ID)
	assert.Equal(t, repositoryrStudent.Group.SpecializationCode, serviceStudent.Group.SpecializationCode)
	assert.Equal(t, repositoryrStudent.Group.GroupNumber, serviceStudent.Group.GroupNumber)
}
