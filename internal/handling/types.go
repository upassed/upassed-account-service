package handling

import (
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type applicationErrorImpl struct {
	Message string
	Code    codes.Code
	Time    time.Time
}

func NewApplicationError(message string, code codes.Code) *applicationErrorImpl {
	return &applicationErrorImpl{
		Message: message,
		Code:    code,
		Time:    time.Now(),
	}
}

func (err *applicationErrorImpl) Error() string {
	return err.Message
}

func (err *applicationErrorImpl) GRPCStatus() *status.Status {
	return status.New(err.Code, err.Message)
}
