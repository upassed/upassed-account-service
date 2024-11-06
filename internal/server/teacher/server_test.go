package teacher_test

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

type mockTeacherService struct {
	mock.Mock
}

func (m *mockTeacherService) Create(ctx context.Context, teacher *business.Teacher) (*business.TeacherCreateResponse, error) {
	args := m.Called(ctx, teacher)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*business.TeacherCreateResponse), args.Error(1)
}

func (m *mockTeacherService) FindByID(ctx context.Context, teacherID uuid.UUID) (*business.Teacher, error) {
	args := m.Called(ctx, teacherID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*business.Teacher), args.Error(1)
}

var (
	teacherClient client.TeacherClient
	teacherSvc    *mockTeacherService
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

	teacherSvc = new(mockTeacherService)
	teacherServer := server.New(server.AppServerCreateParams{
		Config:         cfg,
		Log:            logger,
		TeacherService: teacherSvc,
	})

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	cc, err := grpc.NewClient(fmt.Sprintf(":%s", cfg.GrpcServer.Port), opts...)
	if err != nil {
		log.Fatal("error creating client connection: ", err)
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

func TestFindByID_InvalidRequest(t *testing.T) {
	request := client.TeacherFindByIDRequest{
		TeacherId: "invalid_uuid",
	}

	_, err := teacherClient.FindByID(context.Background(), &request)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, codes.InvalidArgument, convertedError.Code())
}

func TestFindByID_ServiceError(t *testing.T) {
	request := client.TeacherFindByIDRequest{
		TeacherId: uuid.NewString(),
	}

	expectedError := handling.New("some service error", codes.NotFound)
	teacherSvc.On("FindByID", mock.Anything, uuid.MustParse(request.TeacherId)).Return(nil, handling.Process(expectedError))

	_, err := teacherClient.FindByID(context.Background(), &request)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedError.Error(), convertedError.Message())
	assert.Equal(t, codes.NotFound, convertedError.Code())

	clearTeacherServiceMockCalls()
}

func TestFindByID_HappyPath(t *testing.T) {
	teacherID := uuid.New()
	request := client.TeacherFindByIDRequest{
		TeacherId: teacherID.String(),
	}

	foundTeacher := util.RandomBusinessTeacher()
	foundTeacher.ID = teacherID
	teacherSvc.On("FindByID", mock.Anything, teacherID).Return(foundTeacher, nil)

	response, err := teacherClient.FindByID(context.Background(), &request)
	require.NoError(t, err)

	assert.Equal(t, teacherID.String(), response.GetTeacher().GetId())

	clearTeacherServiceMockCalls()
}

func clearTeacherServiceMockCalls() {
	teacherSvc.ExpectedCalls = nil
	teacherSvc.Calls = nil
}
