package group_test

import (
	"context"
	"errors"
	"fmt"
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
	"github.com/upassed/upassed-account-service/internal/logger"
	"github.com/upassed/upassed-account-service/internal/server"
	"github.com/upassed/upassed-account-service/internal/service/student"
	"github.com/upassed/upassed-account-service/pkg/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type mockGroupService struct {
	mock.Mock
}

func (m *mockGroupService) FindStudentsInGroup(ctx context.Context, groupID uuid.UUID) ([]student.Student, error) {
	args := m.Called(ctx, groupID)
	return args.Get(0).([]student.Student), args.Error(1)
}

var (
	groupClient client.GroupClient
	groupSvc    *mockGroupService
)

func TestMain(m *testing.M) {
	projectRoot, err := getProjectRoot()
	if err != nil {
		log.Fatal("error to get project root folder: ", err)
	}

	if err := os.Setenv(config.EnvConfigPath, filepath.Join(projectRoot, "config", "test.yml")); err != nil {
		log.Fatal(err)
	}

	config, err := config.Load()
	if err != nil {
		log.Fatal("config load error: ", err)
	}

	logger := logger.New(config.Env)
	groupSvc = new(mockGroupService)
	groupServer := server.New(server.AppServerCreateParams{
		Config:       config,
		Log:          logger,
		GroupService: groupSvc,
	})

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	cc, err := grpc.NewClient(fmt.Sprintf(":%s", config.GrpcServer.Port), opts...)
	if err != nil {
		log.Fatal("error creating client connection", err)
	}

	groupClient = client.NewGroupClient(cc)
	go func() {
		if err := groupServer.Run(); err != nil {
			os.Exit(1)
		}
	}()

	exitCode := m.Run()
	groupServer.GracefulStop()
	os.Exit(exitCode)
}

func TestFindStudentsInGroup_ServiceError(t *testing.T) {
	request := client.FindStudentsInGroupRequest{
		GroupId: uuid.NewString(),
	}

	expectedError := handling.New("some service error", codes.NotFound)
	groupSvc.On("FindStudentsInGroup", mock.Anything, uuid.MustParse(request.GetGroupId())).Return(make([]student.Student, 0), handling.Process(expectedError))

	_, err := groupClient.FindStudentsInGroup(context.Background(), &request)
	require.NotNil(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedError.Error(), convertedError.Message())
	assert.Equal(t, codes.NotFound, convertedError.Code())

	clearGroupServiceMockCalls()
}

func TestFindStudentsInGroup_HappyPath(t *testing.T) {
	request := client.FindStudentsInGroupRequest{
		GroupId: uuid.NewString(),
	}

	studentsInGroup := []student.Student{randomStudent(), randomStudent(), randomStudent()}
	groupSvc.On("FindStudentsInGroup", mock.Anything, uuid.MustParse(request.GetGroupId())).Return(studentsInGroup, nil)

	response, err := groupClient.FindStudentsInGroup(context.Background(), &request)
	require.Nil(t, err)

	assert.Equal(t, len(studentsInGroup), len(response.GetStudentsInGroup()))
	for idx := range studentsInGroup {
		assertStudentsEqual(t, studentsInGroup[idx], response.GetStudentsInGroup()[idx])
	}

	clearGroupServiceMockCalls()
}

func clearGroupServiceMockCalls() {
	groupSvc.ExpectedCalls = nil
	groupSvc.Calls = nil
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
