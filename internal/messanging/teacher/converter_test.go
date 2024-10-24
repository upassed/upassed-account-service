package teacher_test

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	event "github.com/upassed/upassed-account-service/internal/messanging/model"
	"github.com/upassed/upassed-account-service/internal/messanging/teacher"
	business "github.com/upassed/upassed-account-service/internal/service/model"
	"testing"
)

func TestConvertToTeacher(t *testing.T) {
	request := event.TeacherCreateRequest{
		FirstName:   gofakeit.FirstName(),
		LastName:    gofakeit.LastName(),
		MiddleName:  gofakeit.MiddleName(),
		ReportEmail: gofakeit.Email(),
		Username:    gofakeit.Username(),
	}

	convertedTeacher := teacher.ConvertToTeacher(&request)
	require.NotNil(t, convertedTeacher.ID)

	assert.Equal(t, request.FirstName, convertedTeacher.FirstName)
	assert.Equal(t, request.LastName, convertedTeacher.LastName)
	assert.Equal(t, request.MiddleName, convertedTeacher.MiddleName)
	assert.Equal(t, request.ReportEmail, convertedTeacher.ReportEmail)
	assert.Equal(t, request.Username, convertedTeacher.Username)
}

func TestConvertToTeacherCreateResponse(t *testing.T) {
	response := business.TeacherCreateResponse{
		CreatedTeacherID: uuid.New(),
	}

	convertedResponse := teacher.ConvertToTeacherCreateResponse(response)
	assert.Equal(t, response.CreatedTeacherID.String(), convertedResponse.CreatedTeacherID)
}
