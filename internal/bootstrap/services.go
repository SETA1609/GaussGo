package bootstrap

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"gaussgo/internal/apperrors"
	"gaussgo/internal/contracts"
)

type StdLogger struct{}

func NewStdLogger() contracts.LoggingService {
	return StdLogger{}
}

func NewSilentLogger() contracts.LoggingService {
	return silentLogger{}
}

func (StdLogger) Debug(msg string, fields map[string]any) { logLine("debug", msg, fields) }
func (StdLogger) Info(msg string, fields map[string]any)  { logLine("info", msg, fields) }
func (StdLogger) Warn(msg string, fields map[string]any)  { logLine("warn", msg, fields) }
func (StdLogger) Error(msg string, fields map[string]any) { logLine("error", msg, fields) }

func logLine(level string, msg string, fields map[string]any) {
	_, _ = fmt.Fprintf(os.Stderr, "level=%s msg=%q fields=%v\n", level, msg, fields)
}

type silentLogger struct{}

func (silentLogger) Debug(string, map[string]any) {}
func (silentLogger) Info(string, map[string]any)  {}
func (silentLogger) Warn(string, map[string]any)  {}
func (silentLogger) Error(string, map[string]any) {}

func resolveDefaultLogger() contracts.LoggingService {
	if envBool(os.Getenv("GAUSSGO_QUIET")) || envBool(os.Getenv("GAUSSGO_SILENT")) {
		return NewSilentLogger()
	}
	level := strings.ToLower(strings.TrimSpace(os.Getenv("GAUSSGO_LOG_LEVEL")))
	if level == "" || level == "silent" || level == "off" || level == "none" {
		return NewSilentLogger()
	}
	if level == "debug" || level == "info" || level == "warn" || level == "error" {
		return NewStdLogger()
	}
	return NewStdLogger()
}

func envBool(value string) bool {
	v := strings.TrimSpace(strings.ToLower(value))
	return v == "1" || v == "true" || v == "yes" || v == "on"
}

type InMemoryEventBus struct {
	mu            sync.RWMutex
	subscribers   map[string]contracts.EventHandler
	eventToSubIDs map[string][]string
	nextID        int
}

func NewInMemoryEventBus() contracts.EventBusService {
	return &InMemoryEventBus{
		subscribers:   map[string]contracts.EventHandler{},
		eventToSubIDs: map[string][]string{},
	}
}

func (b *InMemoryEventBus) Emit(event contracts.Event) error {
	if event.Name == "" {
		return apperrors.New(apperrors.CodeValidation, apperrors.ErrorTypeInput, "event name is required")
	}
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now().UTC()
	}

	b.mu.RLock()
	ids := append([]string(nil), b.eventToSubIDs[event.Name]...)
	handlers := make([]contracts.EventHandler, 0, len(ids))
	for _, id := range ids {
		h, ok := b.subscribers[id]
		if ok {
			handlers = append(handlers, h)
		}
	}
	b.mu.RUnlock()

	for _, h := range handlers {
		h(event)
	}

	return nil
}

func (b *InMemoryEventBus) Subscribe(eventName string, handler contracts.EventHandler) (string, error) {
	if eventName == "" {
		return "", apperrors.New(apperrors.CodeValidation, apperrors.ErrorTypeInput, "event name is required")
	}
	if handler == nil {
		return "", apperrors.New(apperrors.CodeValidation, apperrors.ErrorTypeInput, "event handler is required")
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	b.nextID++
	id := fmt.Sprintf("sub-%d", b.nextID)
	b.subscribers[id] = handler
	b.eventToSubIDs[eventName] = append(b.eventToSubIDs[eventName], id)

	return id, nil
}

func (b *InMemoryEventBus) Unsubscribe(subscriptionID string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	delete(b.subscribers, subscriptionID)
	for eventName, ids := range b.eventToSubIDs {
		filtered := ids[:0]
		for _, id := range ids {
			if id != subscriptionID {
				filtered = append(filtered, id)
			}
		}
		if len(filtered) == 0 {
			delete(b.eventToSubIDs, eventName)
			continue
		}
		b.eventToSubIDs[eventName] = filtered
	}

	return nil
}
