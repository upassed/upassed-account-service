package group_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	"github.com/upassed/upassed-account-service/internal/service/group"
	business "github.com/upassed/upassed-account-service/internal/service/model"
)

func TestConvertToServiceStudents(t *testing.T) {
	studentsToConvert := []domain.Student{randomStudent(), randomStudent(), randomStudent()}
	convertedStudents := group.ConvertToServiceStudents(studentsToConvert)

	require.Equal(t, len(studentsToConvert), len(convertedStudents))
	for idx := range studentsToConvert {
		assertStudentsEqual(t, studentsToConvert[idx], convertedStudents[idx])
	}
}

func TestConvertToServiceGroup(t *testing.T) {
	groupToConvert := domain.Group{
		ID:                 uuid.New(),
		SpecializationCode: gofakeit.WeekDay(),
		GroupNumber:        gofakeit.WeekDay(),
	}

	convertedGroup := group.ConvertToServiceGroup(groupToConvert)

	assert.Equal(t, groupToConvert.ID, convertedGroup.ID)
	assert.Equal(t, groupToConvert.SpecializationCode, convertedGroup.SpecializationCode)
	assert.Equal(t, groupToConvert.GroupNumber, convertedGroup.GroupNumber)
}

func TestConvertToGroupFilter(t *testing.T) {
	filterToConvert := business.GroupFilter{
		SpecializationCode: gofakeit.WeekDay(),
		GroupNumber:        gofakeit.WeekDay(),
	}

	convertedFilter := group.ConvertToGroupFilter(filterToConvert)

	assert.Equal(t, filterToConvert.SpecializationCode, convertedFilter.SpecializationCode)
	assert.Equal(t, filterToConvert.GroupNumber, convertedFilter.GroupNumber)
}

func assertStudentsEqual(t *testing.T, left domain.Student, right business.Student) {
	assert.Equal(t, left.ID, right.ID)
	assert.Equal(t, left.FirstName, right.FirstName)
	assert.Equal(t, left.LastName, right.LastName)
	assert.Equal(t, left.MiddleName, right.MiddleName)
	assert.Equal(t, left.EducationalEmail, right.EducationalEmail)
	assert.Equal(t, left.Username, right.Username)
	assert.Equal(t, left.Group.ID, right.Group.ID)
	assert.Equal(t, left.Group.SpecializationCode, right.Group.SpecializationCode)
	assert.Equal(t, left.Group.GroupNumber, right.Group.GroupNumber)
}
