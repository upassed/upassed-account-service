package service_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	config "github.com/upassed/upassed-account-service/internal/config/app"
	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/internal/logger"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	"github.com/upassed/upassed-account-service/internal/service"
	business "github.com/upassed/upassed-account-service/internal/service/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mockTeacherRepository struct {
	mock.Mock
}

func (m *mockTeacherRepository) Save(ctx context.Context, teacher domain.Teacher) error {
	args := m.Called(ctx, teacher)
	return args.Error(0)
}

func (m *mockTeacherRepository) FindByID(ctx context.Context, teacherID uuid.UUID) (domain.Teacher, error) {
	args := m.Called(ctx, teacherID)
	return args.Get(0).(domain.Teacher), args.Error(1)
}

func (m *mockTeacherRepository) CheckDuplicateExists(ctx context.Context, reportEmail, username string) (bool, error) {
	args := m.Called(ctx, reportEmail, username)
	return args.Bool(0), args.Error(1)
}

var _ = Describe("Teacher Service Tests", func() {
	Describe("Create teacher tests", func() {
		It("should return internal status error if there was error checking duplicate teacher", func() {
			log := logger.New(config.EnvDev)
			repository := new(mockTeacherRepository)
			teacher := randomTeacher()

			expectedRepoError := handling.NewApplicationError("repo layer error message", codes.Internal)
			repository.On("CheckDuplicateExists", mock.Anything, teacher.ReportEmail, teacher.Username).Return(false, expectedRepoError)

			service := service.NewTeacherService(log, repository)

			_, err := service.Create(context.Background(), teacher)
			Expect(err).NotTo(BeNil())

			convertedError := status.Convert(err)
			Expect(convertedError.Message()).To(Equal(expectedRepoError.Error()))
			Expect(convertedError.Code()).To(Equal(expectedRepoError.Code))
		})

		It("should return not found status error if there was duplicate teacher by email or username", func() {
			log := logger.New(config.EnvDev)
			repository := new(mockTeacherRepository)
			teacher := randomTeacher()

			repository.On("CheckDuplicateExists", mock.Anything, teacher.ReportEmail, teacher.Username).Return(true, nil)

			service := service.NewTeacherService(log, repository)

			_, err := service.Create(context.Background(), teacher)
			Expect(err).NotTo(BeNil())

			convertedError := status.Convert(err)
			Expect(convertedError.Message()).To(Equal("teacher duplicate found"))
			Expect(convertedError.Code()).To(Equal(codes.AlreadyExists))
		})

		It("should return error if there was error saving a teacher to a database", func() {
			log := logger.New(config.EnvDev)
			repository := new(mockTeacherRepository)
			teacher := randomTeacher()

			repository.On("CheckDuplicateExists", mock.Anything, teacher.ReportEmail, teacher.Username).Return(false, nil)

			expectedRepoError := handling.NewApplicationError("repo layer error message", codes.DeadlineExceeded)
			repository.On("Save", mock.Anything, mock.Anything).Return(expectedRepoError)

			service := service.NewTeacherService(log, repository)

			_, err := service.Create(context.Background(), teacher)
			Expect(err).NotTo(BeNil())

			convertedError := status.Convert(err)
			Expect(convertedError.Message()).To(Equal(expectedRepoError.Error()))
			Expect(convertedError.Code()).To(Equal(expectedRepoError.Code))
		})

		It("should not return error if the teacher was successfully saved to a database", func() {
			log := logger.New(config.EnvDev)
			repository := new(mockTeacherRepository)
			teacher := randomTeacher()

			repository.On("CheckDuplicateExists", mock.Anything, teacher.ReportEmail, teacher.Username).Return(false, nil)
			repository.On("Save", mock.Anything, mock.Anything).Return(nil)

			service := service.NewTeacherService(log, repository)

			response, err := service.Create(context.Background(), teacher)
			Expect(err).To(BeNil())

			Expect(response.CreatedTeacherID).To(Equal(teacher.ID))
		})
	})
})

func randomTeacher() business.Teacher {
	return business.Teacher{
		ID:          uuid.New(),
		FirstName:   gofakeit.FirstName(),
		LastName:    gofakeit.LastName(),
		MiddleName:  gofakeit.MiddleName(),
		ReportEmail: gofakeit.Email(),
		Username:    gofakeit.Username(),
	}
}

func TestService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Teacher Service Test Suite")
}
