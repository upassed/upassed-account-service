package handling

import (
	"errors"
	"fmt"
	"time"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ServiceLayerError interface {
	Error() string
}

type ServiceLayerErrorImpl struct {
	Message string
	Code    codes.Code
	Time    time.Time
}

func NewServiceLayerError(message string, code codes.Code) ServiceLayerError {
	return &ServiceLayerErrorImpl{
		Message: message,
		Code:    code,
		Time:    time.Now(),
	}
}

func (err *ServiceLayerErrorImpl) Error() string {
	return err.Message
}

const timeFormat string = "2006-01-02 15:04:05"

func HandleServiceLayerError(err error) error {
	var serviceLayerErr *ServiceLayerErrorImpl
	if errors.As(err, &serviceLayerErr) {
		convertedErr := status.New(serviceLayerErr.Code, serviceLayerErr.Message)
		timeInfo := errdetails.DebugInfo{
			Detail: fmt.Sprintf("Time: %s", serviceLayerErr.Time.Format(timeFormat)),
		}

		convertedErrWithDetails, err := convertedErr.WithDetails(&timeInfo)
		if err != nil {
			return convertedErr.Err()
		}

		return convertedErrWithDetails.Err()
	}

	return err
}
