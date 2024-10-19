package teacher_test

import (
	"context"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	config "github.com/upassed/upassed-account-service/internal/config"
	"github.com/upassed/upassed-account-service/internal/logger"
	testcontainer "github.com/upassed/upassed-account-service/internal/repository"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	"github.com/upassed/upassed-account-service/internal/repository/teacher"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type teacherRepository interface {
	Save(context.Context, domain.Teacher) error
	FindByID(context.Context, uuid.UUID) (domain.Teacher, error)
	CheckDuplicateExists(ctx context.Context, reportEmail, username string) (bool, error)
}

var (
	repository teacherRepository
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

	repository, err = teacher.New(config, logger)
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
	_, err = teacher.New(config, logger.New(config.Env))
	require.NotNil(t, err)
	assert.ErrorIs(t, err, teacher.ErrorOpeningDbConnection)
}

func TestSave_InvalidUsernameLength(t *testing.T) {
	teacherToSave := randomTeacher()
	teacherToSave.Username = gofakeit.LoremIpsumSentence(50)

	err := repository.Save(context.Background(), teacherToSave)
	require.NotNil(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, codes.Internal, convertedError.Code())
	assert.Equal(t, teacher.ErrorSavingTeacher.Error(), convertedError.Message())
}

func TestSave_HappyPath(t *testing.T) {
	teacherToSave := randomTeacher()

	err := repository.Save(context.Background(), teacherToSave)
	require.Nil(t, err)
}

func TestFindByID_TeacherNotFound(t *testing.T) {
	randomTeacherID := uuid.New()

	_, err := repository.FindByID(context.Background(), randomTeacherID)
	require.NotNil(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, codes.NotFound, convertedError.Code())
	assert.Equal(t, teacher.ErrorTeacherNotFoundByID.Error(), convertedError.Message())
}

func TestFindByID_HappyPath(t *testing.T) {
	teacher := randomTeacher()
	err := repository.Save(context.Background(), teacher)
	require.Nil(t, err)

	foundTeacher, err := repository.FindByID(context.Background(), teacher.ID)
	require.Nil(t, err)

	assert.Equal(t, teacher.ID, foundTeacher.ID)
	assert.Equal(t, teacher.FirstName, foundTeacher.FirstName)
	assert.Equal(t, teacher.LastName, foundTeacher.LastName)
	assert.Equal(t, teacher.MiddleName, foundTeacher.MiddleName)
	assert.Equal(t, teacher.ReportEmail, foundTeacher.ReportEmail)
	assert.Equal(t, teacher.Username, foundTeacher.Username)
}

func TestCheckDuplicates_DuplicatesNotExists(t *testing.T) {
	teacher := randomTeacher()
	duplicatesExists, err := repository.CheckDuplicateExists(context.Background(), teacher.ReportEmail, teacher.Username)
	require.Nil(t, err)
	assert.False(t, duplicatesExists)
}

func TestCheckDuplicates_DuplicatesExists(t *testing.T) {
	teacher := randomTeacher()
	err := repository.Save(context.Background(), teacher)
	require.Nil(t, err)

	duplicatesExists, err := repository.CheckDuplicateExists(context.Background(), teacher.ReportEmail, teacher.Username)
	require.Nil(t, err)

	assert.True(t, duplicatesExists)
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

func randomTeacher() domain.Teacher {
	return domain.Teacher{
		ID:          uuid.New(),
		FirstName:   gofakeit.FirstName(),
		LastName:    gofakeit.LastName(),
		MiddleName:  gofakeit.MiddleName(),
		ReportEmail: gofakeit.Email(),
		Username:    gofakeit.Username(),
	}
}
