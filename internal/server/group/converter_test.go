package group_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/upassed/upassed-account-service/internal/server/group"
	business "github.com/upassed/upassed-account-service/internal/service/model"
	"github.com/upassed/upassed-account-service/pkg/client"
)

func TestConvertToFindStudentsInGroupResponse(t *testing.T) {
	studentsInGroup := []business.Student{randomStudent(), randomStudent(), randomStudent()}

	response := group.ConvertToFindStudentsInGroupResponse(studentsInGroup)
	require.Equal(t, len(studentsInGroup), len(response.GetStudentsInGroup()))
	for idx := range studentsInGroup {
		assertStudentsEqual(t, studentsInGroup[idx], response.GetStudentsInGroup()[idx])
	}
}

func TestConvertToFindByIDResponse(t *testing.T) {
	groupToConvert := business.Group{
		ID:                 uuid.New(),
		SpecializationCode: gofakeit.WeekDay(),
		GroupNumber:        gofakeit.WeekDay(),
	}

	response := group.ConvertToFindByIDResponse(groupToConvert)
	require.NotNil(t, response)

	assert.Equal(t, groupToConvert.ID.String(), response.GetGroup().GetId())
	assert.Equal(t, groupToConvert.SpecializationCode, response.GetGroup().GetSpecializationCode())
	assert.Equal(t, groupToConvert.GroupNumber, response.GetGroup().GetGroupNumber())
}

func TestConvertToGroupFilter(t *testing.T) {
	request := client.GroupSearchByFilterRequest{
		SpecializationCode: gofakeit.WeekDay(),
		GroupNumber:        gofakeit.WeekDay(),
	}

	groupFilter := group.ConvertToGroupFilter(&request)

	assert.Equal(t, request.GetSpecializationCode(), groupFilter.SpecializationCode)
	assert.Equal(t, request.GetGroupNumber(), groupFilter.GroupNumber)
}

func TestConvertToSearchByFilterResponse(t *testing.T) {
	matchedGroups := []business.Group{randomGroup(), randomGroup(), randomGroup()}
	response := group.ConvertToSearchByFilterResponse(matchedGroups)

	require.Equal(t, len(matchedGroups), len(response.GetMatchedGroups()))
	for idx := range matchedGroups {
		assertGroupsEqual(t, matchedGroups[idx], response.GetMatchedGroups()[idx])
	}
}

func assertStudentsEqual(t *testing.T, left business.Student, right *client.StudentDTO) {
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

func assertGroupsEqual(t *testing.T, left business.Group, right *client.GroupDTO) {
	assert.Equal(t, left.ID.String(), right.GetId())
	assert.Equal(t, left.SpecializationCode, right.GetSpecializationCode())
	assert.Equal(t, left.GroupNumber, right.GetGroupNumber())
}

func randomStudent() business.Student {
	return business.Student{
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
}

func randomGroup() business.Group {
	return business.Group{
		ID:                 uuid.New(),
		SpecializationCode: gofakeit.WeekDay(),
		GroupNumber:        gofakeit.WeekDay(),
	}
}
