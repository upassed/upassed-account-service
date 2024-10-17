package student_test

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
	"github.com/upassed/upassed-account-service/internal/config"
	"github.com/upassed/upassed-account-service/internal/logger"
	testcontainer "github.com/upassed/upassed-account-service/internal/repository"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	"github.com/upassed/upassed-account-service/internal/repository/student"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type studentRepository interface {
	Save(context.Context, domain.Student) error
	FindByID(context.Context, uuid.UUID) (domain.Student, error)
	CheckDuplicateExists(ctx context.Context, edicationalEmail, username string) (bool, error)
}

var (
	repository studentRepository
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

	repository, err = student.New(config, logger)
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
	_, err = student.New(config, logger.New(config.Env))
	require.NotNil(t, err)
	assert.ErrorIs(t, err, student.ErrorOpeningDbConnection)
}

func TestSave_InvalidUsernameLength(t *testing.T) {
	studentToSave := randomStudent()
	studentToSave.Username = gofakeit.LoremIpsumSentence(50)

	err := repository.Save(context.Background(), studentToSave)
	require.NotNil(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, codes.Internal, convertedError.Code())
	assert.Equal(t, student.ErrorSavingStudent.Error(), convertedError.Message())
}

func TestSave_HappyPath(t *testing.T) {
	studentToSave := randomStudent()
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
	assert.Equal(t, student.ErrorStudentNotFound.Error(), convertedError.Message())
}

func TestFindByID_HappyPath(t *testing.T) {
	student := randomStudent()
	groupID := uuid.MustParse("5eead8d5-b868-4708-aa25-713ad8399233")
	student.GroupID = groupID
	student.Group.ID = groupID

	err := repository.Save(context.Background(), student)
	require.Nil(t, err)

	foundStudent, err := repository.FindByID(context.Background(), student.ID)
	require.Nil(t, err)

	assert.Equal(t, student.ID, foundStudent.ID)
	assert.Equal(t, student.FirstName, foundStudent.FirstName)
	assert.Equal(t, student.LastName, foundStudent.LastName)
	assert.Equal(t, student.MiddleName, foundStudent.MiddleName)
	assert.Equal(t, student.EducationalEmail, foundStudent.EducationalEmail)
	assert.Equal(t, student.Username, foundStudent.Username)
	assert.Equal(t, student.GroupID, foundStudent.GroupID)
	assert.Equal(t, student.Group.ID, foundStudent.Group.ID)
	assert.Equal(t, "5130904", foundStudent.Group.SpecializationCode)
	assert.Equal(t, "10101", foundStudent.Group.GroupNumber)
}

func TestCheckDuplicates_DuplicatesNotExists(t *testing.T) {
	student := randomStudent()
	duplicatesExists, err := repository.CheckDuplicateExists(context.Background(), student.EducationalEmail, student.Username)
	require.Nil(t, err)
	assert.False(t, duplicatesExists)
}

func TestCheckDuplicates_DuplicatesExists(t *testing.T) {
	student := randomStudent()
	groupID := uuid.MustParse("5eead8d5-b868-4708-aa25-713ad8399233")
	student.GroupID = groupID
	student.Group.ID = groupID

	err := repository.Save(context.Background(), student)
	require.Nil(t, err)

	duplicatesExists, err := repository.CheckDuplicateExists(context.Background(), student.EducationalEmail, student.Username)
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

func randomStudent() domain.Student {
	return domain.Student{
		ID:               uuid.New(),
		FirstName:        gofakeit.FirstName(),
		LastName:         gofakeit.LastName(),
		MiddleName:       gofakeit.MiddleName(),
		EducationalEmail: gofakeit.Email(),
		Username:         gofakeit.Username(),
		Group: domain.Group{
			ID:                 uuid.New(),
			SpecializationCode: gofakeit.WeekDay(),
			GroupNumber:        gofakeit.WeekDay(),
		},
	}
}
