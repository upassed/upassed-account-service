package group_test

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/upassed/upassed-account-service/internal/util"
	"github.com/upassed/upassed-account-service/internal/util/mocks"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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

var (
	groupClient client.GroupClient
	groupSvc    *mocks.GroupService
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

	groupSvc = mocks.NewGroupService(ctrl)
	groupServer := server.New(server.AppServerCreateParams{
		Config:       cfg,
		Log:          logger,
		GroupService: groupSvc,
	})

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	cc, err := grpc.NewClient(fmt.Sprintf(":%s", cfg.GrpcServer.Port), opts...)
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
	groupSvc.EXPECT().
		FindStudentsInGroup(gomock.Any(), uuid.MustParse(request.GetGroupId())).
		Return(nil, handling.Process(expectedError))

	_, err := groupClient.FindStudentsInGroup(context.Background(), &request)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedError.Error(), convertedError.Message())
	assert.Equal(t, codes.NotFound, convertedError.Code())
}

func TestFindStudentsInGroup_HappyPath(t *testing.T) {
	request := client.FindStudentsInGroupRequest{
		GroupId: uuid.NewString(),
	}

	studentsInGroup := []*business.Student{
		util.RandomBusinessStudent(),
		util.RandomBusinessStudent(),
		util.RandomBusinessStudent(),
	}

	groupSvc.EXPECT().
		FindStudentsInGroup(gomock.Any(), uuid.MustParse(request.GetGroupId())).
		Return(studentsInGroup, nil)

	response, err := groupClient.FindStudentsInGroup(context.Background(), &request)
	require.NoError(t, err)

	assert.Equal(t, len(studentsInGroup), len(response.GetStudentsInGroup()))
	for idx := range studentsInGroup {
		assertStudentsEqual(t, studentsInGroup[idx], response.GetStudentsInGroup()[idx])
	}
}

func TestFindByID_ServiceLayerError(t *testing.T) {
	request := client.GroupFindByIDRequest{
		GroupId: uuid.NewString(),
	}

	expectedError := handling.New("some service error", codes.NotFound)
	groupSvc.EXPECT().
		FindByID(gomock.Any(), uuid.MustParse(request.GetGroupId())).
		Return(nil, expectedError)

	_, err := groupClient.FindByID(context.Background(), &request)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedError.Error(), convertedError.Message())
	assert.Equal(t, codes.NotFound, convertedError.Code())
}

func TestFindByID_HappyPath(t *testing.T) {
	request := client.GroupFindByIDRequest{
		GroupId: uuid.NewString(),
	}

	expectedFoundGroup := util.RandomBusinessGroup()
	groupSvc.EXPECT().
		FindByID(gomock.Any(), uuid.MustParse(request.GetGroupId())).
		Return(expectedFoundGroup, nil)

	response, err := groupClient.FindByID(context.Background(), &request)
	require.NoError(t, err)

	assert.Equal(t, expectedFoundGroup.ID.String(), response.GetGroup().GetId())
	assert.Equal(t, expectedFoundGroup.SpecializationCode, response.GetGroup().GetSpecializationCode())
	assert.Equal(t, expectedFoundGroup.GroupNumber, response.GetGroup().GetGroupNumber())
}

func TestFindByFilter_InvalidRequest(t *testing.T) {
	request := client.GroupSearchByFilterRequest{
		SpecializationCode: gofakeit.LoremIpsumSentence(50),
		GroupNumber:        gofakeit.LoremIpsumSentence(50),
	}

	_, err := groupClient.SearchByFilter(context.Background(), &request)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, codes.InvalidArgument, convertedError.Code())
}

func TestFindByFilter_ServiceError(t *testing.T) {
	request := client.GroupSearchByFilterRequest{
		SpecializationCode: "5130904",
		GroupNumber:        "10101",
	}

	expectedServiceError := handling.New("some service error", codes.DeadlineExceeded)
	groupSvc.EXPECT().
		FindByFilter(gomock.Any(), gomock.Any()).
		Return(nil, expectedServiceError)

	_, err := groupClient.SearchByFilter(context.Background(), &request)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, codes.DeadlineExceeded, convertedError.Code())
	assert.Equal(t, expectedServiceError.Error(), convertedError.Message())
}

func TestFindByFilter_HappyPath(t *testing.T) {
	request := client.GroupSearchByFilterRequest{
		SpecializationCode: "5130904",
		GroupNumber:        "10101",
	}

	expectedMatchedGroups := []*business.Group{
		util.RandomBusinessGroup(),
		util.RandomBusinessGroup(),
		util.RandomBusinessGroup(),
	}

	groupSvc.EXPECT().
		FindByFilter(gomock.Any(), gomock.Any()).
		Return(expectedMatchedGroups, nil)

	response, err := groupClient.SearchByFilter(context.Background(), &request)
	require.NoError(t, err)

	assert.Equal(t, len(expectedMatchedGroups), len(response.MatchedGroups))
	for idx := range expectedMatchedGroups {
		assertGroupsEqual(t, expectedMatchedGroups[idx], response.GetMatchedGroups()[idx])
	}
}
