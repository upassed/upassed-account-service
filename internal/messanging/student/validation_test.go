package student_test

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
	"github.com/upassed/upassed-account-service/internal/util"
	"testing"
)

func TestStudentCreateRequestEmailValidation_InvalidEmail(t *testing.T) {
	request := util.RandomEventStudentCreateRequest()
	request.EducationalEmail = "invalid-email"

	err := request.Validate()
	require.NotNil(t, err)
}

func TestStudentCreateRequestEmailValidation_EmptyUsername(t *testing.T) {
	request := util.RandomEventStudentCreateRequest()
	request.Username = ""

	err := request.Validate()
	require.NotNil(t, err)
}

func TestStudentCreateRequestEmailValidation_TooLongUsername(t *testing.T) {
	request := util.RandomEventStudentCreateRequest()
	request.Username = gofakeit.LoremIpsumSentence(50)

	err := request.Validate()
	require.NotNil(t, err)
}

func TestStudentCreateRequestEmailValidation_TooShortUsername(t *testing.T) {
	request := util.RandomEventStudentCreateRequest()
	request.Username = "1"

	err := request.Validate()
	require.NotNil(t, err)
}

func TestStudentCreateRequestEmailValidation_InvalidGroupID(t *testing.T) {
	request := util.RandomEventStudentCreateRequest()
	request.GroupId = "invalid-uuid"

	err := request.Validate()
	require.NotNil(t, err)
}
