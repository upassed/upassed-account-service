package async

import (
	"context"
	"golang.org/x/sync/errgroup"
	"time"
)

func ExecuteWithTimeout[T any](parentContext context.Context, timeout time.Duration, callable func(ctx context.Context) (T, error)) (T, error) {
	ctxWithTimeout, cancel := context.WithTimeout(parentContext, timeout)
	defer cancel()

	errGroup, errGroupContext := errgroup.WithContext(ctxWithTimeout)
	resultChannel := make(chan T)

	var callableResult T
	errGroup.Go(func() error {
		result, err := callable(errGroupContext)
		if err != nil {
			return err
		}

		select {
		case resultChannel <- result:
			return nil
		case <-errGroupContext.Done():
			return errGroupContext.Err()
		}
	})

	select {
	case <-ctxWithTimeout.Done():
		return callableResult, ctxWithTimeout.Err()
	case result := <-resultChannel:
		callableResult = result
	case <-errGroupContext.Done():
		if err := errGroup.Wait(); err != nil {
			return callableResult, err
		}
	}

	if err := errGroup.Wait(); err != nil {
		return callableResult, err
	}

	return callableResult, nil
}
