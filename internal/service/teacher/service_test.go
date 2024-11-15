package teacher_test

import (
	"context"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/golang/mock/gomock"
	"github.com/upassed/upassed-account-service/internal/util"
	"github.com/upassed/upassed-account-service/internal/util/mocks"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/upassed/upassed-account-service/internal/config"
	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/internal/logging"
	"github.com/upassed/upassed-account-service/internal/service/teacher"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	cfg        *config.Config
	repository *mocks.TeacherRepository
	service    teacher.Service
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

	cfg, err = config.Load()
	if err != nil {
		log.Fatal("unable to parse config: ", err)
	}

	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()

	repository = mocks.NewTeacherRepository(ctrl)
	service = teacher.New(cfg, logging.New(config.EnvTesting), repository)

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestCreate_ErrorCheckingDuplicateExistsOccurred(t *testing.T) {
	duplicateTeacher := util.RandomBusinessTeacher()

	expectedRepoError := handling.New("repo layer error message", codes.Internal)
	repository.EXPECT().
		CheckDuplicateExists(gomock.Any(), duplicateTeacher.ReportEmail, duplicateTeacher.Username).
		Return(false, expectedRepoError)

	_, err := service.Create(context.Background(), duplicateTeacher)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedRepoError.Error(), convertedError.Message())
}

func TestCreate_DuplicateExists(t *testing.T) {
	duplicateTeacher := util.RandomBusinessTeacher()

	repository.EXPECT().
		CheckDuplicateExists(gomock.Any(), duplicateTeacher.ReportEmail, duplicateTeacher.Username).
		Return(true, nil)

	_, err := service.Create(context.Background(), duplicateTeacher)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, "teacher duplicate found", convertedError.Message())
	assert.Equal(t, codes.AlreadyExists, convertedError.Code())
}

func TestCreate_ErrorSavingToDatabase(t *testing.T) {
	teacherToSave := util.RandomBusinessTeacher()

	repository.EXPECT().
		CheckDuplicateExists(gomock.Any(), teacherToSave.ReportEmail, teacherToSave.Username).
		Return(false, nil)

	expectedRepoError := handling.New("repo layer error message", codes.DeadlineExceeded)
	repository.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Return(expectedRepoError)

	_, err := service.Create(context.Background(), teacherToSave)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedRepoError.Error(), convertedError.Message())
	assert.Equal(t, expectedRepoError.Code(), convertedError.Code())
}

func TestCreate_DeadlineExceeded(t *testing.T) {
	oldTimeout := cfg.Timeouts.EndpointExecutionTimeoutMS
	cfg.Timeouts.EndpointExecutionTimeoutMS = "0"

	teacherToSave := util.RandomBusinessTeacher()

	repository.EXPECT().
		CheckDuplicateExists(gomock.Any(), teacherToSave.ReportEmail, teacherToSave.Username).
		Return(false, nil)

	repository.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Return(nil)

	_, err := service.Create(context.Background(), teacherToSave)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, codes.DeadlineExceeded, convertedError.Code())

	cfg.Timeouts.EndpointExecutionTimeoutMS = oldTimeout
}

func TestCreate_HappyPath(t *testing.T) {
	teacherToSave := util.RandomBusinessTeacher()

	repository.EXPECT().
		CheckDuplicateExists(gomock.Any(), teacherToSave.ReportEmail, teacherToSave.Username).
		Return(false, nil)

	repository.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Return(nil)

	response, err := service.Create(context.Background(), teacherToSave)
	require.NoError(t, err)

	assert.Equal(t, teacherToSave.ID, response.CreatedTeacherID)
}

func TestFindByUsername_ErrorSearchingTeacherInDatabase(t *testing.T) {
	teacherUsername := gofakeit.Username()

	expectedRepoError := handling.New("repo layer error message", codes.NotFound)
	repository.EXPECT().
		FindByUsername(gomock.Any(), teacherUsername).
		Return(nil, expectedRepoError)

	_, err := service.FindByUsername(context.Background(), teacherUsername)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedRepoError.Code(), convertedError.Code())
	assert.Equal(t, expectedRepoError.Error(), convertedError.Message())
}

func TestFindByUsername_ErrorDeadlineExceeded(t *testing.T) {
	oldTimeout := cfg.Timeouts.EndpointExecutionTimeoutMS
	cfg.Timeouts.EndpointExecutionTimeoutMS = "0"

	teacherUsername := gofakeit.Username()
	expectedFoundTeacher := teacher.ConvertToRepositoryTeacher(util.RandomBusinessTeacher())

	repository.EXPECT().
		FindByUsername(gomock.Any(), teacherUsername).
		Return(expectedFoundTeacher, nil)

	_, err := service.FindByUsername(context.Background(), teacherUsername)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, codes.DeadlineExceeded, convertedError.Code())

	cfg.Timeouts.EndpointExecutionTimeoutMS = oldTimeout
}

func TestFindByUsername_HappyPath(t *testing.T) {
	teacherUsername := gofakeit.Username()
	expectedFoundTeacher := teacher.ConvertToRepositoryTeacher(util.RandomBusinessTeacher())

	repository.EXPECT().
		FindByUsername(gomock.Any(), teacherUsername).
		Return(expectedFoundTeacher, nil)

	foundTeacher, err := service.FindByUsername(context.Background(), teacherUsername)
	require.NoError(t, err)

	assert.Equal(t, teacher.ConvertToServiceTeacher(expectedFoundTeacher), foundTeacher)
}
