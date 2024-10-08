package handling

import (
	"errors"
	"fmt"
	"time"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

const timeFormat string = "2006-01-02 15:04:05"

func HandleApplicationError(err error) error {
	var applicationErr *ApplicationErrorImpl
	if errors.As(err, &applicationErr) {
		convertedErr := status.New(applicationErr.Code, applicationErr.Message)
		timeInfo := errdetails.DebugInfo{
			Detail: fmt.Sprintf("Time: %s", applicationErr.Time.Format(timeFormat)),
		}

		convertedErrWithDetails, err := convertedErr.WithDetails(&timeInfo)
		if err != nil {
			return convertedErr.Err()
		}

		return convertedErrWithDetails.Err()
	}

	return wrapAsApplicationError(err)
}

func wrapAsApplicationError(err error) ApplicationError {
	convertedErr := status.New(codes.Internal, err.Error())
	timeInfo := errdetails.DebugInfo{
		Detail: fmt.Sprintf("Time: %s", time.Now().Format(timeFormat)),
	}

	convertedErrWithDetails, err := convertedErr.WithDetails(&timeInfo)
	if err != nil {
		return convertedErr.Err()
	}

	return convertedErrWithDetails.Err()
}
