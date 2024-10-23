package student_test

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
	"github.com/upassed/upassed-account-service/internal/repository/student"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	repository student.Repository
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
	repository, err = student.New(cfg, logger)
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
	_, err = student.New(cfg, logging.New(cfg.Env))
	require.NotNil(t, err)
	assert.ErrorIs(t, err, student.ErrOpeningDbConnection)
}

func TestSave_InvalidUsernameLength(t *testing.T) {
	studentToSave := util.RandomDomainStudent()
	studentToSave.Username = gofakeit.LoremIpsumSentence(50)

	err := repository.Save(context.Background(), studentToSave)
	require.NotNil(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, codes.Internal, convertedError.Code())
	assert.Equal(t, student.ErrSavingStudent.Error(), convertedError.Message())
}

func TestSave_HappyPath(t *testing.T) {
	studentToSave := util.RandomDomainStudent()
	studentToSave.GroupID = uuid.MustParse("5eead8d5-b868-4708-aa25-713ad8399233")
	studentToSave.Group.ID = uuid.MustParse("5eead8d5-b868-4708-aa25-713ad8399233")

	err := repository.Save(context.Background(), studentToSave)
	require.Nil(t, err)
}

func TestFindByID_StudentNotFound(t *testing.T) {
	randomStudentID := uuid.New()

	_, err := repository.FindByID(context.Background(), randomStudentID)
	require.NotNil(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, codes.NotFound, convertedError.Code())
	assert.Equal(t, student.ErrStudentNotFoundByID.Error(), convertedError.Message())
}

func TestFindByID_HappyPath(t *testing.T) {
	existingStudent := util.RandomDomainStudent()
	groupID := uuid.MustParse("5eead8d5-b868-4708-aa25-713ad8399233")
	existingStudent.GroupID = groupID
	existingStudent.Group.ID = groupID

	err := repository.Save(context.Background(), existingStudent)
	require.Nil(t, err)

	foundStudent, err := repository.FindByID(context.Background(), existingStudent.ID)
	require.Nil(t, err)

	assert.Equal(t, existingStudent.ID, foundStudent.ID)
	assert.Equal(t, existingStudent.FirstName, foundStudent.FirstName)
	assert.Equal(t, existingStudent.LastName, foundStudent.LastName)
	assert.Equal(t, existingStudent.MiddleName, foundStudent.MiddleName)
	assert.Equal(t, existingStudent.EducationalEmail, foundStudent.EducationalEmail)
	assert.Equal(t, existingStudent.Username, foundStudent.Username)
	assert.Equal(t, existingStudent.GroupID, foundStudent.GroupID)
	assert.Equal(t, existingStudent.Group.ID, foundStudent.Group.ID)
	assert.NotNil(t, foundStudent.Group.SpecializationCode)
	assert.NotNil(t, foundStudent.Group.GroupNumber)
}

func TestCheckDuplicates_DuplicatesNotExists(t *testing.T) {
	uniqueStudent := util.RandomDomainStudent()
	duplicatesExists, err := repository.CheckDuplicateExists(context.Background(), uniqueStudent.EducationalEmail, uniqueStudent.Username)
	require.Nil(t, err)
	assert.False(t, duplicatesExists)
}

func TestCheckDuplicates_DuplicatesExists(t *testing.T) {
	duplicateStudent := util.RandomDomainStudent()
	groupID := uuid.MustParse("5eead8d5-b868-4708-aa25-713ad8399233")
	duplicateStudent.GroupID = groupID
	duplicateStudent.Group.ID = groupID

	err := repository.Save(context.Background(), duplicateStudent)
	require.Nil(t, err)

	duplicatesExists, err := repository.CheckDuplicateExists(context.Background(), duplicateStudent.EducationalEmail, duplicateStudent.Username)
	require.Nil(t, err)

	assert.True(t, duplicatesExists)
}
