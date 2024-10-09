package handling

import (
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ApplicationErrorImpl struct {
	Message string
	Code    codes.Code
	Time    time.Time
}

func NewApplicationError(message string, code codes.Code) *ApplicationErrorImpl {
	return &ApplicationErrorImpl{
		Message: message,
		Code:    code,
		Time:    time.Now(),
	}
}

func (err *ApplicationErrorImpl) Error() string {
	return err.Message
}

func (err *ApplicationErrorImpl) GRPCStatus() *status.Status {
	return status.New(err.Code, err.Message)
}
