// Package signal provides OS signal handling utilities for graceful shutdown.
package signal

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const defaultShutdownTimeout = 5 * time.Second

type options struct {
	shutdownTimeout time.Duration
}

// Option configures interrupt handling behavior.
type Option func(*options)

// WithShutdownTimeout sets how long the shutdown callback may take before the
// process exits. The default is 5 seconds.
func WithShutdownTimeout(d time.Duration) Option {
	return func(o *options) {
		o.shutdownTimeout = d
	}
}

// HandleInterrupt sets up signal handling for graceful shutdown on SIGINT and SIGTERM.
// The callback is called with a context that carries the shutdown deadline.
// After the callback executes, the process exits with the appropriate Unix exit code.
func HandleInterrupt(callback func(ctx context.Context), opts ...Option) {
	cfg := options{shutdownTimeout: defaultShutdownTimeout}
	for _, opt := range opts {
		opt(&cfg)
	}

	HandleSignals(func(sig os.Signal) {
		if callback != nil {
			ctx, cancel := context.WithTimeout(context.Background(), cfg.shutdownTimeout)
			callback(ctx)
			cancel()
		}

		//nolint:revive // calls to os.Exit is acceptable here.
		os.Exit(exitCode(sig))
	}, syscall.SIGINT, syscall.SIGTERM)
}

// HandleSignals sets up signal handling for specified signals.
// The callback function will be called with the received signal when the process receives any of the specified signals.
func HandleSignals(callback func(os.Signal), signals ...os.Signal) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, signals...)

	go func() {
		sig := <-sigChan
		callback(sig)
	}()
}

// exitCode returns the Unix exit code following the 128 + signal_number convention.
func exitCode(sig os.Signal) int {
	code := 128
	switch sig {
	case syscall.SIGINT:
		code += int(syscall.SIGINT)
	case syscall.SIGTERM:
		code += int(syscall.SIGTERM)
	}

	return code
}
