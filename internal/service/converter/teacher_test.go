package converter_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	"github.com/upassed/upassed-account-service/internal/service/converter"
	business "github.com/upassed/upassed-account-service/internal/service/model"
)

var _ = Describe("Converter Tests", func() {
	Describe("Convert business-level Teacher to domain-level Teacher", func() {
		It("should convert Teacher properly", func() {
			teacher := business.Teacher{
				ID:          uuid.New(),
				FirstName:   gofakeit.FirstName(),
				LastName:    gofakeit.LastName(),
				MiddleName:  gofakeit.MiddleName(),
				ReportEmail: gofakeit.Email(),
				Username:    gofakeit.Username(),
			}

			domainTeacher := converter.ConvertTeacherToDomain(teacher)

			Expect(domainTeacher.ID).To(Equal(teacher.ID))
			Expect(domainTeacher.FirstName).To(Equal(teacher.FirstName))
			Expect(domainTeacher.LastName).To(Equal(teacher.LastName))
			Expect(domainTeacher.MiddleName).To(Equal(teacher.MiddleName))
			Expect(domainTeacher.ReportEmail).To(Equal(teacher.ReportEmail))
			Expect(domainTeacher.Username).To(Equal(teacher.Username))
		})
	})

	Describe("Convert domain-level Teacher to business-level Teacher", func() {
		It("should convert Teacher properly", func() {
			teacher := domain.Teacher{
				ID:          uuid.New(),
				FirstName:   gofakeit.FirstName(),
				LastName:    gofakeit.LastName(),
				MiddleName:  gofakeit.MiddleName(),
				ReportEmail: gofakeit.Email(),
				Username:    gofakeit.Username(),
			}

			businessTeacher := converter.ConvertTeacherToBusiness(teacher)

			Expect(businessTeacher.ID).To(Equal(teacher.ID))
			Expect(businessTeacher.FirstName).To(Equal(teacher.FirstName))
			Expect(businessTeacher.LastName).To(Equal(teacher.LastName))
			Expect(businessTeacher.MiddleName).To(Equal(teacher.MiddleName))
			Expect(businessTeacher.ReportEmail).To(Equal(teacher.ReportEmail))
			Expect(businessTeacher.Username).To(Equal(teacher.Username))
		})
	})
})

func TestTeacherConverter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Service Layer Teacher Converter Suite")
}
