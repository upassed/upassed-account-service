package converter_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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

			domainTeacher := converter.ConvertTeacher(teacher)

			Expect(domainTeacher.ID).To(Equal(teacher.ID))
			Expect(domainTeacher.FirstName).To(Equal(teacher.FirstName))
			Expect(domainTeacher.LastName).To(Equal(teacher.LastName))
			Expect(domainTeacher.MiddleName).To(Equal(teacher.MiddleName))
			Expect(domainTeacher.ReportEmail).To(Equal(teacher.ReportEmail))
			Expect(domainTeacher.Username).To(Equal(teacher.Username))
		})
	})
})

func TestTeacherConverter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Service Layer Teacher Converter Suite")
}
