package group_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	server "github.com/upassed/upassed-account-service/internal/server/group"
	"github.com/upassed/upassed-account-service/internal/service/group"
	"github.com/upassed/upassed-account-service/internal/service/student"
	"github.com/upassed/upassed-account-service/pkg/client"
)

func TestConvertToFindStudentsInGroupResponse(t *testing.T) {
	studentsInGroup := []student.Student{randomStudent(), randomStudent(), randomStudent()}

	response := server.ConvertToFindStudentsInGroupResponse(studentsInGroup)
	require.Equal(t, len(studentsInGroup), len(response.GetStudentsInGroup()))
	for idx := range studentsInGroup {
		assertStudentsEqual(t, studentsInGroup[idx], response.GetStudentsInGroup()[idx])
	}
}

func assertStudentsEqual(t *testing.T, left student.Student, right *client.StudentDTO) {
	assert.Equal(t, left.ID.String(), right.GetId())
	assert.Equal(t, left.FirstName, right.GetFirstName())
	assert.Equal(t, left.LastName, right.GetLastName())
	assert.Equal(t, left.MiddleName, right.GetMiddleName())
	assert.Equal(t, left.EducationalEmail, right.GetEducationalEmail())
	assert.Equal(t, left.Username, right.GetUsername())
	assert.Equal(t, left.Group.ID.String(), right.GetGroup().GetId())
	assert.Equal(t, left.Group.SpecializationCode, right.GetGroup().GetSpecializationCode())
	assert.Equal(t, left.Group.GroupNumber, right.GetGroup().GetGroupNumber())
}

func randomStudent() student.Student {
	return student.Student{
		ID:               uuid.New(),
		FirstName:        gofakeit.FirstName(),
		LastName:         gofakeit.LastName(),
		MiddleName:       gofakeit.MiddleName(),
		EducationalEmail: gofakeit.Email(),
		Username:         gofakeit.Username(),
		Group: group.Group{
			ID:                 uuid.New(),
			SpecializationCode: gofakeit.WeekDay(),
			GroupNumber:        gofakeit.WeekDay(),
		},
	}
}
