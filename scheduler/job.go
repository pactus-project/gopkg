package scheduler

import "context"

// Job represents a schedulable task that can be run with context support.
type Job interface {
	Run(ctx context.Context) error
}
