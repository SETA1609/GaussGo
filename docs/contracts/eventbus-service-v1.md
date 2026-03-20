# EventBus Service Contract v1

## Version
- `eventbus-service/v1`

## Required Interface and Invariants
- Event bus supports publishing and subscribing to named events.
- Subscribers receive events asynchronously or synchronously per implementation, but ordering is preserved per event name in a single emitter flow.
- Event payload uses an envelope with `name`, `timestamp`, `source`, and `data`.
- Event handlers should not panic the bus; panic-safe dispatch is required.

## Validation Rules
- Reject empty event names.
- Reject nil handlers on subscribe.
- Unsubscribe by token/ID must be idempotent.

## Example Interface (Go)
```go
type Event struct {
    Name      string
    Timestamp time.Time
    Source    string
    Data      map[string]any
}

type EventBusService interface {
    Emit(event Event) error
    Subscribe(eventName string, handler func(Event)) (subscriptionID string, err error)
    Unsubscribe(subscriptionID string) error
}
```

## Example Events
- `state.loaded`
- `state.saved`
- `locale.changed`
- `mods.refreshed`
- `quiz.started`
- `quiz.finished`

## Backward Compatibility Notes
- Adding new event names is non-breaking.
- Changing event envelope required fields is breaking and requires a major contract revision.
