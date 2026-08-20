// Package net provides context-aware helpers for dialing and listening on
// network connections, with functional options for common configurations.
package net

import (
	"context"
	"net"
	"time"
)

// ListenOption is a functional option for configuring NetworkListen.
type ListenOption func(*listenOptions)

type listenOptions struct {
	keepAlive time.Duration
}

// WithKeepAlive sets the keep-alive period for network connections
// accepted by the listener. It only applies to TCP listeners.
func WithKeepAlive(keepAlive time.Duration) ListenOption {
	return func(o *listenOptions) {
		o.keepAlive = keepAlive
	}
}

// NetworkListen creates a network listener on the given address.
// The context can be used to cancel the listen operation before it completes.
// Options can be used to configure the listener, such as WithKeepAlive.
func NetworkListen(ctx context.Context, network, address string,
	opts ...ListenOption,
) (net.Listener, error) {
	cfg := &listenOptions{}
	for _, opt := range opts {
		opt(cfg)
	}

	lc := net.ListenConfig{
		KeepAlive: cfg.keepAlive,
	}

	return lc.Listen(ctx, network, address)
}

// DialOption is a functional option for configuring NetworkDial.
type DialOption func(*dialOptions)

type dialOptions struct {
	timeout time.Duration
}

// WithTimeout sets the maximum amount of time a dial will wait for a
// connection to complete.
func WithTimeout(timeout time.Duration) DialOption {
	return func(o *dialOptions) {
		o.timeout = timeout
	}
}

// NetworkDial establishes a connection to the given network address,
// honoring context cancellation. Options can be used to configure the dial,
// such as WithTimeout.
func NetworkDial(ctx context.Context, network, address string,
	opts ...DialOption,
) (net.Conn, error) {
	cfg := &dialOptions{}
	for _, opt := range opts {
		opt(cfg)
	}

	d := net.Dialer{
		Timeout: cfg.timeout,
	}

	return d.DialContext(ctx, network, address)
}

// NetworkDialTimeout dials a network address with a context and a timeout.
// It is a convenience wrapper around NetworkDial with WithTimeout.
func NetworkDialTimeout(ctx context.Context, network, address string,
	timeout time.Duration,
) (net.Conn, error) {
	return NetworkDial(ctx, network, address, WithTimeout(timeout))
}
