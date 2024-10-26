package group_test

import (
	"context"
	"errors"
	"github.com/upassed/upassed-account-service/internal/util"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/upassed/upassed-account-service/internal/config"
	"github.com/upassed/upassed-account-service/internal/logging"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	"github.com/upassed/upassed-account-service/internal/service/group"
	business "github.com/upassed/upassed-account-service/internal/service/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mockGroupRepository struct {
	mock.Mock
}

func (m *mockGroupRepository) FindStudentsInGroup(ctx context.Context, groupID uuid.UUID) ([]*domain.Student, error) {
	args := m.Called(ctx, groupID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*domain.Student), args.Error(1)
}

func (m *mockGroupRepository) FindByID(ctx context.Context, groupID uuid.UUID) (*domain.Group, error) {
	args := m.Called(ctx, groupID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.Group), args.Error(1)
}

func (m *mockGroupRepository) FindByFilter(ctx context.Context, filter *domain.GroupFilter) ([]*domain.Group, error) {
	args := m.Called(ctx, filter)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*domain.Group), args.Error(1)
}

var (
	cfg *config.Config
)

func TestMain(m *testing.M) {
	currentDir, _ := os.Getwd()
	projectRoot, err := util.GetProjectRoot(currentDir)
	if err != nil {
		log.Fatal("error to get project root folder: ", err)
	}

	if err := os.Setenv(config.EnvConfigPath, filepath.Join(projectRoot, "config", "test.yml")); err != nil {
		log.Fatal(err)
	}

	cfg, err = config.Load()
	if err != nil {
		log.Fatal("unable to parse config: ", err)
	}

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestFindStudentsInGroup_ErrorInRepositoryLayer(t *testing.T) {
	groupRepository := new(mockGroupRepository)

	groupID := uuid.New()
	expectedRepositoryError := errors.New("some repo error")
	groupRepository.On("FindStudentsInGroup", mock.Anything, groupID).Return(nil, expectedRepositoryError)

	service := group.New(cfg, logging.New(config.EnvTesting), groupRepository)
	_, err := service.FindStudentsInGroup(context.Background(), groupID)
	require.NotNil(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedRepositoryError.Error(), convertedError.Message())
	assert.Equal(t, codes.Internal, convertedError.Code())
}

func TestFindStudentsInGroup_HappyPath(t *testing.T) {
	groupRepository := new(mockGroupRepository)

	groupID := uuid.New()
	expectedStudentsInGroup := []*domain.Student{
		util.RandomDomainStudent(),
		util.RandomDomainStudent(),
		util.RandomDomainStudent(),
	}

	groupRepository.On("FindStudentsInGroup", mock.Anything, groupID).Return(expectedStudentsInGroup, nil)

	service := group.New(cfg, logging.New(config.EnvTesting), groupRepository)
	actualFoundStudentsInGroup, err := service.FindStudentsInGroup(context.Background(), groupID)
	require.Nil(t, err)

	assert.Equal(t, len(expectedStudentsInGroup), len(actualFoundStudentsInGroup))
}

func TestFindByID_RepositoryError(t *testing.T) {
	groupRepository := new(mockGroupRepository)

	groupID := uuid.New()
	expectedRepositoryError := errors.New("some repo error")
	groupRepository.On("FindByID", mock.Anything, groupID).Return(nil, expectedRepositoryError)

	service := group.New(cfg, logging.New(config.EnvTesting), groupRepository)
	_, err := service.FindByID(context.Background(), groupID)
	require.NotNil(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedRepositoryError.Error(), convertedError.Message())
	assert.Equal(t, codes.Internal, convertedError.Code())
}

func TestFindByID_HappyPath(t *testing.T) {
	groupRepository := new(mockGroupRepository)

	groupID := uuid.New()
	expectedFoundGroup := util.RandomDomainGroup()
	expectedFoundGroup.ID = groupID

	groupRepository.On("FindByID", mock.Anything, groupID).Return(expectedFoundGroup, nil)

	service := group.New(cfg, logging.New(config.EnvTesting), groupRepository)
	foundGroup, err := service.FindByID(context.Background(), groupID)
	require.Nil(t, err)

	assert.Equal(t, expectedFoundGroup.ID, foundGroup.ID)
	assert.Equal(t, expectedFoundGroup.SpecializationCode, foundGroup.SpecializationCode)
	assert.Equal(t, expectedFoundGroup.GroupNumber, foundGroup.GroupNumber)
}

func TestFindByFilter_RepositoryError(t *testing.T) {
	groupRepository := new(mockGroupRepository)

	groupFilter := &business.GroupFilter{
		SpecializationCode: gofakeit.WeekDay(),
		GroupNumber:        gofakeit.WeekDay(),
	}

	expectedRepositoryError := errors.New("some repo error")
	groupRepository.On("FindByFilter", mock.Anything, mock.Anything).Return(nil, expectedRepositoryError)

	service := group.New(cfg, logging.New(config.EnvTesting), groupRepository)
	_, err := service.FindByFilter(context.Background(), groupFilter)
	require.NotNil(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedRepositoryError.Error(), convertedError.Message())
	assert.Equal(t, codes.Internal, convertedError.Code())
}

func TestFindByFilter_HappyPath(t *testing.T) {
	groupRepository := new(mockGroupRepository)

	groupFilter := &business.GroupFilter{
		SpecializationCode: gofakeit.WeekDay(),
		GroupNumber:        gofakeit.WeekDay(),
	}

	foundMatchedGroups := []*domain.Group{util.RandomDomainGroup(), util.RandomDomainGroup(), util.RandomDomainGroup()}
	groupRepository.On("FindByFilter", mock.Anything, mock.Anything).Return(foundMatchedGroups, nil)

	service := group.New(cfg, logging.New(config.EnvTesting), groupRepository)
	response, err := service.FindByFilter(context.Background(), groupFilter)
	require.Nil(t, err)

	assert.Equal(t, len(foundMatchedGroups), len(response))
}
