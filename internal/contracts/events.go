package contracts

import "time"

type Event struct {
	Name      string
	Timestamp time.Time
	Source    string
	Data      map[string]any
}

type EventHandler func(Event)

type EventBusService interface {
	Emit(event Event) error
	Subscribe(eventName string, handler EventHandler) (subscriptionID string, err error)
	Unsubscribe(subscriptionID string) error
}
