package student_test

import (
	"context"
	"fmt"
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

func (m *mockStudentService) FindByID(ctx context.Context, studentID uuid.UUID) (*business.Student, error) {
	args := m.Called(ctx, studentID)

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

func TestFindByID_InvalidRequest(t *testing.T) {
	request := client.StudentFindByIDRequest{
		StudentId: "invalid_uuid",
	}

	_, err := studentClient.FindByID(context.Background(), &request)
	require.NotNil(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, codes.InvalidArgument, convertedError.Code())
}

func TestFindByID_ServiceError(t *testing.T) {
	request := client.StudentFindByIDRequest{
		StudentId: uuid.NewString(),
	}

	expectedError := handling.New("some service error", codes.NotFound)
	studentSvc.On("FindByID", mock.Anything, uuid.MustParse(request.StudentId)).Return(nil, handling.Process(expectedError))

	_, err := studentClient.FindByID(context.Background(), &request)
	require.NotNil(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedError.Error(), convertedError.Message())
	assert.Equal(t, codes.NotFound, convertedError.Code())

	clearStudentServiceMockCalls()
}

func TestFindByID_HappyPath(t *testing.T) {
	studentID := uuid.New()
	request := client.StudentFindByIDRequest{
		StudentId: studentID.String(),
	}

	foundStudent := util.RandomBusinessStudent()
	foundStudent.ID = studentID
	studentSvc.On("FindByID", mock.Anything, studentID).Return(foundStudent, nil)

	response, err := studentClient.FindByID(context.Background(), &request)
	require.Nil(t, err)

	assert.Equal(t, studentID.String(), response.GetStudent().GetId())

	clearStudentServiceMockCalls()
}

func clearStudentServiceMockCalls() {
	studentSvc.ExpectedCalls = nil
	studentSvc.Calls = nil
}
