package handling_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/upassed/upassed-account-service/internal/handling"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Handling Application Errors Tests", func() {
	Describe("Convert Application Error", func() {
		It("should convert ApplicationError and add a custom additional details", func() {
			message := gofakeit.Error().Error()
			code := codes.Internal
			applicationError := handling.NewApplicationError(message, code)

			handledError := handling.HandleApplicationError(applicationError)

			st := status.Convert(handledError)
			Expect(st.Code()).To(Equal(code))
			Expect(st.Message()).To(Equal(message))
			Expect(st.Details()).To(HaveLen(1))
		})
	})

	Describe("Convert not an Application Error", func() {
		It("should wrap not an ApplicationError and add a custom additional details", func() {
			initialError := gofakeit.Error()

			handledError := handling.HandleApplicationError(initialError)

			st := status.Convert(handledError)
			Expect(st.Code()).To(Equal(codes.Internal))
			Expect(st.Message()).To(Equal(initialError.Error()))
			Expect(st.Details()).To(HaveLen(1))
		})
	})
})

func TestHandling(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Handling Application Errors Suite")
}
