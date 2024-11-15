package student_test

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/golang/mock/gomock"
	"github.com/upassed/upassed-account-service/internal/util"
	"github.com/upassed/upassed-account-service/internal/util/mocks"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/upassed/upassed-account-service/internal/config"
	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/internal/logging"
	"github.com/upassed/upassed-account-service/internal/server"
	"github.com/upassed/upassed-account-service/pkg/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

var (
	studentClient client.StudentClient
	studentSvc    *mocks.StudentService
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
	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()

	studentSvc = mocks.NewStudentService(ctrl)
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
	studentSvc.EXPECT().
		FindByUsername(gomock.Any(), request.GetStudentUsername()).
		Return(nil, handling.Process(expectedError))

	_, err := studentClient.FindByUsername(context.Background(), &request)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedError.Error(), convertedError.Message())
	assert.Equal(t, codes.NotFound, convertedError.Code())
}

func TestFindByUsername_HappyPath(t *testing.T) {
	studentUsername := gofakeit.Username()
	request := client.StudentFindByUsernameRequest{
		StudentUsername: studentUsername,
	}

	foundStudent := util.RandomBusinessStudent()
	foundStudent.Username = studentUsername
	studentSvc.EXPECT().
		FindByUsername(gomock.Any(), request.GetStudentUsername()).
		Return(foundStudent, nil)

	response, err := studentClient.FindByUsername(context.Background(), &request)
	require.NoError(t, err)

	assert.Equal(t, studentUsername, response.GetStudent().GetUsername())
}
