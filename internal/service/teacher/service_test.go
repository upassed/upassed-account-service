package teacher_test

import (
	"context"
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
	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/internal/logging"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	"github.com/upassed/upassed-account-service/internal/service/teacher"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mockTeacherRepository struct {
	mock.Mock
}

func (m *mockTeacherRepository) Save(ctx context.Context, teacher domain.Teacher) error {
	args := m.Called(ctx, teacher)
	return args.Error(0)
}

func (m *mockTeacherRepository) FindByID(ctx context.Context, teacherID uuid.UUID) (domain.Teacher, error) {
	args := m.Called(ctx, teacherID)
	return args.Get(0).(domain.Teacher), args.Error(1)
}

func (m *mockTeacherRepository) CheckDuplicateExists(ctx context.Context, reportEmail, username string) (bool, error) {
	args := m.Called(ctx, reportEmail, username)
	return args.Bool(0), args.Error(1)
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

func TestCreate_ErrorCheckingDuplicateExistsOccurred(t *testing.T) {
	logger := logging.New(config.EnvTesting)
	repository := new(mockTeacherRepository)
	duplicateTeacher := util.RandomBusinessTeacher()

	expectedRepoError := handling.New("repo layer error message", codes.Internal)
	repository.On("CheckDuplicateExists", mock.Anything, duplicateTeacher.ReportEmail, duplicateTeacher.Username).Return(false, expectedRepoError)

	service := teacher.New(cfg, logger, repository)

	_, err := service.Create(context.Background(), duplicateTeacher)
	require.NotNil(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedRepoError.Error(), convertedError.Message())
}

func TestCreate_DuplicateExists(t *testing.T) {
	logger := logging.New(config.EnvTesting)
	repository := new(mockTeacherRepository)
	duplicateTeacher := util.RandomBusinessTeacher()

	repository.On("CheckDuplicateExists", mock.Anything, duplicateTeacher.ReportEmail, duplicateTeacher.Username).Return(true, nil)

	service := teacher.New(cfg, logger, repository)

	_, err := service.Create(context.Background(), duplicateTeacher)
	require.NotNil(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, "teacher duplicate found", convertedError.Message())
	assert.Equal(t, codes.AlreadyExists, convertedError.Code())
}

func TestCreate_ErrorSavingToDatabase(t *testing.T) {
	logger := logging.New(config.EnvTesting)
	repository := new(mockTeacherRepository)
	teacherToSave := util.RandomBusinessTeacher()

	repository.On("CheckDuplicateExists", mock.Anything, teacherToSave.ReportEmail, teacherToSave.Username).Return(false, nil)

	expectedRepoError := handling.New("repo layer error message", codes.DeadlineExceeded)
	repository.On("Save", mock.Anything, mock.Anything).Return(expectedRepoError)

	service := teacher.New(cfg, logger, repository)

	_, err := service.Create(context.Background(), teacherToSave)
	require.NotNil(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedRepoError.Error(), convertedError.Message())
	assert.Equal(t, expectedRepoError.Code(), convertedError.Code())
}

func TestCreate_HappyPath(t *testing.T) {
	logger := logging.New(config.EnvTesting)
	repository := new(mockTeacherRepository)
	teacherToSave := util.RandomBusinessTeacher()

	repository.On("CheckDuplicateExists", mock.Anything, teacherToSave.ReportEmail, teacherToSave.Username).Return(false, nil)
	repository.On("Save", mock.Anything, mock.Anything).Return(nil)

	service := teacher.New(cfg, logger, repository)

	response, err := service.Create(context.Background(), teacherToSave)
	require.Nil(t, err)

	assert.Equal(t, teacherToSave.ID, response.CreatedTeacherID)
}

func TestFindByID_ErrorSearchingTeacherInDatabase(t *testing.T) {
	logger := logging.New(config.EnvTesting)
	teacherRepository := new(mockTeacherRepository)
	teacherID := uuid.New()

	expectedRepoError := handling.New("repo layer error message", codes.NotFound)
	teacherRepository.On("FindByID", mock.Anything, teacherID).Return(domain.Teacher{}, expectedRepoError)
	service := teacher.New(cfg, logger, teacherRepository)

	_, err := service.FindByID(context.Background(), teacherID)
	require.NotNil(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedRepoError.Code(), convertedError.Code())
	assert.Equal(t, expectedRepoError.Error(), convertedError.Message())
}

func TestFindByID_HappyPath(t *testing.T) {
	logger := logging.New(config.EnvTesting)
	repository := new(mockTeacherRepository)
	teacherID := uuid.New()
	expectedFoundTeacher := teacher.ConvertToRepositoryTeacher(util.RandomBusinessTeacher())

	repository.On("FindByID", mock.Anything, teacherID).Return(expectedFoundTeacher, nil)
	service := teacher.New(cfg, logger, repository)

	foundTeacher, err := service.FindByID(context.Background(), teacherID)
	require.Nil(t, err)

	assert.Equal(t, teacher.ConvertToServiceTeacher(expectedFoundTeacher), foundTeacher)
}
