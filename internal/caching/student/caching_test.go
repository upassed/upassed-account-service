package student_test

import (
	"context"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/upassed/upassed-account-service/internal/caching"
	"github.com/upassed/upassed-account-service/internal/caching/student"
	"github.com/upassed/upassed-account-service/internal/config"
	"github.com/upassed/upassed-account-service/internal/logging"
	"github.com/upassed/upassed-account-service/internal/testcontainer"
	"github.com/upassed/upassed-account-service/internal/util"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

var (
	redisClient *student.RedisClient
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
	logger := logging.New(cfg.Env)

	redisTestcontainer, err := testcontainer.NewRedisTestcontainer(ctx, cfg)
	if err != nil {
		log.Fatal("unable to run redis testcontainer: ", err)
	}

	port, err := redisTestcontainer.Start(ctx)
	if err != nil {
		log.Fatal("unable to get a postgres testcontainer real port: ", err)
	}

	cfg.Redis.Port = strconv.Itoa(port)
	redis, err := caching.OpenRedisConnection(cfg, logger)
	if err != nil {
		log.Fatal("unable to open connections to redis: ", err)
	}

	redisClient = student.New(redis, cfg, logger)
	exitCode := m.Run()
	if err := redisTestcontainer.Stop(ctx); err != nil {
		log.Println("unable to stop redis testcontainer: ", err)
	}

	os.Exit(exitCode)
}

func TestSaveStudent_HappyPath(t *testing.T) {
	studentToSave := util.RandomDomainStudent()
	ctx := context.Background()
	err := redisClient.Save(ctx, studentToSave)
	require.NoError(t, err)

	studentFromCache, err := redisClient.GetByUsername(ctx, studentToSave.Username)
	require.NoError(t, err)

	assert.Equal(t, *studentToSave, *studentFromCache)
}

func TestFindStudentByUsername_StudentNotFound(t *testing.T) {
	studentUsername := gofakeit.Username()
	foundStudent, err := redisClient.GetByUsername(context.Background(), studentUsername)
	require.Error(t, err)

	assert.ErrorIs(t, err, student.ErrStudentIsNotPresentInCache)
	assert.Nil(t, foundStudent)
}

func TestFindStudentByUsername_StudentFound(t *testing.T) {
	studentToSave := util.RandomDomainStudent()
	ctx := context.Background()
	err := redisClient.Save(ctx, studentToSave)
	require.NoError(t, err)

	foundStudent, err := redisClient.GetByUsername(context.Background(), studentToSave.Username)
	require.NoError(t, err)

	assert.Equal(t, *studentToSave, *foundStudent)
}
