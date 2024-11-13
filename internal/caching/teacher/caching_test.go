package teacher_test

import (
	"context"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/upassed/upassed-account-service/internal/caching"
	"github.com/upassed/upassed-account-service/internal/caching/teacher"
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
	redisClient *teacher.RedisClient
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

	redisClient = teacher.New(redis, cfg, logger)
	exitCode := m.Run()
	if err := redisTestcontainer.Stop(ctx); err != nil {
		log.Println("unable to stop redis testcontainer: ", err)
	}

	os.Exit(exitCode)
}

func TestSaveTeacher_HappyPath(t *testing.T) {
	teacherToSave := util.RandomDomainTeacher()
	ctx := context.Background()
	err := redisClient.Save(ctx, teacherToSave)
	require.NoError(t, err)

	teacherFromCache, err := redisClient.GetByUsername(ctx, teacherToSave.Username)
	require.NoError(t, err)

	assert.Equal(t, *teacherToSave, *teacherFromCache)
}

func TestFindTeacherByUsername_TeacherNotFound(t *testing.T) {
	teacherUsername := gofakeit.Username()
	foundTeacher, err := redisClient.GetByUsername(context.Background(), teacherUsername)
	require.Error(t, err)

	assert.ErrorIs(t, err, teacher.ErrTeacherUsernameIsNotPresentInCache)
	assert.Nil(t, foundTeacher)
}

func TestFindTeacherByUsername_TeacherFound(t *testing.T) {
	teacherToSave := util.RandomDomainTeacher()
	ctx := context.Background()
	err := redisClient.Save(ctx, teacherToSave)
	require.NoError(t, err)

	foundTeacher, err := redisClient.GetByUsername(ctx, teacherToSave.Username)
	require.NoError(t, err)
	assert.NotNil(t, foundTeacher)
}
