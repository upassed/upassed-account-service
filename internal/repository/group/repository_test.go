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
)

type groupRepository interface {
	Exists(context.Context, uuid.UUID) (bool, error)
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
