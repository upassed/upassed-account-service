package teacher_test

import (
	"github.com/upassed/upassed-account-service/internal/util"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/upassed/upassed-account-service/internal/server/teacher"
)

func TestConvertToFindByIDResponse(t *testing.T) {
	teacherToConvert := util.RandomBusinessTeacher()
	response := teacher.ConvertToFindByIDResponse(teacherToConvert)
	require.NotNil(t, response.GetTeacher())

	assert.Equal(t, teacherToConvert.FirstName, response.GetTeacher().GetFirstName())
	assert.Equal(t, teacherToConvert.LastName, response.GetTeacher().GetLastName())
	assert.Equal(t, teacherToConvert.MiddleName, response.GetTeacher().GetMiddleName())
	assert.Equal(t, teacherToConvert.ReportEmail, response.GetTeacher().GetReportEmail())
	assert.Equal(t, teacherToConvert.Username, response.GetTeacher().GetUsername())
}
