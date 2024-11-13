package student_test

import (
	"github.com/upassed/upassed-account-service/internal/util"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/upassed/upassed-account-service/internal/server/student"
)

func TestConvertToFindByUsernameResponse(t *testing.T) {
	studentToConvert := util.RandomBusinessStudent()
	response := student.ConvertToFindByUsernameResponse(studentToConvert)
	require.NotNil(t, response)

	assert.Equal(t, studentToConvert.ID.String(), response.GetStudent().GetId())
	assert.Equal(t, studentToConvert.FirstName, response.GetStudent().GetFirstName())
	assert.Equal(t, studentToConvert.LastName, response.GetStudent().GetLastName())
	assert.Equal(t, studentToConvert.MiddleName, response.GetStudent().GetMiddleName())
	assert.Equal(t, studentToConvert.EducationalEmail, response.GetStudent().GetEducationalEmail())
	assert.Equal(t, studentToConvert.Username, response.GetStudent().GetUsername())
	assert.Equal(t, studentToConvert.Group.ID.String(), response.GetStudent().GetGroup().GetId())
	assert.Equal(t, studentToConvert.Group.SpecializationCode, response.GetStudent().GetGroup().GetSpecializationCode())
	assert.Equal(t, studentToConvert.Group.GroupNumber, response.GetStudent().GetGroup().GetGroupNumber())
}
