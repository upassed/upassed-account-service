package teacher_test

import (
	"context"
	"github.com/upassed/upassed-account-service/internal/caching"
	"github.com/upassed/upassed-account-service/internal/repository"
	"github.com/upassed/upassed-account-service/internal/testcontainer"
	"github.com/upassed/upassed-account-service/internal/util"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/upassed/upassed-account-service/internal/config"
	"github.com/upassed/upassed-account-service/internal/logging"
	"github.com/upassed/upassed-account-service/internal/repository/teacher"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	teacherRepository teacher.Repository
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
		log.Fatal("unable to parse cfg: ", err)
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
	db, err := repository.OpenGormDbConnection(cfg, logger)
	if err != nil {
		log.Fatal("unable to open connection to postgres: ", err)
	}

	redis, err := caching.OpenRedisConnection(cfg, logger)
	if err != nil {
		log.Fatal("unable to open connection to redis: ", err)
	}

	teacherRepository = teacher.New(db, redis, cfg, logger)
	exitCode := m.Run()
	if err := postgresTestcontainer.Stop(ctx); err != nil {
		log.Println("unable to stop postgres testcontainer: ", err)
	}

	if err := redisTestcontainer.Stop(ctx); err != nil {
		log.Println("unable to stop redis testcontainer: ", err)
	}

	os.Exit(exitCode)
}

func TestSave_InvalidUsernameLength(t *testing.T) {
	teacherToSave := util.RandomDomainTeacher()
	teacherToSave.Username = gofakeit.LoremIpsumSentence(50)

	err := teacherRepository.Save(context.Background(), teacherToSave)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, codes.Internal, convertedError.Code())
	assert.Equal(t, teacher.ErrSavingTeacher.Error(), convertedError.Message())
}

func TestSave_HappyPath(t *testing.T) {
	teacherToSave := util.RandomDomainTeacher()

	err := teacherRepository.Save(context.Background(), teacherToSave)
	require.NoError(t, err)
}

func TestFindByUsername_TeacherNotFound(t *testing.T) {
	randomTeacherUsername := gofakeit.Username()

	_, err := teacherRepository.FindByUsername(context.Background(), randomTeacherUsername)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, codes.NotFound, convertedError.Code())
	assert.Equal(t, teacher.ErrTeacherNotFoundByUsername.Error(), convertedError.Message())
}

func TestFindByUsername_HappyPath(t *testing.T) {
	existingTeacher := util.RandomDomainTeacher()
	err := teacherRepository.Save(context.Background(), existingTeacher)
	require.NoError(t, err)

	foundTeacher, err := teacherRepository.FindByUsername(context.Background(), existingTeacher.Username)
	require.NoError(t, err)

	assert.Equal(t, existingTeacher.ID, foundTeacher.ID)
	assert.Equal(t, existingTeacher.FirstName, foundTeacher.FirstName)
	assert.Equal(t, existingTeacher.LastName, foundTeacher.LastName)
	assert.Equal(t, existingTeacher.MiddleName, foundTeacher.MiddleName)
	assert.Equal(t, existingTeacher.ReportEmail, foundTeacher.ReportEmail)
	assert.Equal(t, existingTeacher.Username, foundTeacher.Username)
}

func TestCheckDuplicates_DuplicatesNotExists(t *testing.T) {
	uniqueTeacher := util.RandomDomainTeacher()
	duplicatesExists, err := teacherRepository.CheckDuplicateExists(context.Background(), uniqueTeacher.ReportEmail, uniqueTeacher.Username)
	require.NoError(t, err)
	assert.False(t, duplicatesExists)
}

func TestCheckDuplicates_DuplicatesExists(t *testing.T) {
	duplicateTeacher := util.RandomDomainTeacher()
	err := teacherRepository.Save(context.Background(), duplicateTeacher)
	require.NoError(t, err)

	duplicatesExists, err := teacherRepository.CheckDuplicateExists(context.Background(), duplicateTeacher.ReportEmail, duplicateTeacher.Username)
	require.NoError(t, err)

	assert.True(t, duplicatesExists)
}
