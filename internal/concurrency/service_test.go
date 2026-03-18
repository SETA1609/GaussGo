package concurrency

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"gaussgo/internal/contracts"
)

func TestRunParallelAllSucceed(t *testing.T) {
	svc := NewService(4)
	tasks := []contracts.Task{
		func(ctx context.Context) error { return nil },
		func(ctx context.Context) error { return nil },
		func(ctx context.Context) error { return nil },
	}
	if errs := svc.RunParallel(context.Background(), tasks); len(errs) != 0 {
		t.Fatalf("expected no errors, got %+v", errs)
	}
}

func TestRunParallelCollectsAllErrors(t *testing.T) {
	svc := NewService(4)
	errA := errors.New("task a failed")
	errC := errors.New("task c failed")
	tasks := []contracts.Task{
		func(ctx context.Context) error { return errA },
		func(ctx context.Context) error { return nil },
		func(ctx context.Context) error { return errC },
	}
	results := svc.RunParallel(context.Background(), tasks)
	if len(results) != 2 {
		t.Fatalf("expected 2 errors, got %d", len(results))
	}
	if results[0].Index != 0 || !errors.Is(results[0].Err, errA) {
		t.Fatalf("unexpected result[0]: %+v", results[0])
	}
	if results[1].Index != 2 || !errors.Is(results[1].Err, errC) {
		t.Fatalf("unexpected result[1]: %+v", results[1])
	}
}

func TestRunParallelResultsSortedByIndex(t *testing.T) {
	svc := NewService(1) // serial execution to get deterministic completion order
	sentinel := errors.New("fail")
	tasks := make([]contracts.Task, 5)
	for i := range tasks {
		tasks[i] = func(ctx context.Context) error { return sentinel }
	}
	results := svc.RunParallel(context.Background(), tasks)
	for i, r := range results {
		if r.Index != i {
			t.Fatalf("expected sorted index %d, got %d", i, r.Index)
		}
	}
}

func TestRunParallelEmpty(t *testing.T) {
	svc := NewService(4)
	if errs := svc.RunParallel(context.Background(), nil); errs != nil {
		t.Fatalf("expected nil for empty tasks, got %+v", errs)
	}
}

func TestRunParallelRespectsContextCancellation(t *testing.T) {
	svc := NewService(4)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // already cancelled

	var ran int32
	tasks := []contracts.Task{
		func(ctx context.Context) error {
			atomic.AddInt32(&ran, 1)
			return nil
		},
		func(ctx context.Context) error {
			atomic.AddInt32(&ran, 1)
			return nil
		},
	}
	results := svc.RunParallel(ctx, tasks)
	if len(results) == 0 {
		t.Fatal("expected context cancellation errors")
	}
	for _, r := range results {
		if !errors.Is(r.Err, context.Canceled) {
			t.Fatalf("expected context.Canceled, got %v", r.Err)
		}
	}
	if atomic.LoadInt32(&ran) != 0 {
		t.Fatal("expected no tasks to execute after cancellation")
	}
}

func TestRunParallelBoundedConcurrency(t *testing.T) {
	const maxWorkers = 2
	svc := NewService(maxWorkers)

	var active int32
	var peak int32
	block := make(chan struct{})

	tasks := make([]contracts.Task, 6)
	for i := range tasks {
		tasks[i] = func(ctx context.Context) error {
			cur := atomic.AddInt32(&active, 1)
			for {
				p := atomic.LoadInt32(&peak)
				if cur <= p || atomic.CompareAndSwapInt32(&peak, p, cur) {
					break
				}
			}
			<-block
			atomic.AddInt32(&active, -1)
			return nil
		}
	}

	done := make(chan []contracts.TaskResult, 1)
	go func() {
		done <- svc.RunParallel(context.Background(), tasks)
	}()

	// Give goroutines time to reach the block channel.
	time.Sleep(20 * time.Millisecond)
	close(block)
	<-done

	if p := atomic.LoadInt32(&peak); p > maxWorkers {
		t.Fatalf("peak concurrent workers %d exceeded maxWorkers %d", p, maxWorkers)
	}
}

func TestRunWithTimeoutCompletes(t *testing.T) {
	svc := NewService(1)
	err := svc.RunWithTimeout(context.Background(), time.Second, func(ctx context.Context) error {
		return nil
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestRunWithTimeoutExceeds(t *testing.T) {
	svc := NewService(1)
	err := svc.RunWithTimeout(context.Background(), 10*time.Millisecond, func(ctx context.Context) error {
		<-ctx.Done()
		return ctx.Err()
	})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected DeadlineExceeded, got %v", err)
	}
}

func TestRunWithTimeoutParentCancels(t *testing.T) {
	svc := NewService(1)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()
	err := svc.RunWithTimeout(ctx, time.Hour, func(ctx context.Context) error {
		<-ctx.Done()
		return ctx.Err()
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected Canceled, got %v", err)
	}
}

func TestNewServicePanicsOnZeroWorkers(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for maxWorkers=0")
		}
	}()
	NewService(0)
}
