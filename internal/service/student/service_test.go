package student_test

import (
	"context"
	"errors"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/upassed/upassed-account-service/internal/util"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/upassed/upassed-account-service/internal/config"
	"github.com/upassed/upassed-account-service/internal/logging"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	"github.com/upassed/upassed-account-service/internal/service/student"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mockStudentRepository struct {
	mock.Mock
}

func (m *mockStudentRepository) Save(ctx context.Context, student *domain.Student) error {
	args := m.Called(ctx, student)
	return args.Error(0)
}

func (m *mockStudentRepository) FindByUsername(ctx context.Context, studentUsername string) (*domain.Student, error) {
	args := m.Called(ctx, studentUsername)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.Student), args.Error(1)
}

func (m *mockStudentRepository) CheckDuplicateExists(ctx context.Context, educationalEmail, username string) (bool, error) {
	args := m.Called(ctx, educationalEmail, username)
	return args.Bool(0), args.Error(1)
}

type mockGroupRepository struct {
	mock.Mock
}

func (m *mockGroupRepository) Exists(ctx context.Context, groupID uuid.UUID) (bool, error) {
	args := m.Called(ctx, groupID)
	return args.Bool(0), args.Error(1)
}

func (m *mockGroupRepository) FindByID(ctx context.Context, groupID uuid.UUID) (*domain.Group, error) {
	args := m.Called(ctx, groupID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.Group), args.Error(1)
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

func TestCreate_ErrorCheckingDuplicateExists(t *testing.T) {
	studentRepository := new(mockStudentRepository)
	groupRepository := new(mockGroupRepository)

	studentToCreate := util.RandomBusinessStudent()
	expectedRepositoryError := errors.New("some repo error")
	studentRepository.On(
		"CheckDuplicateExists",
		mock.Anything,
		studentToCreate.EducationalEmail,
		studentToCreate.Username,
	).Return(false, expectedRepositoryError)

	service := student.New(cfg, logging.New(config.EnvTesting), studentRepository, groupRepository)
	_, err := service.Create(context.Background(), studentToCreate)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedRepositoryError.Error(), convertedError.Message())
	assert.Equal(t, codes.Internal, convertedError.Code())
}

func TestCreate_DuplicateExists(t *testing.T) {
	studentRepository := new(mockStudentRepository)
	groupRepository := new(mockGroupRepository)

	studentToCreate := util.RandomBusinessStudent()
	studentRepository.On(
		"CheckDuplicateExists",
		mock.Anything,
		studentToCreate.EducationalEmail,
		studentToCreate.Username,
	).Return(true, nil)

	service := student.New(cfg, logging.New(config.EnvTesting), studentRepository, groupRepository)
	_, err := service.Create(context.Background(), studentToCreate)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, "student duplicate found", convertedError.Message())
	assert.Equal(t, codes.AlreadyExists, convertedError.Code())
}

func TestCreate_ErrorCheckingGroupExists(t *testing.T) {
	studentRepository := new(mockStudentRepository)
	groupRepository := new(mockGroupRepository)

	studentToCreate := util.RandomBusinessStudent()
	studentRepository.On(
		"CheckDuplicateExists",
		mock.Anything,
		studentToCreate.EducationalEmail,
		studentToCreate.Username,
	).Return(false, nil)

	expectedRepositoryError := errors.New("some repo error")
	groupRepository.On("Exists", mock.Anything, studentToCreate.Group.ID).Return(false, expectedRepositoryError)

	service := student.New(cfg, logging.New(config.EnvTesting), studentRepository, groupRepository)
	_, err := service.Create(context.Background(), studentToCreate)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedRepositoryError.Error(), convertedError.Message())
	assert.Equal(t, codes.Internal, convertedError.Code())
}

func TestCreate_GroupNotExists(t *testing.T) {
	studentRepository := new(mockStudentRepository)
	groupRepository := new(mockGroupRepository)

	studentToCreate := util.RandomBusinessStudent()
	studentRepository.On(
		"CheckDuplicateExists",
		mock.Anything,
		studentToCreate.EducationalEmail,
		studentToCreate.Username,
	).Return(false, nil)

	groupRepository.On("Exists", mock.Anything, studentToCreate.Group.ID).Return(false, nil)

	service := student.New(cfg, logging.New(config.EnvTesting), studentRepository, groupRepository)
	_, err := service.Create(context.Background(), studentToCreate)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, "group does not exists by id", convertedError.Message())
	assert.Equal(t, codes.NotFound, convertedError.Code())
}

