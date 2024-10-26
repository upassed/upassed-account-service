package student_test

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/upassed/upassed-account-service/internal/messanging/student"
	business "github.com/upassed/upassed-account-service/internal/service/model"
	"github.com/upassed/upassed-account-service/internal/util"
	"testing"
)

func TestConvertToStudentCreateRequest_InvalidBytes(t *testing.T) {
	invalidBytes := make([]byte, 10)
	_, err := student.ConvertToStudentCreateRequest(invalidBytes)
	require.NotNil(t, err)
}

func TestConvertToStudentCreateRequest_ValidBytes(t *testing.T) {
	initialRequest := util.RandomEventStudentCreateRequest()
	initialRequestBytes, err := json.Marshal(initialRequest)
	require.Nil(t, err)

	convertedRequest, err := student.ConvertToStudentCreateRequest(initialRequestBytes)
	require.Nil(t, err)

	assert.Equal(t, *initialRequest, *convertedRequest)
}

func TestConvertToStudent(t *testing.T) {
	request := util.RandomEventStudentCreateRequest()
	convertedStudent := student.ConvertToStudent(request)
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

	convertedResponse := student.ConvertToStudentCreateResponse(&responseFromService)
	require.NotNil(t, convertedResponse)

	assert.Equal(t, responseFromService.CreatedStudentID.String(), convertedResponse.CreatedStudentID)
}
