package student_test

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/upassed/upassed-account-service/internal/util"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/upassed/upassed-account-service/internal/config"
	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/internal/logging"
	"github.com/upassed/upassed-account-service/internal/server"
	business "github.com/upassed/upassed-account-service/internal/service/model"
	"github.com/upassed/upassed-account-service/pkg/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type mockStudentService struct {
	mock.Mock
}

func (m *mockStudentService) Create(ctx context.Context, student *business.Student) (*business.StudentCreateResponse, error) {
	args := m.Called(ctx, student)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*business.StudentCreateResponse), args.Error(1)
}

func (m *mockStudentService) FindByUsername(ctx context.Context, studentUsername string) (*business.Student, error) {
	args := m.Called(ctx, studentUsername)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*business.Student), args.Error(1)
}

var (
	studentClient client.StudentClient
	studentSvc    *mockStudentService
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

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("cfg load error: ", err)
	}

	logger := logging.New(cfg.Env)
	studentSvc = new(mockStudentService)
	studentServer := server.New(server.AppServerCreateParams{
		Config:         cfg,
		Log:            logger,
		StudentService: studentSvc,
	})

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	cc, err := grpc.NewClient(fmt.Sprintf(":%s", cfg.GrpcServer.Port), opts...)
	if err != nil {
		log.Fatal("error creating client connection", err)
	}

	studentClient = client.NewStudentClient(cc)
	go func() {
		if err := studentServer.Run(); err != nil {
			os.Exit(1)
		}
	}()

	exitCode := m.Run()
	studentServer.GracefulStop()
	os.Exit(exitCode)
}

func TestFindByUsername_InvalidRequest(t *testing.T) {
	request := client.StudentFindByUsernameRequest{
		StudentUsername: "0",
	}

	_, err := studentClient.FindByUsername(context.Background(), &request)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, codes.InvalidArgument, convertedError.Code())
}

func TestFindByUsername_ServiceError(t *testing.T) {
	request := client.StudentFindByUsernameRequest{
		StudentUsername: gofakeit.Username(),
	}

	expectedError := handling.New("some service error", codes.NotFound)
	studentSvc.On("FindByUsername", mock.Anything, request.GetStudentUsername()).Return(nil, handling.Process(expectedError))

	_, err := studentClient.FindByUsername(context.Background(), &request)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedError.Error(), convertedError.Message())
	assert.Equal(t, codes.NotFound, convertedError.Code())

	clearStudentServiceMockCalls()
}

func TestFindByUsername_HappyPath(t *testing.T) {
	studentUsername := gofakeit.Username()
	request := client.StudentFindByUsernameRequest{
		StudentUsername: studentUsername,
	}

	foundStudent := util.RandomBusinessStudent()
	foundStudent.Username = studentUsername
	studentSvc.On("FindByUsername", mock.Anything, studentUsername).Return(foundStudent, nil)

	response, err := studentClient.FindByUsername(context.Background(), &request)
	require.NoError(t, err)

	assert.Equal(t, studentUsername, response.GetStudent().GetUsername())

	clearStudentServiceMockCalls()
}

func clearStudentServiceMockCalls() {
	studentSvc.ExpectedCalls = nil
	studentSvc.Calls = nil
}
