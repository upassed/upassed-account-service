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
	Describe("ConvertTeacherCreateRequest", func() {
		It("should convert TeacherCreateRequest properly", func() {
			request := client.TeacherCreateRequest{
				FirstName:   gofakeit.FirstName(),
				LastName:    gofakeit.LastName(),
				MiddleName:  gofakeit.MiddleName(),
				ReportEmail: gofakeit.Email(),
				Username:    gofakeit.Username(),
			}

			convertedRequest := converter.ConvertTeacherCreateRequest(&request)

			Expect(convertedRequest.FirstName).To(Equal(request.GetFirstName()))
			Expect(convertedRequest.LastName).To(Equal(request.GetLastName()))
			Expect(convertedRequest.MiddleName).To(Equal(request.GetMiddleName()))
			Expect(convertedRequest.ReportEmail).To(Equal(request.GetReportEmail()))
			Expect(convertedRequest.Username).To(Equal(request.GetUsername()))
		})
	})

	Describe("ConvertTeacherCreateResponse", func() {
		It("should convert TeacherCreateResponse properly", func() {
			response := business.TeacherCreateResponse{
				CreatedTeacherID: uuid.New(),
			}

			convertedResponse := converter.TestConvertTeacherCreateResponse(response)

			Expect(convertedResponse.GetCreatedTeacherId()).To(Equal(response.CreatedTeacherID.String()))
		})
	})
})

func TestTeacherConverter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Server Layer Teacher Converter Suite")
}
