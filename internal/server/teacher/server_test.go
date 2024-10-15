package teacher_test

import (
	"context"
	"errors"
	"fmt"
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
	"github.com/upassed/upassed-account-service/internal/logger"
	"github.com/upassed/upassed-account-service/internal/server"
	service "github.com/upassed/upassed-account-service/internal/service/teacher"
	"github.com/upassed/upassed-account-service/pkg/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type mockTeacherService struct {
	mock.Mock
}

func (m *mockTeacherService) Create(ctx context.Context, teacher service.Teacher) (service.TeacherCreateResponse, error) {
	args := m.Called(ctx, teacher)
	return args.Get(0).(service.TeacherCreateResponse), args.Error(1)
}

func (m *mockTeacherService) FindByID(ctx context.Context, teacherID string) (service.Teacher, error) {
	args := m.Called(ctx, teacherID)
	return args.Get(0).(service.Teacher), args.Error(1)
}

var (
	teacherClient client.TeacherClient
	teacherSvc    *mockTeacherService
)

func TestMain(m *testing.M) {
	projectRoot, err := getProjectRoot()
	if err != nil {
		log.Fatal("error to get project root folder", err)
	}

	if err := os.Setenv(config.EnvConfigPath, filepath.Join(projectRoot, "config", "test.yml")); err != nil {
		log.Fatal(err)
	}

	config, err := config.Load()
	if err != nil {
		log.Fatal("config load error", err)
	}

	logger := logger.New(config.Env)

	teacherSvc = new(mockTeacherService)
	teacherServer := server.New(server.AppServerCreateParams{
		Config:         config,
		Log:            logger,
		TeacherService: teacherSvc,
	})

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	cc, err := grpc.NewClient(fmt.Sprintf(":%s", config.GrpcServer.Port), opts...)
	if err != nil {
		log.Fatal("error creating client connection", err)
	}

	teacherClient = client.NewTeacherClient(cc)
	go func() {
		if err := teacherServer.Run(); err != nil {
			os.Exit(1)
		}
	}()

	exitCode := m.Run()
	teacherServer.GracefulStop()
	os.Exit(exitCode)
}

func TestCreate_InvalidRequest(t *testing.T) {
	request := client.TeacherCreateRequest{
		FirstName:   gofakeit.FirstName(),
		LastName:    gofakeit.LastName(),
		MiddleName:  gofakeit.MiddleName(),
		ReportEmail: "invalid_email",
		Username:    gofakeit.Username(),
	}

	_, err := teacherClient.Create(context.Background(), &request)
	require.NotNil(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, codes.InvalidArgument, convertedError.Code())
}

func TestCreate_ServiceError(t *testing.T) {
	request := client.TeacherCreateRequest{
		FirstName:   gofakeit.FirstName(),
		LastName:    gofakeit.LastName(),
		MiddleName:  gofakeit.MiddleName(),
		ReportEmail: gofakeit.Email(),
		Username:    gofakeit.Username(),
	}

	expectedError := handling.NewApplicationError("some service error", codes.AlreadyExists)
	teacherSvc.On("Create", mock.Anything, mock.Anything).Return(service.TeacherCreateResponse{}, handling.HandleApplicationError(expectedError))

	_, err := teacherClient.Create(context.Background(), &request)
	require.NotNil(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedError.Error(), convertedError.Message())
	assert.Equal(t, codes.AlreadyExists, convertedError.Code())
}

func TeacherCreate_HappyPath(t *testing.T) {
	request := client.TeacherCreateRequest{
		FirstName:   gofakeit.FirstName(),
		LastName:    gofakeit.LastName(),
		MiddleName:  gofakeit.MiddleName(),
		ReportEmail: gofakeit.Email(),
		Username:    gofakeit.Username(),
	}

	createdTeacherID := uuid.New()
	teacherSvc.On("Create", mock.Anything, mock.Anything).Return(service.TeacherCreateResponse{
		CreatedTeacherID: createdTeacherID,
	}, nil)

	response, err := teacherClient.Create(context.Background(), &request)
	require.Nil(t, err)

	assert.Equal(t, createdTeacherID, response.GetCreatedTeacherId())
}

func TestFindByID_InvalidRequest(t *testing.T) {
	request := client.TeacherFindByIDRequest{
		TeacherId: "invalid_uuid",
	}

	_, err := teacherClient.FindByID(context.Background(), &request)
	require.NotNil(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, codes.InvalidArgument, convertedError.Code())
}

func TestFindByID_ServiceError(t *testing.T) {
	request := client.TeacherFindByIDRequest{
		TeacherId: uuid.NewString(),
	}

	expectedError := handling.NewApplicationError("some service error", codes.NotFound)
	teacherSvc.On("FindByID", mock.Anything, request.TeacherId).Return(service.Teacher{}, handling.HandleApplicationError(expectedError))

	_, err := teacherClient.FindByID(context.Background(), &request)
	require.NotNil(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedError.Error(), convertedError.Message())
	assert.Equal(t, codes.NotFound, convertedError.Code())
}

func TestFindByID_HappyPath(t *testing.T) {
	teacherID := uuid.New()
	request := client.TeacherFindByIDRequest{
		TeacherId: teacherID.String(),
	}

	teacherSvc.On("FindByID", mock.Anything, teacherID.String()).Return(service.Teacher{
		ID:          teacherID,
		FirstName:   gofakeit.FirstName(),
		LastName:    gofakeit.LastName(),
		MiddleName:  gofakeit.MiddleName(),
		ReportEmail: gofakeit.Email(),
		Username:    gofakeit.Username(),
	}, nil)

	response, err := teacherClient.FindByID(context.Background(), &request)
	require.Nil(t, err)

	assert.Equal(t, teacherID.String(), response.GetTeacher().GetId())
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
