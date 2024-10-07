package converter_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/upassed/upassed-account-service/internal/server/converter"
	business "github.com/upassed/upassed-account-service/internal/service/model"
	"github.com/upassed/upassed-account-service/pkg/client"
)

func TestConvertTeacherCreateRequest(t *testing.T) {
	request := client.TeacherCreateRequest{
		FirstName:   gofakeit.FirstName(),
		LastName:    gofakeit.LastName(),
		MiddleName:  gofakeit.MiddleName(),
		ReportEmail: gofakeit.Email(),
		Username:    gofakeit.Username(),
	}

	convertedRequest := converter.ConvertTeacherCreateRequest(&request)

	assert.Equal(t, request.GetFirstName(), convertedRequest.FirstName)
	assert.Equal(t, request.GetLastName(), convertedRequest.LastName)
	assert.Equal(t, request.GetMiddleName(), convertedRequest.MiddleName)
	assert.Equal(t, request.GetReportEmail(), convertedRequest.ReportEmail)
	assert.Equal(t, request.GetUsername(), convertedRequest.Username)
}

func TestConvertTeacherCreateResponse(t *testing.T) {
	response := business.TeacherCreateResponse{
		CreatedTeacherID: uuid.New(),
	}

	convertedResponse := converter.TestConvertTeacherCreateResponse(response)
	assert.Equal(t, response.CreatedTeacherID.String(), convertedResponse.GetCreatedTeacherId())
}
