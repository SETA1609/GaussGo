package concurrency

import (
	"context"
	"runtime"
	"sort"
	"sync"
	"time"

	"gaussgo/internal/contracts"
)

var _ contracts.ConcurrencyService = (*service)(nil)

type service struct {
	maxWorkers int
}

// NewService returns a ConcurrencyService with a fixed worker budget.
// maxWorkers controls how many tasks may run concurrently in RunParallel.
// Panics if maxWorkers < 1.
func NewService(maxWorkers int) contracts.ConcurrencyService {
	if maxWorkers < 1 {
		panic("concurrency: maxWorkers must be at least 1")
	}
	return &service{maxWorkers: maxWorkers}
}

// DefaultService returns a ConcurrencyService sized to the available CPU count.
func DefaultService() contracts.ConcurrencyService {
	return NewService(runtime.GOMAXPROCS(0))
}

// RunParallel runs every task in its own goroutine, bounded by the service's
// worker budget. It waits for all tasks to finish and returns the errors of
// every task that failed, sorted by task index.
//
// If ctx is already cancelled when RunParallel is called, every task receives
// a context-cancelled result immediately without executing.
func (s *service) RunParallel(ctx context.Context, tasks []contracts.Task) []contracts.TaskResult {
	if len(tasks) == 0 {
		return nil
	}

	results := make(chan contracts.TaskResult, len(tasks))
	sem := make(chan struct{}, s.maxWorkers)

	var wg sync.WaitGroup
	for i, task := range tasks {
		wg.Add(1)
		go func(idx int, t contracts.Task) {
			defer wg.Done()

			// Explicit pre-check: if the context is already done, skip the
			// select entirely so the task never executes.
			if err := ctx.Err(); err != nil {
				results <- contracts.TaskResult{Index: idx, Err: err}
				return
			}

			// Acquire a worker slot, still respecting cancellation in case
			// it fires while waiting for a free slot.
			select {
			case <-ctx.Done():
				results <- contracts.TaskResult{Index: idx, Err: ctx.Err()}
				return
			case sem <- struct{}{}:
			}
			defer func() { <-sem }()

			if err := t(ctx); err != nil {
				results <- contracts.TaskResult{Index: idx, Err: err}
			}
		}(i, task)
	}

	wg.Wait()
	close(results)

	var out []contracts.TaskResult
	for r := range results {
		out = append(out, r)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Index < out[j].Index
	})
	return out
}

// RunWithTimeout runs task under a deadline derived from ctx and timeout.
// The tighter of the two deadlines wins.
func (s *service) RunWithTimeout(ctx context.Context, timeout time.Duration, task contracts.Task) error {
	tctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return task(tctx)
}
