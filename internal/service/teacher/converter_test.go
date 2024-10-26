package teacher_test

import (
	"github.com/upassed/upassed-account-service/internal/util"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/upassed/upassed-account-service/internal/service/teacher"
)

func TestConvertToRepositoryTeacher(t *testing.T) {
	teacherToConvert := util.RandomBusinessTeacher()
	repositoryTeacher := teacher.ConvertToRepositoryTeacher(teacherToConvert)

	assert.Equal(t, teacherToConvert.ID, repositoryTeacher.ID)
	assert.Equal(t, teacherToConvert.FirstName, repositoryTeacher.FirstName)
	assert.Equal(t, teacherToConvert.LastName, repositoryTeacher.LastName)
	assert.Equal(t, teacherToConvert.MiddleName, repositoryTeacher.MiddleName)
	assert.Equal(t, teacherToConvert.ReportEmail, repositoryTeacher.ReportEmail)
	assert.Equal(t, teacherToConvert.Username, repositoryTeacher.Username)
}

func TestConvertToServiceTeacher(t *testing.T) {
	teacherToConvert := util.RandomDomainTeacher()
	serviceTeacher := teacher.ConvertToServiceTeacher(teacherToConvert)

	assert.Equal(t, teacherToConvert.ID, serviceTeacher.ID)
	assert.Equal(t, teacherToConvert.FirstName, serviceTeacher.FirstName)
	assert.Equal(t, teacherToConvert.LastName, serviceTeacher.LastName)
	assert.Equal(t, teacherToConvert.MiddleName, serviceTeacher.MiddleName)
	assert.Equal(t, teacherToConvert.ReportEmail, serviceTeacher.ReportEmail)
	assert.Equal(t, teacherToConvert.Username, serviceTeacher.Username)
}