func TestCreate_ErrorSavingToDatabase(t *testing.T) {
	studentRepository := new(mockStudentRepository)
	groupRepository := new(mockGroupRepository)

	studentToCreate := util.RandomBusinessStudent()
	studentRepository.On(
		"CheckDuplicateExists",
		mock.Anything,
		studentToCreate.EducationalEmail,
		studentToCreate.Username,
	).Return(false, nil)

	groupRepository.On("Exists", mock.Anything, studentToCreate.Group.ID).Return(true, nil)
	groupRepository.On("FindByID", mock.Anything, studentToCreate.Group.ID).Return(util.RandomDomainGroup(), nil)

	expectedRepositoryError := errors.New("some repo error")
	studentRepository.On("Save", mock.Anything, mock.Anything).Return(expectedRepositoryError)

	service := student.New(cfg, logging.New(config.EnvTesting), studentRepository, groupRepository)
	_, err := service.Create(context.Background(), studentToCreate)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedRepositoryError.Error(), convertedError.Message())
	assert.Equal(t, codes.Internal, convertedError.Code())
}

func TestCreate_HappyPath(t *testing.T) {
	studentRepository := new(mockStudentRepository)
	groupRepository := new(mockGroupRepository)

	studentToCreate := util.RandomBusinessStudent()
	studentRepository.On(
		"CheckDuplicateExists",
		mock.Anything,
		studentToCreate.EducationalEmail,
		studentToCreate.Username,
	).Return(false, nil)

	groupRepository.On("Exists", mock.Anything, studentToCreate.Group.ID).Return(true, nil)
	groupRepository.On("FindByID", mock.Anything, studentToCreate.Group.ID).Return(util.RandomDomainGroup(), nil)
	studentRepository.On("Save", mock.Anything, mock.Anything).Return(nil)

	service := student.New(cfg, logging.New(config.EnvTesting), studentRepository, groupRepository)
	response, err := service.Create(context.Background(), studentToCreate)
	require.NoError(t, err)

	assert.Equal(t, studentToCreate.ID, response.CreatedStudentID)
}

func TestFindByUsername_ErrorSearchingStudentInDatabase(t *testing.T) {
	studentRepository := new(mockStudentRepository)
	studentUsername := gofakeit.Username()

	expectedRepositoryError := errors.New("some repo error")
	studentRepository.On("FindByUsername", mock.Anything, studentUsername).Return(nil, expectedRepositoryError)

	service := student.New(cfg, logging.New(config.EnvTesting), studentRepository, new(mockGroupRepository))
	_, err := service.FindByUsername(context.Background(), studentUsername)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedRepositoryError.Error(), convertedError.Message())
	assert.Equal(t, codes.Internal, convertedError.Code())
}

func TestFindByUsername_DeadlineExceeded(t *testing.T) {
	oldTimeout := cfg.Timeouts.EndpointExecutionTimeoutMS
	cfg.Timeouts.EndpointExecutionTimeoutMS = "0"

	studentRepository := new(mockStudentRepository)
	studentUsername := gofakeit.Username()
	foundStudent := util.RandomDomainStudent()

	studentRepository.On("FindByUsername", mock.Anything, studentUsername).Return(foundStudent, nil)

	studentService := student.New(cfg, logging.New(config.EnvTesting), studentRepository, new(mockGroupRepository))
	_, err := studentService.FindByUsername(context.Background(), studentUsername)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, codes.DeadlineExceeded, convertedError.Code())

	cfg.Timeouts.EndpointExecutionTimeoutMS = oldTimeout
}

func TestFindByUsername_HappyPath(t *testing.T) {
	studentRepository := new(mockStudentRepository)
	studentUsername := gofakeit.Username()
	foundStudent := util.RandomDomainStudent()

	studentRepository.On("FindByUsername", mock.Anything, studentUsername).Return(foundStudent, nil)

	studentService := student.New(cfg, logging.New(config.EnvTesting), studentRepository, new(mockGroupRepository))
	response, err := studentService.FindByUsername(context.Background(), studentUsername)
	require.NoError(t, err)

	assert.Equal(t, foundStudent, student.ConvertToRepositoryStudent(response))
}
