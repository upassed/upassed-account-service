package group_test

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/upassed/upassed-account-service/internal/config"
	"github.com/upassed/upassed-account-service/internal/logger"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	"github.com/upassed/upassed-account-service/internal/service/group"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mockGroupRepository struct {
	mock.Mock
}

func (m *mockGroupRepository) FindStudentsInGroup(ctx context.Context, groupID uuid.UUID) ([]domain.Student, error) {
	args := m.Called(ctx, groupID)
	return args.Get(0).([]domain.Student), args.Error(1)
}

func TestFindStudentsInGroup_ErrorInRepositoryLayer(t *testing.T) {
	groupRepository := new(mockGroupRepository)

	groupID := uuid.New()
	expectedReposotiryError := errors.New("some repo error")
	groupRepository.On("FindStudentsInGroup", mock.Anything, groupID).Return([]domain.Student{}, expectedReposotiryError)

	service := group.New(logger.New(config.EnvTesting), groupRepository)
	_, err := service.FindStudentsInGroup(context.Background(), groupID)
	require.NotNil(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedReposotiryError.Error(), convertedError.Message())
	assert.Equal(t, codes.Internal, convertedError.Code())
}

func TestFindStudentsInGroup_HappyPath(t *testing.T) {
	groupRepository := new(mockGroupRepository)

	groupID := uuid.New()
	expectedStudentsInGroup := []domain.Student{randomStudent(), randomStudent(), randomStudent()}
	groupRepository.On("FindStudentsInGroup", mock.Anything, groupID).Return(expectedStudentsInGroup, nil)

	service := group.New(logger.New(config.EnvTesting), groupRepository)
	actualFoundStudentsInGroup, err := service.FindStudentsInGroup(context.Background(), groupID)
	require.Nil(t, err)

	assert.Equal(t, len(expectedStudentsInGroup), len(actualFoundStudentsInGroup))
}

func randomStudent() domain.Student {
	return domain.Student{
		ID:               uuid.New(),
		FirstName:        gofakeit.FirstName(),
		LastName:         gofakeit.LastName(),
		MiddleName:       gofakeit.MiddleName(),
		EducationalEmail: gofakeit.Email(),
		Username:         gofakeit.Username(),
		Group: domain.Group{
			ID:                 uuid.New(),
			SpecializationCode: gofakeit.WeekDay(),
			GroupNumber:        gofakeit.WeekDay(),
		},
	}
}
