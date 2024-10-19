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
	repositoryStudent := domain.Student{
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

	serviceStudent := student.ConvertToServiceStudent(repositoryStudent)
	require.NotNil(t, serviceStudent)

	assert.Equal(t, repositoryStudent.ID, serviceStudent.ID)
	assert.Equal(t, repositoryStudent.FirstName, serviceStudent.FirstName)
	assert.Equal(t, repositoryStudent.LastName, serviceStudent.LastName)
	assert.Equal(t, repositoryStudent.MiddleName, serviceStudent.MiddleName)
	assert.Equal(t, repositoryStudent.EducationalEmail, serviceStudent.EducationalEmail)
	assert.Equal(t, repositoryStudent.Username, serviceStudent.Username)
	assert.Equal(t, repositoryStudent.Group.ID, serviceStudent.Group.ID)
	assert.Equal(t, repositoryStudent.Group.SpecializationCode, serviceStudent.Group.SpecializationCode)
	assert.Equal(t, repositoryStudent.Group.GroupNumber, serviceStudent.Group.GroupNumber)
}
