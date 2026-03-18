package controllers

import (
	"time"

	"gaussgo/internal/contracts"
)

const (
	EventStateLoaded   = "state.loaded"
	EventStateCreated  = "state.created"
	EventLocaleChanged = "locale.changed"
	EventModsRefreshed = "mods.refreshed"
	EventModsToggled   = "mods.toggled"
)

func emitEvent(bus contracts.EventBusService, logger contracts.LoggingService, name string, source string, data map[string]any) {
	if bus == nil {
		return
	}
	err := bus.Emit(contracts.Event{
		Name:      name,
		Timestamp: time.Now().UTC(),
		Source:    source,
		Data:      data,
	})
	if err != nil && logger != nil {
		logger.Warn("event emit failed", map[string]any{"name": name, "source": source, "error": err.Error()})
	}
}
