package converter_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/upassed/upassed-account-service/internal/server/converter"
	business "github.com/upassed/upassed-account-service/internal/service/model"
	"github.com/upassed/upassed-account-service/pkg/client"
)

var _ = Describe("Converter Tests", func() {
	Describe("Convert TeacherCreateRequest to business-level Teacher model", func() {
		It("should convert TeacherCreateRequest properly", func() {
			request := client.TeacherCreateRequest{
				FirstName:   gofakeit.FirstName(),
				LastName:    gofakeit.LastName(),
				MiddleName:  gofakeit.MiddleName(),
				ReportEmail: gofakeit.Email(),
				Username:    gofakeit.Username(),
			}

			convertedTeacher := converter.ConvertTeacherCreateRequest(&request)

			Expect(convertedTeacher.ID).NotTo(BeNil())
			Expect(convertedTeacher.FirstName).To(Equal(request.GetFirstName()))
			Expect(convertedTeacher.LastName).To(Equal(request.GetLastName()))
			Expect(convertedTeacher.MiddleName).To(Equal(request.GetMiddleName()))
			Expect(convertedTeacher.ReportEmail).To(Equal(request.GetReportEmail()))
			Expect(convertedTeacher.Username).To(Equal(request.GetUsername()))
		})
	})

	Describe("Convert business-level TeacherCreateResponse to gRPC-level TeacherCreateResponse", func() {
		It("should convert TeacherCreateResponse properly", func() {
			response := business.TeacherCreateResponse{
				CreatedTeacherID: uuid.New(),
			}

			convertedResponse := converter.ConvertTeacherCreateResponse(response)

			Expect(convertedResponse.GetCreatedTeacherId()).To(Equal(response.CreatedTeacherID.String()))
		})
	})

	Describe("Convert business-level Teacher to gRPC-level TeacherFindByIDResponse", func() {
		It("should convert Teacher properly", func() {
			teacher := business.Teacher{
				ID:          uuid.New(),
				FirstName:   gofakeit.FirstName(),
				LastName:    gofakeit.LastName(),
				MiddleName:  gofakeit.MiddleName(),
				ReportEmail: gofakeit.Email(),
				Username:    gofakeit.Username(),
			}

			response := converter.ConvertTeacher(teacher)

			Expect(response.GetTeacher()).NotTo(BeNil())
			Expect(response.GetTeacher().GetFirstName()).To(Equal(teacher.FirstName))
			Expect(response.GetTeacher().GetLastName()).To(Equal(teacher.LastName))
			Expect(response.GetTeacher().GetMiddleName()).To(Equal(teacher.MiddleName))
			Expect(response.GetTeacher().GetReportEmail()).To(Equal(teacher.ReportEmail))
			Expect(response.GetTeacher().GetUsername()).To(Equal(teacher.Username))
		})
	})
})

func TestTeacherConverter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Server Layer Teacher Converter Suite")
}
