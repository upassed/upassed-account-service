package teacher_test

import (
	"encoding/json"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	event "github.com/upassed/upassed-account-service/internal/messanging/model"
	"github.com/upassed/upassed-account-service/internal/messanging/teacher"
	business "github.com/upassed/upassed-account-service/internal/service/model"
	"github.com/upassed/upassed-account-service/internal/util"
	"testing"
)

func TestConvertToTeacherCreateRequest_InvalidBytes(t *testing.T) {
	invalidBytes := make([]byte, 10)
	_, err := teacher.ConvertToTeacherCreateRequest(invalidBytes)
	require.NotNil(t, err)
}

func TestConvertToTeacherCreateRequest_ValidBytes(t *testing.T) {
	initialRequest := util.RandomEventTeacherCreateRequest()
	initialRequestBytes, err := json.Marshal(initialRequest)
	require.Nil(t, err)

	convertedRequest, err := teacher.ConvertToTeacherCreateRequest(initialRequestBytes)
	require.Nil(t, err)

	assert.Equal(t, *initialRequest, *convertedRequest)
}

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

	convertedResponse := teacher.ConvertToTeacherCreateResponse(&response)
	assert.Equal(t, response.CreatedTeacherID.String(), convertedResponse.CreatedTeacherID)
}
