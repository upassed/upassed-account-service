package student_test

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
	studentRepository student.Repository
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

	studentRepository = student.New(db, redis, cfg, logger)
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
	studentToSave := util.RandomDomainStudent()
	studentToSave.Username = gofakeit.LoremIpsumSentence(50)

	err := studentRepository.Save(context.Background(), studentToSave)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, codes.Internal, convertedError.Code())
	assert.Equal(t, student.ErrSavingStudent.Error(), convertedError.Message())
}

func TestSave_HappyPath(t *testing.T) {
	studentToSave := util.RandomDomainStudent()
	studentToSave.GroupID = uuid.MustParse("5eead8d5-b868-4708-aa25-713ad8399233")
	studentToSave.Group.ID = uuid.MustParse("5eead8d5-b868-4708-aa25-713ad8399233")

	err := studentRepository.Save(context.Background(), studentToSave)
	require.NoError(t, err)
}

func TestFindByID_StudentNotFound(t *testing.T) {
	randomStudentID := uuid.New()

	_, err := studentRepository.FindByID(context.Background(), randomStudentID)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, codes.NotFound, convertedError.Code())
	assert.Equal(t, student.ErrStudentNotFoundByID.Error(), convertedError.Message())
}

func TestFindByID_HappyPath(t *testing.T) {
	existingStudent := util.RandomDomainStudent()
	groupID := uuid.MustParse("5eead8d5-b868-4708-aa25-713ad8399233")
	existingStudent.GroupID = groupID
	existingStudent.Group.ID = groupID

	err := studentRepository.Save(context.Background(), existingStudent)
	require.NoError(t, err)

	foundStudent, err := studentRepository.FindByID(context.Background(), existingStudent.ID)
	require.NoError(t, err)

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
	duplicatesExists, err := studentRepository.CheckDuplicateExists(context.Background(), uniqueStudent.EducationalEmail, uniqueStudent.Username)
	require.NoError(t, err)
	assert.False(t, duplicatesExists)
}

func TestCheckDuplicates_DuplicatesExists(t *testing.T) {
	duplicateStudent := util.RandomDomainStudent()
	groupID := uuid.MustParse("5eead8d5-b868-4708-aa25-713ad8399233")
	duplicateStudent.GroupID = groupID
	duplicateStudent.Group.ID = groupID

	err := studentRepository.Save(context.Background(), duplicateStudent)
	require.NoError(t, err)

	duplicatesExists, err := studentRepository.CheckDuplicateExists(context.Background(), duplicateStudent.EducationalEmail, duplicateStudent.Username)
	require.NoError(t, err)

	assert.True(t, duplicatesExists)
}
