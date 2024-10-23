package group_test

import (
	"context"
	"github.com/upassed/upassed-account-service/internal/testcontainer"
	"github.com/upassed/upassed-account-service/internal/util"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/upassed/upassed-account-service/internal/config"
	"github.com/upassed/upassed-account-service/internal/logging"
	"github.com/upassed/upassed-account-service/internal/repository/group"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	repository group.Repository
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
		log.Fatal("unable to parse config: ", err)
	}

	ctx := context.Background()
	postgresTestcontainer, err := testcontainer.NewPostgresTestcontainer(ctx)
	if err != nil {
		log.Fatal("unable to create a testcontainer: ", err)
	}

	port, err := postgresTestcontainer.Start(ctx)
	if err != nil {
		log.Fatal("unable to get a postgres testcontainer real port: ", err)
	}

	cfg.Storage.Port = strconv.Itoa(port)
	logger := logging.New(cfg.Env)
	if err := postgresTestcontainer.Migrate(cfg, logger); err != nil {
		log.Fatal("unable to run migrations: ", err)
	}

	redisTestcontainer, err := testcontainer.NewRedisTestcontainer(ctx, cfg)
	if err != nil {
		log.Fatal("unable to run redis testcontainer: ", err)
	}

	port, err = redisTestcontainer.Start(ctx)
	if err != nil {
		log.Fatal("unable to get a postgres testcontainer real port: ", err)
	}

	cfg.Redis.Port = strconv.Itoa(port)
	repository, err = group.New(cfg, logger)
	if err != nil {
		log.Fatal("unable to create repository: ", err)
	}

	exitCode := m.Run()
	if err := postgresTestcontainer.Stop(ctx); err != nil {
		log.Fatal("unable to stop postgres testcontainer: ", err)
	}

	if err := redisTestcontainer.Stop(ctx); err != nil {
		log.Fatal("unable to stop redis testcontainer: ", err)
	}

	os.Exit(exitCode)
}

func TestConnectToDatabase_InvalidCredentials(t *testing.T) {
	cfg, err := config.Load()
	require.Nil(t, err)

	cfg.Storage.DatabaseName = "invalid-db-name"
	_, err = group.New(cfg, logging.New(cfg.Env))
	require.NotNil(t, err)
	assert.ErrorIs(t, err, group.ErrOpeningDbConnection)
}

func TestExists_GroupNotExists(t *testing.T) {
	groupID := uuid.New()
	exists, err := repository.Exists(context.Background(), groupID)
	require.Nil(t, err)
	assert.False(t, exists)
}

func TestExists_GroupExists(t *testing.T) {
	groupID := uuid.MustParse("5eead8d5-b868-4708-aa25-713ad8399233")
	exists, err := repository.Exists(context.Background(), groupID)
	require.Nil(t, err)
	assert.True(t, exists)
}

func TestFindStudentsInGroup_StudentsNotFound(t *testing.T) {
	groupID := uuid.New()
	studentsInGroup, err := repository.FindStudentsInGroup(context.Background(), groupID)
	require.Nil(t, err)
	assert.Equal(t, 0, len(studentsInGroup))
}

func TestFindStudentsInGroup_StudentsExistsInGroup(t *testing.T) {
	// TODO - need to run test migration script with students creation
}

func TestFindByID_GroupNotFound(t *testing.T) {
	groupID := uuid.New()
	_, err := repository.FindByID(context.Background(), groupID)
	require.NotNil(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, codes.NotFound, convertedError.Code())
	assert.Equal(t, group.ErrGroupNotFoundByID.Error(), convertedError.Message())
}

func TestFindByID_GroupFound(t *testing.T) {
	groupID := uuid.MustParse("5eead8d5-b868-4708-aa25-713ad8399233")
	foundGroup, err := repository.FindByID(context.Background(), groupID)
	require.Nil(t, err)

	assert.Equal(t, groupID, foundGroup.ID)
	assert.Equal(t, "5130904", foundGroup.SpecializationCode)
	assert.Equal(t, "10101", foundGroup.GroupNumber)
}

func TestFindByFilter_NothingMatched(t *testing.T) {
	filter := domain.GroupFilter{
		SpecializationCode: "----",
		GroupNumber:        "----",
	}

	matchedGroups, err := repository.FindByFilter(context.Background(), filter)
	require.Nil(t, err)

	assert.Equal(t, 0, len(matchedGroups))
}

func TestFindByFilter_HasMatchedGroups(t *testing.T) {
	filter := domain.GroupFilter{
		SpecializationCode: "513",
		GroupNumber:        "10101",
	}

	matchedGroups, err := repository.FindByFilter(context.Background(), filter)
	require.Nil(t, err)

	assert.Equal(t, 1, len(matchedGroups))
	assert.Equal(t, uuid.MustParse("5eead8d5-b868-4708-aa25-713ad8399233"), matchedGroups[0].ID)
}
