package controllers

import (
	"testing"

	"gaussgo/internal/types"
)

func TestLocaleControllerSetLocalePersists(t *testing.T) {
	_, store, created := prepareStateAndMods(t)
	recorder := &eventRecorder{}
	ctx := &types.RuntimeContext{ActiveStateID: created.StateID}

	c := NewLocaleController(store, ctx, noopLogger{}, recorder, map[string]struct{}{"en": {}, "es": {}})
	if err := c.SetLocale("es"); err != nil {
		t.Fatalf("set locale failed: %v", err)
	}

	reloaded, err := store.Load(created.StateID)
	if err != nil {
		t.Fatalf("reload failed: %v", err)
	}
	if reloaded.UI.Locale != "es" {
		t.Fatalf("expected persisted locale es, got %s", reloaded.UI.Locale)
	}
	if ctx.CurrentLocale != "es" {
		t.Fatalf("expected context locale es, got %s", ctx.CurrentLocale)
	}
	if len(recorder.events) != 1 || recorder.events[0].Name != EventLocaleChanged {
		t.Fatalf("expected one locale changed event, got %#v", recorder.events)
	}
}

func TestLocaleControllerSetLocaleRejectsUnsupported(t *testing.T) {
	_, store, created := prepareStateAndMods(t)
	ctx := &types.RuntimeContext{ActiveStateID: created.StateID}
	c := NewLocaleController(store, ctx, noopLogger{}, &eventRecorder{}, map[string]struct{}{"en": {}, "es": {}})
	if err := c.SetLocale("fr"); err == nil {
		t.Fatal("expected unsupported locale error")
	}
}
