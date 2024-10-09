package repository_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	config "github.com/upassed/upassed-account-service/internal/config/app"
	"github.com/upassed/upassed-account-service/internal/logger"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	repository "github.com/upassed/upassed-account-service/internal/repository/postgres"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Teacher Repository Tests", Ordered, func() {
	BeforeAll(func() {
		projectRoot, err := getProjectRoot()
		Expect(err).To(BeNil())

		err = os.Setenv(config.EnvConfigPath, filepath.Join(projectRoot, "config", "app", "test.yml"))
		Expect(err).To(BeNil())
	})

	Describe("Save teacher tests", func() {
		It("should return internal status error if some check is not satisfied (ex: username is too long)", func() {
			config, err := config.Load()
			Expect(err).To(BeNil())

			log := logger.New(config.Env)

			teacherToSave := randomTeacher()
			teacherToSave.Username = gofakeit.LoremIpsumSentence(50)

			repo, err := repository.NewTeacherRepository(config, log)
			Expect(err).To(BeNil())

			err = repo.Save(context.Background(), teacherToSave)
			Expect(err).NotTo(BeNil())

			convertedError := status.Convert(err)
			Expect(convertedError.Code()).To(Equal(codes.Internal))
			Expect(convertedError.Message()).To(Equal(repository.ErrorSavingTeacher.Error()))
		})

		It("should not return error if the teacher was successfully saved to a database", func() {
			config, err := config.Load()
			Expect(err).To(BeNil())

			log := logger.New(config.Env)

			teacherToSave := randomTeacher()

			repo, err := repository.NewTeacherRepository(config, log)
			Expect(err).To(BeNil())

			err = repo.Save(context.Background(), teacherToSave)
			Expect(err).To(BeNil())
		})
	})

	Describe("Find teacher by id tests", func() {
		It("should return not found status error if teacher was not found in database by id", func() {
			config, err := config.Load()
			Expect(err).To(BeNil())

			log := logger.New(config.Env)

			randomTeacherID := uuid.New()

			repo, err := repository.NewTeacherRepository(config, log)
			Expect(err).To(BeNil())

			_, err = repo.FindByID(context.Background(), randomTeacherID)
			Expect(err).NotTo(BeNil())

			convertedError := status.Convert(err)
			Expect(convertedError.Code()).To(Equal(codes.NotFound))
			Expect(convertedError.Message()).To(Equal(repository.ErrorTeacherNotFound.Error()))
		})

		It("should return teacher data if the teacher was found by id", func() {
			config, err := config.Load()
			Expect(err).To(BeNil())

			log := logger.New(config.Env)

			repo, err := repository.NewTeacherRepository(config, log)
			Expect(err).To(BeNil())

			teacher := randomTeacher()
			err = repo.Save(context.Background(), teacher)
			Expect(err).To(BeNil())

			foundTeacher, err := repo.FindByID(context.Background(), teacher.ID)
			Expect(err).To(BeNil())

			Expect(teacher.ID).To(Equal(foundTeacher.ID))
			Expect(teacher.FirstName).To(Equal(foundTeacher.FirstName))
			Expect(teacher.LastName).To(Equal(foundTeacher.LastName))
			Expect(teacher.MiddleName).To(Equal(foundTeacher.MiddleName))
			Expect(teacher.ReportEmail).To(Equal(foundTeacher.ReportEmail))
			Expect(teacher.Username).To(Equal(foundTeacher.Username))
		})
	})

	Describe("Check teacher duplicates tests", func() {
		It("should return false if duplicates does not exists in database", func() {
			config, err := config.Load()
			Expect(err).To(BeNil())

			log := logger.New(config.Env)

			repo, err := repository.NewTeacherRepository(config, log)
			Expect(err).To(BeNil())

			teacher := randomTeacher()
			duplicatesExists, err := repo.CheckDuplicateExists(context.Background(), teacher.ReportEmail, teacher.Username)
			Expect(err).To(BeNil())
			Expect(duplicatesExists).To(BeFalse())
		})

		It("should return true if duplicates exists in database", func() {
			config, err := config.Load()
			Expect(err).To(BeNil())

			log := logger.New(config.Env)

			repo, err := repository.NewTeacherRepository(config, log)
			Expect(err).To(BeNil())

			teacher := randomTeacher()
			err = repo.Save(context.Background(), teacher)
			Expect(err).To(BeNil())

			duplicatesExists, err := repo.CheckDuplicateExists(context.Background(), teacher.ReportEmail, teacher.Username)
			Expect(err).To(BeNil())
			Expect(duplicatesExists).To(BeTrue())
		})
	})
})

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

func TestTeacherRepository(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Teacher Repository Test Suite")
}
