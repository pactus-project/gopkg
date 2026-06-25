package scheduler

import (
	"context"
	"log"
	"time"

	"golang.org/x/sync/errgroup"
)

// Scheduler manages a collection of jobs and runs them on a configured interval.
type Scheduler struct {
	jobs      []Job
	onSuccess func()
}

// Option is a functional option for configuring the Scheduler.
type Option func(*Scheduler)

// NewScheduler creates a new empty Scheduler instance.
func NewScheduler() Scheduler {
	return Scheduler{
		jobs: make([]Job, 0),
	}
}

// WithOnSuccess registers a callback to run after all jobs succeed in a tick.
func WithOnSuccess(cb func()) Option {
	return func(s *Scheduler) {
		s.onSuccess = cb
	}
}

// AddJob registers a new job to be executed on each scheduler tick.
func (s *Scheduler) AddJob(job Job) {
	s.jobs = append(s.jobs, job)
}

// Start starts the scheduler and runs the jobs on the given interval.
func (s *Scheduler) Start(ctx context.Context, interval time.Duration, opts ...Option) {
	for _, opt := range opts {
		opt(s)
	}

	Every(interval).Do(ctx, func(ctx context.Context) {
		s.runJobs(ctx)
	})
}

func (s *Scheduler) runJobs(ctx context.Context) {
	group, _ := errgroup.WithContext(ctx)

	for _, j := range s.jobs {
		job := j
		group.Go(func() error {
			if err := job.Run(ctx); err != nil {
				log.Printf("job failed: %v", err)

				return err
			}

			return nil
		})
	}

	err := group.Wait()
	if err == nil && s.onSuccess != nil {
		s.onSuccess()
	}
}
