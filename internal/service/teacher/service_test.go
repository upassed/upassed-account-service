package teacher_test

import (
	"context"
	"errors"
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
	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/internal/logging"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	business "github.com/upassed/upassed-account-service/internal/service/model"
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
	projectRoot, err := getProjectRoot()
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
	duplicateTeacher := randomTeacher()

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
	duplicateTeacher := randomTeacher()

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
	teacherToSave := randomTeacher()

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
	teacherToSave := randomTeacher()

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
	expectedFoundTeacher := teacher.ConvertToRepositoryTeacher(randomTeacher())

	repository.On("FindByID", mock.Anything, teacherID).Return(expectedFoundTeacher, nil)
	service := teacher.New(cfg, logger, repository)

	foundTeacher, err := service.FindByID(context.Background(), teacherID)
	require.Nil(t, err)

	assert.Equal(t, teacher.ConvertToServiceTeacher(expectedFoundTeacher), foundTeacher)
}

func randomTeacher() business.Teacher {
	return business.Teacher{
		ID:          uuid.New(),
		FirstName:   gofakeit.FirstName(),
		LastName:    gofakeit.LastName(),
		MiddleName:  gofakeit.MiddleName(),
		ReportEmail: gofakeit.Email(),
		Username:    gofakeit.Username(),
	}
}

func getProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			return "", errors.New("project root not found")
		}

		dir = parentDir
	}
}
