package student_test

import (
	"context"
	"errors"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/golang/mock/gomock"
	"github.com/upassed/upassed-account-service/internal/util"
	"github.com/upassed/upassed-account-service/internal/util/mocks"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/upassed/upassed-account-service/internal/config"
	"github.com/upassed/upassed-account-service/internal/logging"
	"github.com/upassed/upassed-account-service/internal/service/student"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	cfg               *config.Config
	service           student.Service
	studentRepository *mocks.StudentRepository
	groupRepository   *mocks.GroupRepository
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

	studentRepository = mocks.NewStudentRepository(ctrl)
	groupRepository = mocks.NewGroupRepository(ctrl)
	service = student.New(cfg, logging.New(config.EnvTesting), studentRepository, groupRepository)

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestCreate_ErrorCheckingDuplicateExists(t *testing.T) {
	studentToCreate := util.RandomBusinessStudent()
	expectedRepositoryError := errors.New("some repo error")

	studentRepository.EXPECT().
		CheckDuplicateExists(gomock.Any(), studentToCreate.EducationalEmail, studentToCreate.Username).
		Return(false, expectedRepositoryError)

	_, err := service.Create(context.Background(), studentToCreate)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedRepositoryError.Error(), convertedError.Message())
	assert.Equal(t, codes.Internal, convertedError.Code())
}

func TestCreate_DuplicateExists(t *testing.T) {
	studentToCreate := util.RandomBusinessStudent()

	studentRepository.EXPECT().
		CheckDuplicateExists(gomock.Any(), studentToCreate.EducationalEmail, studentToCreate.Username).
		Return(true, nil)

	_, err := service.Create(context.Background(), studentToCreate)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, "student duplicate found", convertedError.Message())
	assert.Equal(t, codes.AlreadyExists, convertedError.Code())
}

func TestCreate_ErrorCheckingGroupExists(t *testing.T) {
	studentToCreate := util.RandomBusinessStudent()
	studentRepository.EXPECT().
		CheckDuplicateExists(gomock.Any(), studentToCreate.EducationalEmail, studentToCreate.Username).
		Return(false, nil)

	expectedRepositoryError := errors.New("some repo error")
	groupRepository.EXPECT().
		Exists(gomock.Any(), studentToCreate.Group.ID).
		Return(false, expectedRepositoryError)

	_, err := service.Create(context.Background(), studentToCreate)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedRepositoryError.Error(), convertedError.Message())
	assert.Equal(t, codes.Internal, convertedError.Code())
}

func TestCreate_GroupNotExists(t *testing.T) {
	studentToCreate := util.RandomBusinessStudent()
	studentRepository.EXPECT().
		CheckDuplicateExists(gomock.Any(), studentToCreate.EducationalEmail, studentToCreate.Username).
		Return(false, nil)

	groupRepository.EXPECT().
		Exists(gomock.Any(), studentToCreate.Group.ID).
		Return(false, nil)

	_, err := service.Create(context.Background(), studentToCreate)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, "group does not exists by id", convertedError.Message())
	assert.Equal(t, codes.NotFound, convertedError.Code())
}

func TestCreate_ErrorSavingToDatabase(t *testing.T) {
	studentToCreate := util.RandomBusinessStudent()
	studentRepository.EXPECT().
		CheckDuplicateExists(gomock.Any(), studentToCreate.EducationalEmail, studentToCreate.Username).
		Return(false, nil)

	groupRepository.EXPECT().
		Exists(gomock.Any(), studentToCreate.Group.ID).
		Return(true, nil)

	groupRepository.EXPECT().
		FindByID(gomock.Any(), studentToCreate.Group.ID).
		Return(util.RandomDomainGroup(), nil)

	expectedRepositoryError := errors.New("some repo error")
	studentRepository.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Return(expectedRepositoryError)

	_, err := service.Create(context.Background(), studentToCreate)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedRepositoryError.Error(), convertedError.Message())
	assert.Equal(t, codes.Internal, convertedError.Code())
}

func TestCreate_HappyPath(t *testing.T) {
	studentToCreate := util.RandomBusinessStudent()
	studentRepository.EXPECT().
		CheckDuplicateExists(gomock.Any(), studentToCreate.EducationalEmail, studentToCreate.Username).
		Return(false, nil)

	groupRepository.EXPECT().
		Exists(gomock.Any(), studentToCreate.Group.ID).
		Return(true, nil)

	groupRepository.EXPECT().
		FindByID(gomock.Any(), studentToCreate.Group.ID).
		Return(util.RandomDomainGroup(), nil)

	studentRepository.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Return(nil)

	response, err := service.Create(context.Background(), studentToCreate)
	require.NoError(t, err)

	assert.Equal(t, studentToCreate.ID, response.CreatedStudentID)
}

func TestFindByUsername_ErrorSearchingStudentInDatabase(t *testing.T) {
	studentUsername := gofakeit.Username()

	expectedRepositoryError := errors.New("some repo error")
	studentRepository.EXPECT().
		FindByUsername(gomock.Any(), studentUsername).
		Return(nil, expectedRepositoryError)

	_, err := service.FindByUsername(context.Background(), studentUsername)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedRepositoryError.Error(), convertedError.Message())
	assert.Equal(t, codes.Internal, convertedError.Code())
}

func TestFindByUsername_DeadlineExceeded(t *testing.T) {
	oldTimeout := cfg.Timeouts.EndpointExecutionTimeoutMS
	cfg.Timeouts.EndpointExecutionTimeoutMS = "0"

	studentUsername := gofakeit.Username()
	foundStudent := util.RandomDomainStudent()

	studentRepository.EXPECT().
		FindByUsername(gomock.Any(), studentUsername).
		Return(foundStudent, nil)

	_, err := service.FindByUsername(context.Background(), studentUsername)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, codes.DeadlineExceeded, convertedError.Code())

	cfg.Timeouts.EndpointExecutionTimeoutMS = oldTimeout
}

func TestFindByUsername_HappyPath(t *testing.T) {
	studentUsername := gofakeit.Username()
	foundStudent := util.RandomDomainStudent()

	studentRepository.EXPECT().
		FindByUsername(gomock.Any(), studentUsername).
		Return(foundStudent, nil)

	response, err := service.FindByUsername(context.Background(), studentUsername)
	require.NoError(t, err)

	assert.Equal(t, foundStudent, student.ConvertToRepositoryStudent(response))
}
