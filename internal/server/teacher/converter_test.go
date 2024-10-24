package teacher_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/upassed/upassed-account-service/internal/server/teacher"
	business "github.com/upassed/upassed-account-service/internal/service/model"
)

func TestConvertToFindByIDResponse(t *testing.T) {
	teacherToConvert := business.Teacher{
		ID:          uuid.New(),
		FirstName:   gofakeit.FirstName(),
		LastName:    gofakeit.LastName(),
		MiddleName:  gofakeit.MiddleName(),
		ReportEmail: gofakeit.Email(),
		Username:    gofakeit.Username(),
	}

	response := teacher.ConvertToFindByIDResponse(teacherToConvert)
	require.NotNil(t, response.GetTeacher())

	assert.Equal(t, teacherToConvert.FirstName, response.GetTeacher().GetFirstName())
	assert.Equal(t, teacherToConvert.LastName, response.GetTeacher().GetLastName())
	assert.Equal(t, teacherToConvert.MiddleName, response.GetTeacher().GetMiddleName())
	assert.Equal(t, teacherToConvert.ReportEmail, response.GetTeacher().GetReportEmail())
	assert.Equal(t, teacherToConvert.Username, response.GetTeacher().GetUsername())
}
