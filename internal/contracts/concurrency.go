package contracts

import (
	"context"
	"time"
)

// Task is a unit of work that accepts a context and returns an error.
// Use closures to capture inputs and write outputs — this keeps the interface
// generic without requiring interface{} values.
type Task func(ctx context.Context) error

// TaskResult pairs a task's position in the input slice with its error.
// Only failing tasks produce a TaskResult; successful tasks are omitted.
type TaskResult struct {
	Index int
	Err   error
}

// ConcurrencyService abstracts parallel task execution and timeout control.
//
// It is intended for use by framework internals and mod authors who need
// structured parallelism without writing WaitGroup or channel boilerplate.
// Inject it via RuntimeServices; never import the concurrency package directly
// from mod code.
type ConcurrencyService interface {
	// RunParallel runs all tasks concurrently under a shared worker budget and
	// waits for every task to complete or for ctx to be cancelled. It returns
	// one TaskResult per failing task, sorted by task index.
	RunParallel(ctx context.Context, tasks []Task) []TaskResult

	// RunWithTimeout runs task under a deadline derived from ctx and timeout.
	// The tighter of the two deadlines wins. The task's context is cancelled
	// when the timeout elapses or when the parent ctx is cancelled.
	RunWithTimeout(ctx context.Context, timeout time.Duration, task Task) error
}
