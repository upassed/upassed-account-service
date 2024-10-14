package teacher_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	repository "github.com/upassed/upassed-account-service/internal/repository/teacher"
	"github.com/upassed/upassed-account-service/internal/service/teacher"
)

func TestConvertToRepositoryTeacher(t *testing.T) {
	teacherToConvert := teacher.Teacher{
		ID:          uuid.New(),
		FirstName:   gofakeit.FirstName(),
		LastName:    gofakeit.LastName(),
		MiddleName:  gofakeit.MiddleName(),
		ReportEmail: gofakeit.Email(),
		Username:    gofakeit.Username(),
	}

	repositoryTeacher := teacher.ConvertToRepositoryTeacher(teacherToConvert)

	assert.Equal(t, teacherToConvert.ID, repositoryTeacher.ID)
	assert.Equal(t, teacherToConvert.FirstName, repositoryTeacher.FirstName)
	assert.Equal(t, teacherToConvert.LastName, repositoryTeacher.LastName)
	assert.Equal(t, teacherToConvert.MiddleName, repositoryTeacher.MiddleName)
	assert.Equal(t, teacherToConvert.ReportEmail, repositoryTeacher.ReportEmail)
	assert.Equal(t, teacherToConvert.Username, repositoryTeacher.Username)
}

func TestConvertToServiceTeacher(t *testing.T) {
	teacherToConvert := repository.Teacher{
		ID:          uuid.New(),
		FirstName:   gofakeit.FirstName(),
		LastName:    gofakeit.LastName(),
		MiddleName:  gofakeit.MiddleName(),
		ReportEmail: gofakeit.Email(),
		Username:    gofakeit.Username(),
	}

	serviceTeacher := teacher.ConvertToServiceTeacher(teacherToConvert)

	assert.Equal(t, teacherToConvert.ID, serviceTeacher.ID)
	assert.Equal(t, teacherToConvert.FirstName, serviceTeacher.FirstName)
	assert.Equal(t, teacherToConvert.LastName, serviceTeacher.LastName)
	assert.Equal(t, teacherToConvert.MiddleName, serviceTeacher.MiddleName)
	assert.Equal(t, teacherToConvert.ReportEmail, serviceTeacher.ReportEmail)
	assert.Equal(t, teacherToConvert.Username, serviceTeacher.Username)
}
