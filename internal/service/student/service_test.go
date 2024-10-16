package student_test

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
	"github.com/upassed/upassed-account-service/internal/repository/group"
	"github.com/upassed/upassed-account-service/internal/repository/student"
	groupService "github.com/upassed/upassed-account-service/internal/service/group"
	service "github.com/upassed/upassed-account-service/internal/service/student"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mockStudentRepository struct {
	mock.Mock
}

func (m *mockStudentRepository) Save(ctx context.Context, student student.Student) error {
	args := m.Called(ctx, student)
	return args.Error(0)
}

func (m *mockStudentRepository) FindByID(ctx context.Context, studentID uuid.UUID) (student.Student, error) {
	args := m.Called(ctx, studentID)
	return args.Get(0).(student.Student), args.Error(1)
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

func (m *mockGroupRepository) FindByID(ctx context.Context, groupID uuid.UUID) (group.Group, error) {
	args := m.Called(ctx, groupID)
	return args.Get(0).(group.Group), args.Error(1)
}

func TestCreate_ErrorCheckingDuplicateExists(t *testing.T) {
	studentRepository := new(mockStudentRepository)
	groupRepository := new(mockGroupRepository)

	studentToCreate := randomStudent()
	expectedReposotiryError := errors.New("some repo error")
	studentRepository.On(
		"CheckDuplicateExists",
		mock.Anything,
		studentToCreate.EducationalEmail,
		studentToCreate.Username,
	).Return(false, expectedReposotiryError)

	service := service.New(logger.New(config.EnvTesting), studentRepository, groupRepository)
	_, err := service.Create(context.Background(), studentToCreate)
	require.NotNil(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedReposotiryError.Error(), convertedError.Message())
	assert.Equal(t, codes.Internal, convertedError.Code())
}

func TestCreate_DuplicateExists(t *testing.T) {
	studentRepository := new(mockStudentRepository)
	groupRepository := new(mockGroupRepository)

	studentToCreate := randomStudent()
	studentRepository.On(
		"CheckDuplicateExists",
		mock.Anything,
		studentToCreate.EducationalEmail,
		studentToCreate.Username,
	).Return(true, nil)

	service := service.New(logger.New(config.EnvTesting), studentRepository, groupRepository)
	_, err := service.Create(context.Background(), studentToCreate)
	require.NotNil(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, "student duplicate found", convertedError.Message())
	assert.Equal(t, codes.AlreadyExists, convertedError.Code())
}

func TestCreate_ErrorCheckingGroupExists(t *testing.T) {
	studentRepository := new(mockStudentRepository)
	groupRepository := new(mockGroupRepository)

	studentToCreate := randomStudent()
	studentRepository.On(
		"CheckDuplicateExists",
		mock.Anything,
		studentToCreate.EducationalEmail,
		studentToCreate.Username,
	).Return(false, nil)

	expectedRepositoryError := errors.New("some repo error")
	groupRepository.On("Exists", mock.Anything, studentToCreate.Group.ID).Return(false, expectedRepositoryError)

	service := service.New(logger.New(config.EnvTesting), studentRepository, groupRepository)
	_, err := service.Create(context.Background(), studentToCreate)
	require.NotNil(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedRepositoryError.Error(), convertedError.Message())
	assert.Equal(t, codes.Internal, convertedError.Code())
}

func TestCreate_GroupNotExists(t *testing.T) {
	studentRepository := new(mockStudentRepository)
	groupRepository := new(mockGroupRepository)

	studentToCreate := randomStudent()
	studentRepository.On(
		"CheckDuplicateExists",
		mock.Anything,
		studentToCreate.EducationalEmail,
		studentToCreate.Username,
	).Return(false, nil)

	groupRepository.On("Exists", mock.Anything, studentToCreate.Group.ID).Return(false, nil)

	service := service.New(logger.New(config.EnvTesting), studentRepository, groupRepository)
	_, err := service.Create(context.Background(), studentToCreate)
	require.NotNil(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, "group does not exists by id", convertedError.Message())
	assert.Equal(t, codes.NotFound, convertedError.Code())
}

func TestCreate_ErrorSavingToDatabase(t *testing.T) {
	studentRepository := new(mockStudentRepository)
	groupRepository := new(mockGroupRepository)

	studentToCreate := randomStudent()
	studentRepository.On(
		"CheckDuplicateExists",
		mock.Anything,
		studentToCreate.EducationalEmail,
		studentToCreate.Username,
	).Return(false, nil)

	groupRepository.On("Exists", mock.Anything, studentToCreate.Group.ID).Return(true, nil)

	expectedRepositoryError := errors.New("some repo error")
	studentRepository.On("Save", mock.Anything, mock.Anything).Return(expectedRepositoryError)

	service := service.New(logger.New(config.EnvTesting), studentRepository, groupRepository)
	_, err := service.Create(context.Background(), studentToCreate)
	require.NotNil(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedRepositoryError.Error(), convertedError.Message())
	assert.Equal(t, codes.Internal, convertedError.Code())
}

func TestCreate_HappyPath(t *testing.T) {
	studentRepository := new(mockStudentRepository)
	groupRepository := new(mockGroupRepository)

	studentToCreate := randomStudent()
	studentRepository.On(
		"CheckDuplicateExists",
		mock.Anything,
		studentToCreate.EducationalEmail,
		studentToCreate.Username,
	).Return(false, nil)

	groupRepository.On("Exists", mock.Anything, studentToCreate.Group.ID).Return(true, nil)
	studentRepository.On("Save", mock.Anything, mock.Anything).Return(nil)

	service := service.New(logger.New(config.EnvTesting), studentRepository, groupRepository)
	response, err := service.Create(context.Background(), studentToCreate)
	require.Nil(t, err)

	assert.Equal(t, studentToCreate.ID, response.CreatedStudentID)
}

func randomStudent() service.Student {
	return service.Student{
		ID:               uuid.New(),
		FirstName:        gofakeit.FirstName(),
		LastName:         gofakeit.LastName(),
		MiddleName:       gofakeit.MiddleName(),
		EducationalEmail: gofakeit.Email(),
		Username:         gofakeit.UserAgent(),
		Group: groupService.Group{
			ID: uuid.New(),
		},
	}
}
