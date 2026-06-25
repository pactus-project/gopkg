package retry

import (
	"context"
	"time"
)

type (
	// SyncTask is a function that returns an error, used for synchronous retry.
	SyncTask         func() error
	// SyncTaskT is a generic function that returns a value and an error, used for synchronous retry.
	SyncTaskT[T any] func() (T, error)
)

// Options is a functional option for configuring sync execution.
type Options func(*syncOptions)

type syncOptions struct {
	maxRetries int
	retryDelay time.Duration
}

func defaultSyncOpts() *syncOptions {
	return &syncOptions{
		maxRetries: 3,
		retryDelay: 2 * time.Second,
	}
}

// WithSyncMaxRetries sets the maximum number of retry attempts for sync tasks.
func WithSyncMaxRetries(maxRetries int) Options {
	return func(o *syncOptions) {
		o.maxRetries = maxRetries
	}
}

// WithSyncRetryDelay sets the delay between retry attempts for sync tasks.
func WithSyncRetryDelay(retryDelay time.Duration) Options {
	return func(o *syncOptions) {
		o.retryDelay = retryDelay
	}
}

// ExecuteSync executes a function synchronously with retry logic
// It respects context cancellation and timeout
// Returns nil if the function succeeds, or the last error if all retries are exhausted.
func ExecuteSync(ctx context.Context,
	task SyncTask,
	opts ...Options,
) error {
	_, err := ExecuteSyncT(ctx, func() (any, error) {
		return nil, task()
	}, opts...)

	return err
}

// ExecuteSyncT executes a function synchronously with retry logic and returns a result
// It respects context cancellation and timeout
// Returns the result if the function succeeds, or the last error if all retries are exhausted.
//
func ExecuteSyncT[T any](ctx context.Context,
	task SyncTaskT[T], opts ...Options,
) (T, error) {
	conf := defaultSyncOpts()
	for _, opt := range opts {
		opt(conf)
	}

	var result T
	var err error
	for attempt := 0; attempt < conf.maxRetries; attempt++ {
		result, err = task()
		if err == nil {
			return result, nil
		}

		// Don't wait after the last attempt
		if attempt < conf.maxRetries-1 {
			// Wait before retry, but respect context cancellation
			select {
			case <-ctx.Done():
				return result, ctx.Err()

			case <-time.After(conf.retryDelay):
				// Continue to next retry
			}
		}
	}

	// All retries exhausted
	return result, err
}
