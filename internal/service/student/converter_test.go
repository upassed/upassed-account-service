package student_test

import (
	"github.com/upassed/upassed-account-service/internal/util"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/upassed/upassed-account-service/internal/service/student"
)

func TestConvertToRepositoryStudent(t *testing.T) {
	serviceStudent := util.RandomBusinessStudent()
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
	repositoryStudent := util.RandomDomainStudent()
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
