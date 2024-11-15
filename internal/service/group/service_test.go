package group_test

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/upassed/upassed-account-service/internal/util"
	"github.com/upassed/upassed-account-service/internal/util/mocks"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/upassed/upassed-account-service/internal/config"
	"github.com/upassed/upassed-account-service/internal/logging"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	"github.com/upassed/upassed-account-service/internal/service/group"
	business "github.com/upassed/upassed-account-service/internal/service/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	cfg        *config.Config
	service    group.Service
	repository *mocks.GroupRepository
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

	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()

	repository = mocks.NewGroupRepository(ctrl)
	service = group.New(cfg, logging.New(config.EnvTesting), repository)

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestFindStudentsInGroup_ErrorInRepositoryLayer(t *testing.T) {
	groupID := uuid.New()
	expectedRepositoryError := errors.New("some repo error")
	repository.EXPECT().
		FindStudentsInGroup(gomock.Any(), groupID).
		Return(nil, expectedRepositoryError)

	_, err := service.FindStudentsInGroup(context.Background(), groupID)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedRepositoryError.Error(), convertedError.Message())
	assert.Equal(t, codes.Internal, convertedError.Code())
}

func TestFindStudentsInGroup_HappyPath(t *testing.T) {
	groupID := uuid.New()
	expectedStudentsInGroup := []*domain.Student{
		util.RandomDomainStudent(),
		util.RandomDomainStudent(),
		util.RandomDomainStudent(),
	}

	repository.EXPECT().
		FindStudentsInGroup(gomock.Any(), groupID).
		Return(expectedStudentsInGroup, nil)

	actualFoundStudentsInGroup, err := service.FindStudentsInGroup(context.Background(), groupID)
	require.NoError(t, err)

	assert.Equal(t, len(expectedStudentsInGroup), len(actualFoundStudentsInGroup))
}

func TestFindByID_RepositoryError(t *testing.T) {
	groupID := uuid.New()
	expectedRepositoryError := errors.New("some repo error")

	repository.EXPECT().
		FindByID(gomock.Any(), groupID).
		Return(nil, expectedRepositoryError)

	_, err := service.FindByID(context.Background(), groupID)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedRepositoryError.Error(), convertedError.Message())
	assert.Equal(t, codes.Internal, convertedError.Code())
}

func TestFindByID_HappyPath(t *testing.T) {
	groupID := uuid.New()
	expectedFoundGroup := util.RandomDomainGroup()
	expectedFoundGroup.ID = groupID

	repository.EXPECT().
		FindByID(gomock.Any(), groupID).
		Return(expectedFoundGroup, nil)

	foundGroup, err := service.FindByID(context.Background(), groupID)
	require.NoError(t, err)

	assert.Equal(t, expectedFoundGroup.ID, foundGroup.ID)
	assert.Equal(t, expectedFoundGroup.SpecializationCode, foundGroup.SpecializationCode)
	assert.Equal(t, expectedFoundGroup.GroupNumber, foundGroup.GroupNumber)
}

func TestFindByFilter_RepositoryError(t *testing.T) {
	groupFilter := &business.GroupFilter{
		SpecializationCode: gofakeit.WeekDay(),
		GroupNumber:        gofakeit.WeekDay(),
	}

	expectedRepositoryError := errors.New("some repo error")
	repository.EXPECT().
		FindByFilter(gomock.Any(), gomock.Any()).
		Return(nil, expectedRepositoryError)

	_, err := service.FindByFilter(context.Background(), groupFilter)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedRepositoryError.Error(), convertedError.Message())
	assert.Equal(t, codes.Internal, convertedError.Code())
}

func TestFindByFilter_HappyPath(t *testing.T) {
	groupFilter := &business.GroupFilter{
		SpecializationCode: gofakeit.WeekDay(),
		GroupNumber:        gofakeit.WeekDay(),
	}

	foundMatchedGroups := []*domain.Group{util.RandomDomainGroup(), util.RandomDomainGroup(), util.RandomDomainGroup()}
	repository.EXPECT().
		FindByFilter(gomock.Any(), gomock.Any()).
		Return(foundMatchedGroups, nil)

	response, err := service.FindByFilter(context.Background(), groupFilter)
	require.NoError(t, err)

	assert.Equal(t, len(foundMatchedGroups), len(response))
}
