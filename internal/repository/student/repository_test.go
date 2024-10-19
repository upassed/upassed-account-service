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
	"github.com/upassed/upassed-account-service/internal/logging"
	testcontainer "github.com/upassed/upassed-account-service/internal/repository"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	"github.com/upassed/upassed-account-service/internal/repository/student"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type studentRepository interface {
	Save(context.Context, domain.Student) error
	FindByID(context.Context, uuid.UUID) (domain.Student, error)
	CheckDuplicateExists(ctx context.Context, educationalEmail, username string) (bool, error)
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

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("unable to parse cfg: ", err)
	}

	ctx := context.Background()
	container, err := testcontainer.NewPostgresTestcontainer(ctx)
	if err != nil {
		log.Fatal("unable to create a testcontainer: ", err)
	}

	port, err := container.Start(ctx)
	if err != nil {
		log.Fatal("unable to get a postgres testcontainer real port: ", err)
	}

	cfg.Storage.Port = strconv.Itoa(port)
	logger := logging.New(cfg.Env)
	if err := container.Migrate(cfg, logger); err != nil {
		log.Fatal("unable to run migrations: ", err)
	}

	repository, err = student.New(cfg, logger)
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
	cfg, err := config.Load()
	require.Nil(t, err)

	cfg.Storage.DatabaseName = "invalid-db-name"
	_, err = student.New(cfg, logging.New(cfg.Env))
	require.NotNil(t, err)
	assert.ErrorIs(t, err, student.ErrOpeningDbConnection)
}

func TestSave_InvalidUsernameLength(t *testing.T) {
	studentToSave := randomStudent()
	studentToSave.Username = gofakeit.LoremIpsumSentence(50)

	err := repository.Save(context.Background(), studentToSave)
	require.NotNil(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, codes.Internal, convertedError.Code())
	assert.Equal(t, student.ErrSavingStudent.Error(), convertedError.Message())
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
	assert.Equal(t, student.ErrStudentNotFoundByID.Error(), convertedError.Message())
}

func TestFindByID_HappyPath(t *testing.T) {
	existingStudent := randomStudent()
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
	assert.Equal(t, "5130904", foundStudent.Group.SpecializationCode)
	assert.Equal(t, "10101", foundStudent.Group.GroupNumber)
}

func TestCheckDuplicates_DuplicatesNotExists(t *testing.T) {
	uniqueStudent := randomStudent()
	duplicatesExists, err := repository.CheckDuplicateExists(context.Background(), uniqueStudent.EducationalEmail, uniqueStudent.Username)
	require.Nil(t, err)
	assert.False(t, duplicatesExists)
}

func TestCheckDuplicates_DuplicatesExists(t *testing.T) {
	duplicateStudent := randomStudent()
	groupID := uuid.MustParse("5eead8d5-b868-4708-aa25-713ad8399233")
	duplicateStudent.GroupID = groupID
	duplicateStudent.Group.ID = groupID

	err := repository.Save(context.Background(), duplicateStudent)
	require.Nil(t, err)

	duplicatesExists, err := repository.CheckDuplicateExists(context.Background(), duplicateStudent.EducationalEmail, duplicateStudent.Username)
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
