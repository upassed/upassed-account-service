package teacher_test

import (
	"context"
	"github.com/upassed/upassed-account-service/internal/testcontainer"
	"github.com/upassed/upassed-account-service/internal/util"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/upassed/upassed-account-service/internal/config"
	"github.com/upassed/upassed-account-service/internal/logging"
	"github.com/upassed/upassed-account-service/internal/repository/teacher"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	repository teacher.Repository
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
	repository, err = teacher.New(cfg, logger)
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
	_, err = teacher.New(cfg, logging.New(cfg.Env))
	require.NotNil(t, err)
	assert.ErrorIs(t, err, teacher.ErrOpeningDbConnection)
}

func TestSave_InvalidUsernameLength(t *testing.T) {
	teacherToSave := util.RandomDomainTeacher()
	teacherToSave.Username = gofakeit.LoremIpsumSentence(50)

	err := repository.Save(context.Background(), teacherToSave)
	require.NotNil(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, codes.Internal, convertedError.Code())
	assert.Equal(t, teacher.ErrSavingTeacher.Error(), convertedError.Message())
}

func TestSave_HappyPath(t *testing.T) {
	teacherToSave := util.RandomDomainTeacher()

	err := repository.Save(context.Background(), teacherToSave)
	require.Nil(t, err)
}

func TestFindByID_TeacherNotFound(t *testing.T) {
	randomTeacherID := uuid.New()

	_, err := repository.FindByID(context.Background(), randomTeacherID)
	require.NotNil(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, codes.NotFound, convertedError.Code())
	assert.Equal(t, teacher.ErrTeacherNotFoundByID.Error(), convertedError.Message())
}

func TestFindByID_HappyPath(t *testing.T) {
	existingTeacher := util.RandomDomainTeacher()
	err := repository.Save(context.Background(), existingTeacher)
	require.Nil(t, err)

	foundTeacher, err := repository.FindByID(context.Background(), existingTeacher.ID)
	require.Nil(t, err)

	assert.Equal(t, existingTeacher.ID, foundTeacher.ID)
	assert.Equal(t, existingTeacher.FirstName, foundTeacher.FirstName)
	assert.Equal(t, existingTeacher.LastName, foundTeacher.LastName)
	assert.Equal(t, existingTeacher.MiddleName, foundTeacher.MiddleName)
	assert.Equal(t, existingTeacher.ReportEmail, foundTeacher.ReportEmail)
	assert.Equal(t, existingTeacher.Username, foundTeacher.Username)
}

func TestCheckDuplicates_DuplicatesNotExists(t *testing.T) {
	uniqueTeacher := util.RandomDomainTeacher()
	duplicatesExists, err := repository.CheckDuplicateExists(context.Background(), uniqueTeacher.ReportEmail, uniqueTeacher.Username)
	require.Nil(t, err)
	assert.False(t, duplicatesExists)
}

func TestCheckDuplicates_DuplicatesExists(t *testing.T) {
	duplicateTeacher := util.RandomDomainTeacher()
	err := repository.Save(context.Background(), duplicateTeacher)
	require.Nil(t, err)

	duplicatesExists, err := repository.CheckDuplicateExists(context.Background(), duplicateTeacher.ReportEmail, duplicateTeacher.Username)
	require.Nil(t, err)

	assert.True(t, duplicatesExists)
}
