package student_test

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	event "github.com/upassed/upassed-account-service/internal/messanging/model"
	"github.com/upassed/upassed-account-service/internal/messanging/student"
	business "github.com/upassed/upassed-account-service/internal/service/model"
	"testing"
)

func TestConvertToStudent(t *testing.T) {
	request := event.StudentCreateRequest{
		FirstName:        gofakeit.FirstName(),
		LastName:         gofakeit.LastName(),
		MiddleName:       gofakeit.MiddleName(),
		EducationalEmail: gofakeit.Email(),
		Username:         gofakeit.Username(),
		GroupId:          uuid.NewString(),
	}

	convertedStudent := student.ConvertToStudent(&request)
	require.NotNil(t, convertedStudent)

	assert.Equal(t, request.FirstName, convertedStudent.FirstName)
	assert.Equal(t, request.LastName, convertedStudent.LastName)
	assert.Equal(t, request.MiddleName, convertedStudent.MiddleName)
	assert.Equal(t, request.EducationalEmail, convertedStudent.EducationalEmail)
	assert.Equal(t, request.Username, convertedStudent.Username)
	assert.Equal(t, request.GroupId, convertedStudent.Group.ID.String())
}

func TestConvertToStudentCreateResponse(t *testing.T) {
	responseFromService := business.StudentCreateResponse{
		CreatedStudentID: uuid.New(),
	}

	convertedResponse := student.ConvertToStudentCreateResponse(responseFromService)
	require.NotNil(t, convertedResponse)

	assert.Equal(t, responseFromService.CreatedStudentID.String(), convertedResponse.CreatedStudentID)
}
