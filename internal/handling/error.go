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

func HandleApplicationError(err error, options ...Option) error {
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

	return WrapAsApplicationError(err, options...)
}

func WrapAsApplicationError(err error, options ...Option) error {
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
