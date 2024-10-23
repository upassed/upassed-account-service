package caching_test

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/upassed/upassed-account-service/internal/caching"
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
	redisClient *caching.RedisClient
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
	redisClient, err = caching.New(cfg, logger)
	if err != nil {
		log.Fatal("unable to create repository: ", err)
	}

	exitCode := m.Run()
	if err := redisTestcontainer.Stop(ctx); err != nil {
		log.Fatal("unable to stop redis testcontainer: ", err)
	}

	os.Exit(exitCode)
}

func TestSaveGroup_HappyPath(t *testing.T) {
	groupToSave := util.RandomDomainGroup()
	ctx := context.Background()
	err := redisClient.SaveGroup(ctx, groupToSave)
	require.Nil(t, err)

	groupFromCache, err := redisClient.GetGroupByID(ctx, groupToSave.ID)
	require.Nil(t, err)

	assert.Equal(t, groupToSave, groupFromCache)
}

func TestFindGroupByID_GroupNotFound(t *testing.T) {
	groupID := uuid.New()
	_, err := redisClient.GetGroupByID(context.Background(), groupID)
	require.NotNil(t, err)

	assert.ErrorIs(t, err, caching.ErrGroupIsNotPresentInCache)
}

func TestSaveStudent_HappyPath(t *testing.T) {
	studentToSave := util.RandomDomainStudent()
	ctx := context.Background()
	err := redisClient.SaveStudent(ctx, studentToSave)
	require.Nil(t, err)

	studentFromCache, err := redisClient.GetStudentByID(ctx, studentToSave.ID)
	require.Nil(t, err)

	assert.Equal(t, studentToSave, studentFromCache)
}

func TestFindStudentByID_StudentNotFound(t *testing.T) {
	studentID := uuid.New()
	_, err := redisClient.GetStudentByID(context.Background(), studentID)
	require.NotNil(t, err)

	assert.ErrorIs(t, err, caching.ErrStudentIsNotPresentInCache)
}

func TestSaveTeacher_HappyPath(t *testing.T) {
	teacherToSave := util.RandomDomainTeacher()
	ctx := context.Background()
	err := redisClient.SaveTeacher(ctx, teacherToSave)
	require.Nil(t, err)

	teacherFromCache, err := redisClient.GetTeacherByID(ctx, teacherToSave.ID)
	require.Nil(t, err)

	assert.Equal(t, teacherToSave, teacherFromCache)
}

func TestFindTeacherByID_TeacherNotFound(t *testing.T) {
	teacherID := uuid.New()
	_, err := redisClient.GetTeacherByID(context.Background(), teacherID)
	require.NotNil(t, err)

	assert.ErrorIs(t, err, caching.ErrTeacherIsNotPresentInCache)
}
