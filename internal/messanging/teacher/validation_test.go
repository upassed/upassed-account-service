package teacher_test

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
	"github.com/upassed/upassed-account-service/internal/util"
	"testing"
)

func TestTeacherCreateRequestEmailValidation_InvalidEmail(t *testing.T) {
	request := util.RandomEventTeacherCreateRequest()
	request.ReportEmail = "invalid-email"

	err := request.Validate()
	require.Error(t, err)
}

func TestTeacherCreateRequestEmailValidation_EmptyUsername(t *testing.T) {
	request := util.RandomEventTeacherCreateRequest()
	request.Username = ""

	err := request.Validate()
	require.Error(t, err)
}

func TestTeacherCreateRequestEmailValidation_TooLongUsername(t *testing.T) {
	request := util.RandomEventTeacherCreateRequest()
	request.Username = gofakeit.LoremIpsumSentence(50)

	err := request.Validate()
	require.Error(t, err)
}

func TestTeacherCreateRequestEmailValidation_TooShortUsername(t *testing.T) {
	request := util.RandomEventTeacherCreateRequest()
	request.Username = "1"

	err := request.Validate()
	require.Error(t, err)
}
