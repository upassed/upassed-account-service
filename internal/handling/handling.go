package handling

import (
	"errors"
	"fmt"
	"time"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const timeFormat string = "2006-01-02 15:04:05"

type Option func(*wrapOptions)

type wrapOptions struct {
	code codes.Code
	time time.Time
}

type applicationError struct {
	Message string
	Code    codes.Code
	Time    time.Time
}

func New(message string, code codes.Code) *applicationError {
	return &applicationError{
		Message: message,
		Code:    code,
		Time:    time.Now(),
	}
}

func (err *applicationError) Error() string {
	return err.Message
}

func (err *applicationError) GRPCStatus() *status.Status {
	return status.New(err.Code, err.Message)
}

func Process(err error, options ...Option) error {
	var applicationErr *applicationError
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

	return Wrap(err, options...)
}

func Wrap(err error, options ...Option) error {
	opts := defaultOptions()
	for _, opt := range options {
		opt(opts)
	}

	convertedErr := status.New(opts.code, err.Error())
	timeInfo := errdetails.DebugInfo{
		Detail: fmt.Sprintf("Time: %s", opts.time.Format(timeFormat)),
	}

	convertedErrWithDetails, err := convertedErr.WithDetails(&timeInfo)
	if err != nil {
		return convertedErr.Err()
	}

	return convertedErrWithDetails.Err()
}

func defaultOptions() *wrapOptions {
	return &wrapOptions{
		code: codes.Internal,
		time: time.Now(),
	}
}

func WithCode(code codes.Code) Option {
	return func(opts *wrapOptions) {
		opts.code = code
	}
}
