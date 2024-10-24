package student_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/upassed/upassed-account-service/internal/server/student"
	business "github.com/upassed/upassed-account-service/internal/service/model"
)

func TestConvertToFindByIDResponse(t *testing.T) {
	studentToConvert := business.Student{
		ID:               uuid.New(),
		FirstName:        gofakeit.FirstName(),
		LastName:         gofakeit.LastName(),
		MiddleName:       gofakeit.MiddleName(),
		EducationalEmail: gofakeit.Email(),
		Username:         gofakeit.Username(),
		Group: business.Group{
			ID:                 uuid.New(),
			SpecializationCode: gofakeit.WeekDay(),
			GroupNumber:        gofakeit.WeekDay(),
		},
	}

	response := student.ConvertToFindByIDResponse(studentToConvert)
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
