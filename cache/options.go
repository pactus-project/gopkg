package cache

import "time"

var defaultConfig = options{
	cleanUpInterval: 10 * time.Second,
}

type options struct {
	cleanUpInterval time.Duration
}

// WithCleanUpInterval sets the interval at which expired entries are purged.
func WithCleanUpInterval(interval time.Duration) Option {
	return func(cfg *options) {
		cfg.cleanUpInterval = interval
	}
}

// Option is a functional option for configuring a Cache.
type Option func(*options)
