package student_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/upassed/upassed-account-service/internal/service/group"
	"github.com/upassed/upassed-account-service/internal/service/student"
)

func TestConvertToRepositoryStudent(t *testing.T) {
	serviceStudent := student.Student{
		ID:               uuid.New(),
		FirstName:        gofakeit.FirstName(),
		LastName:         gofakeit.LastName(),
		MiddleName:       gofakeit.MiddleName(),
		EducationalEmail: gofakeit.Email(),
		Username:         gofakeit.Username(),
		Group: group.Group{
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
}
