package types

import "time"

type Event struct {
	Name      string
	Timestamp time.Time
	Source    string
	Data      map[string]any
}
