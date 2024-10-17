package group_test

import (
	"context"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/upassed/upassed-account-service/internal/config"
	"github.com/upassed/upassed-account-service/internal/logger"
	testcontainer "github.com/upassed/upassed-account-service/internal/repository"
	"github.com/upassed/upassed-account-service/internal/repository/group"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type groupRepository interface {
	Exists(context.Context, uuid.UUID) (bool, error)
	FindStudentsInGroup(context.Context, uuid.UUID) ([]domain.Student, error)
	FindByID(context.Context, uuid.UUID) (domain.Group, error)
}

var (
	repository groupRepository
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
		log.Fatal("unable to parse config: ", err)
	}

	ctx := context.Background()
	container, err := testcontainer.NewPostgresTestontainer(ctx)
	if err != nil {
		log.Fatal("unable to create a testcontainer: ", err)
	}

	port, err := container.Start(ctx)
	if err != nil {
		log.Fatal("unable to get a postgres testcontainer real port: ", err)
	}

	config.Storage.Port = strconv.Itoa(port)
	logger := logger.New(config.Env)
	if err := container.Migrate(config, logger); err != nil {
		log.Fatal("unable to run migrations: ", err)
	}

	repository, err = group.New(config, logger)
	if err != nil {
		log.Fatal("unable to create repository: ", err)
	}

	exitCode := m.Run()
	if err := container.Stop(ctx); err != nil {
		log.Fatal("unable to stop postgres testcontainer: ", err)
	}

	os.Exit(exitCode)
}

func TestConnectToDatabase_InvalidCredentials(t *testing.T) {
	config, err := config.Load()
	require.Nil(t, err)

	config.Storage.DatabaseName = "invalid-db-name"
	_, err = group.New(config, logger.New(config.Env))
	require.NotNil(t, err)
	assert.ErrorIs(t, err, group.ErrorOpeningDbConnection)
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
	assert.Equal(t, group.ErrorGroupNotFound.Error(), convertedError.Message())
}

func TestFindByID_GroupFound(t *testing.T) {
	groupID := uuid.MustParse("5eead8d5-b868-4708-aa25-713ad8399233")
	group, err := repository.FindByID(context.Background(), groupID)
	require.Nil(t, err)

	assert.Equal(t, groupID, group.ID)
	assert.Equal(t, "5130904", group.SpecializationCode)
	assert.Equal(t, "10101", group.GroupNumber)
}
