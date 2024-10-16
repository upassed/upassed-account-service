package teacher_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	config "github.com/upassed/upassed-account-service/internal/config"
	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/internal/logger"
	repository "github.com/upassed/upassed-account-service/internal/repository/teacher"
	"github.com/upassed/upassed-account-service/internal/service/teacher"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mockTeacherRepository struct {
	mock.Mock
}

func (m *mockTeacherRepository) Save(ctx context.Context, teacher repository.Teacher) error {
	args := m.Called(ctx, teacher)
	return args.Error(0)
}

func (m *mockTeacherRepository) FindByID(ctx context.Context, teacherID uuid.UUID) (repository.Teacher, error) {
	args := m.Called(ctx, teacherID)
	return args.Get(0).(repository.Teacher), args.Error(1)
}

func (m *mockTeacherRepository) CheckDuplicateExists(ctx context.Context, reportEmail, username string) (bool, error) {
	args := m.Called(ctx, reportEmail, username)
	return args.Bool(0), args.Error(1)
}

func TestCreate_ErrorCheckingDuplicateExistsOccured(t *testing.T) {
	log := logger.New(config.EnvTesting)
	repository := new(mockTeacherRepository)
	duplicateTeacher := randomTeacher()

	expectedRepoError := handling.New("repo layer error message", codes.Internal)
	repository.On("CheckDuplicateExists", mock.Anything, duplicateTeacher.ReportEmail, duplicateTeacher.Username).Return(false, expectedRepoError)

	service := teacher.New(log, repository)

	_, err := service.Create(context.Background(), duplicateTeacher)
	require.NotNil(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedRepoError.Error(), convertedError.Message())
}

func TestCreate_DiplicateExists(t *testing.T) {
	log := logger.New(config.EnvTesting)
	repository := new(mockTeacherRepository)
	duplicateTeacher := randomTeacher()

	repository.On("CheckDuplicateExists", mock.Anything, duplicateTeacher.ReportEmail, duplicateTeacher.Username).Return(true, nil)

	service := teacher.New(log, repository)

	_, err := service.Create(context.Background(), duplicateTeacher)
	require.NotNil(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, "teacher duplicate found", convertedError.Message())
	assert.Equal(t, codes.AlreadyExists, convertedError.Code())
}

func TestCreate_ErrorSavingToDatabase(t *testing.T) {
	log := logger.New(config.EnvTesting)
	repository := new(mockTeacherRepository)
	teacherToSave := randomTeacher()

	repository.On("CheckDuplicateExists", mock.Anything, teacherToSave.ReportEmail, teacherToSave.Username).Return(false, nil)

	expectedRepoError := handling.New("repo layer error message", codes.DeadlineExceeded)
	repository.On("Save", mock.Anything, mock.Anything).Return(expectedRepoError)

	service := teacher.New(log, repository)

	_, err := service.Create(context.Background(), teacherToSave)
	require.NotNil(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedRepoError.Error(), convertedError.Message())
	assert.Equal(t, expectedRepoError.Code, convertedError.Code())
}

func TestCreate_HappyPath(t *testing.T) {
	log := logger.New(config.EnvTesting)
	repository := new(mockTeacherRepository)
	teacherToSave := randomTeacher()

	repository.On("CheckDuplicateExists", mock.Anything, teacherToSave.ReportEmail, teacherToSave.Username).Return(false, nil)
	repository.On("Save", mock.Anything, mock.Anything).Return(nil)

	service := teacher.New(log, repository)

	response, err := service.Create(context.Background(), teacherToSave)
	require.Nil(t, err)

	assert.Equal(t, teacherToSave.ID, response.CreatedTeacherID)
}

func TestFindByID_ErrorSearchingTeacherInDatabase(t *testing.T) {
	log := logger.New(config.EnvTesting)
	teacherRepository := new(mockTeacherRepository)
	teacherID := uuid.New()

	expectedRepoError := handling.New("repo layer error message", codes.NotFound)
	teacherRepository.On("FindByID", mock.Anything, teacherID).Return(repository.Teacher{}, expectedRepoError)
	service := teacher.New(log, teacherRepository)

	_, err := service.FindByID(context.Background(), teacherID)
	require.NotNil(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedRepoError.Code, convertedError.Code())
	assert.Equal(t, expectedRepoError.Message, convertedError.Message())
}

func TestFindByID_HappyPath(t *testing.T) {
	log := logger.New(config.EnvTesting)
	repository := new(mockTeacherRepository)
	teacherID := uuid.New()
	expectedFoundTeacher := teacher.ConvertToRepositoryTeacher(randomTeacher())

	repository.On("FindByID", mock.Anything, teacherID).Return(expectedFoundTeacher, nil)
	service := teacher.New(log, repository)

	foundTeacher, err := service.FindByID(context.Background(), teacherID)
	require.Nil(t, err)

	assert.Equal(t, teacher.ConvertToServiceTeacher(expectedFoundTeacher), foundTeacher)
}

func randomTeacher() teacher.Teacher {
	return teacher.Teacher{
		ID:          uuid.New(),
		FirstName:   gofakeit.FirstName(),
		LastName:    gofakeit.LastName(),
		MiddleName:  gofakeit.MiddleName(),
		ReportEmail: gofakeit.Email(),
		Username:    gofakeit.Username(),
	}
}
