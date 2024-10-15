package handling_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/upassed/upassed-account-service/internal/handling"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestConvertApplicationError_AddDetails(t *testing.T) {
	message := gofakeit.Error().Error()
	code := codes.Internal
	applicationError := handling.NewApplicationError(message, code)

	handledError := handling.HandleApplicationError(applicationError)

	st := status.Convert(handledError)
	assert.Equal(t, code, st.Code())
	assert.Equal(t, message, st.Message())
	assert.Equal(t, 1, len(st.Details()))
}

func TestConvertApplicationError_WrapOptions(t *testing.T) {
	message := "error message"
	code := codes.AlreadyExists

	applicationError := handling.NewApplicationError(message, code)
	wrappedError := handling.HandleApplicationError(applicationError, handling.WithCode(codes.OK))

	st := status.Convert(wrappedError)
	require.NotNil(t, wrappedError)

	assert.Equal(t, message, st.Message())
	assert.Equal(t, code, st.Code())
}

func TestConvertApplicationError_WrappingNotAnApplicationError(t *testing.T) {
	initialError := gofakeit.Error()

	handledError := handling.HandleApplicationError(initialError)

	st := status.Convert(handledError)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Equal(t, initialError.Error(), st.Message())
	assert.Equal(t, 1, len(st.Details()))
}

func TestCreateAnApplicationError(t *testing.T) {
	message := "error message"
	code := codes.AlreadyExists

	applicationError := handling.NewApplicationError(message, code)

	require.NotNil(t, applicationError)

	assert.Equal(t, message, applicationError.Error())
	assert.Equal(t, message, applicationError.GRPCStatus().Message())
	assert.Equal(t, code, applicationError.GRPCStatus().Code())
}
