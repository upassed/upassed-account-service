package handling

import (
	"time"

	"google.golang.org/grpc/codes"
)

type ApplicationError interface {
	Error() string
}

type ApplicationErrorImpl struct {
	Message string
	Code    codes.Code
	Time    time.Time
}

func NewApplicationError(message string, code codes.Code) ApplicationError {
	return &ApplicationErrorImpl{
		Message: message,
		Code:    code,
		Time:    time.Now(),
	}
}

func (err *ApplicationErrorImpl) Error() string {
	return err.Message
}
