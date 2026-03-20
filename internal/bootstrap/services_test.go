package bootstrap

import (
	"sync"
	"testing"
	"time"

	"gaussgo/internal/contracts"
)

func TestEventBusSubscribeEmitUnsubscribe(t *testing.T) {
	bus := NewInMemoryEventBus()

	called := 0
	id, err := bus.Subscribe("quiz.started", func(event contracts.Event) {
		called++
	})
	if err != nil {
		t.Fatalf("subscribe failed: %v", err)
	}

	if err := bus.Emit(contracts.Event{Name: "quiz.started"}); err != nil {
		t.Fatalf("emit failed: %v", err)
	}
	if called != 1 {
		t.Fatalf("expected handler called once, got %d", called)
	}

	if err := bus.Unsubscribe(id); err != nil {
		t.Fatalf("unsubscribe failed: %v", err)
	}
	if err := bus.Emit(contracts.Event{Name: "quiz.started"}); err != nil {
		t.Fatalf("emit failed after unsubscribe: %v", err)
	}
	if called != 1 {
		t.Fatalf("expected no new calls after unsubscribe, got %d", called)
	}
}

func TestEventBusEmptyNameErrors(t *testing.T) {
	bus := NewInMemoryEventBus()
	if err := bus.Emit(contracts.Event{}); err == nil {
		t.Fatal("expected emit error for empty name")
	}
}

func TestEventBusEmitSetsTimestamp(t *testing.T) {
	bus := NewInMemoryEventBus()
	var got contracts.Event
	_, err := bus.Subscribe("state.loaded", func(event contracts.Event) {
		got = event
	})
	if err != nil {
		t.Fatalf("subscribe failed: %v", err)
	}
	if err := bus.Emit(contracts.Event{Name: "state.loaded"}); err != nil {
		t.Fatalf("emit failed: %v", err)
	}
	if got.Timestamp.IsZero() {
		t.Fatal("expected timestamp to be auto-populated")
	}
}

func TestEventBusConcurrentEmitSafe(t *testing.T) {
	bus := NewInMemoryEventBus()
	_, err := bus.Subscribe("mods.refreshed", func(event contracts.Event) {})
	if err != nil {
		t.Fatalf("subscribe failed: %v", err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = bus.Emit(contracts.Event{Name: "mods.refreshed", Timestamp: time.Now()})
		}()
	}
	wg.Wait()

	if err := bus.Unsubscribe("unknown-id"); err != nil {
		t.Fatalf("expected idempotent unsubscribe, got %v", err)
	}
}
