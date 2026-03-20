package controllers

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"gaussgo/internal/contracts"
	"gaussgo/internal/mods"
	"gaussgo/internal/types"
)

func TestModControllerDisableCoreRejected(t *testing.T) {
	modsDir, store, created := prepareStateAndMods(t)
	ctx := &types.RuntimeContext{ActiveStateID: created.StateID, EnabledMods: []string{"core", "linearAlgebra"}}
	c := NewModController(modsDir, store, ctx, noopLogger{}, &eventRecorder{}, nil)

	if err := c.Disable("core"); err == nil {
		t.Fatal("expected core disable rejection")
	}
}

func TestModControllerEnableDisablePersists(t *testing.T) {
	modsDir, store, created := prepareStateAndMods(t)
	recorder := &eventRecorder{}
	ctx := &types.RuntimeContext{ActiveStateID: created.StateID, EnabledMods: []string{"core"}}
	c := NewModController(modsDir, store, ctx, noopLogger{}, recorder, nil)

	if err := c.Enable("linearAlgebra"); err != nil {
		t.Fatalf("enable failed: %v", err)
	}
	reloaded, err := store.Load(created.StateID)
	if err != nil {
		t.Fatalf("load after enable failed: %v", err)
	}
	if !hasValue(reloaded.EnabledMods, "linearAlgebra") {
		t.Fatalf("expected enabled mods to include linearAlgebra, got %v", reloaded.EnabledMods)
	}

	if err := c.Disable("linearAlgebra"); err != nil {
		t.Fatalf("disable failed: %v", err)
	}
	reloaded, err = store.Load(created.StateID)
	if err != nil {
		t.Fatalf("load after disable failed: %v", err)
	}
	if hasValue(reloaded.EnabledMods, "linearAlgebra") {
		t.Fatalf("expected enabled mods to not include linearAlgebra, got %v", reloaded.EnabledMods)
	}

	if len(recorder.events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(recorder.events))
	}
}

func TestModControllerDisableClearsCurrentSelection(t *testing.T) {
	modsDir, store, created := prepareStateAndMods(t)
	ctx := &types.RuntimeContext{
		ActiveStateID:    created.StateID,
		EnabledMods:      []string{"core", "linearAlgebra"},
		CurrentModID:     "linearAlgebra",
		CurrentUnitID:    "vectors",
		CurrentConceptID: "vectors.dot-product",
	}
	c := NewModController(modsDir, store, ctx, noopLogger{}, &eventRecorder{}, nil)

	if err := c.Disable("linearAlgebra"); err != nil {
		t.Fatalf("disable failed: %v", err)
	}

	if ctx.CurrentModID != "" || ctx.CurrentUnitID != "" || ctx.CurrentConceptID != "" {
		t.Fatalf("expected selection cleared, got mod=%q unit=%q concept=%q", ctx.CurrentModID, ctx.CurrentUnitID, ctx.CurrentConceptID)
	}
}

func TestModControllerEnableRejectsUnknownMod(t *testing.T) {
	modsDir, store, created := prepareStateAndMods(t)
	ctx := &types.RuntimeContext{ActiveStateID: created.StateID, EnabledMods: []string{"core"}}
	c := NewModController(modsDir, store, ctx, noopLogger{}, &eventRecorder{}, nil)

	if err := c.Enable("not-a-real-mod"); err == nil {
		t.Fatal("expected unknown mod enable error")
	}
}

func TestModControllerRefreshReturnsStatuses(t *testing.T) {
	modsDir, store, created := prepareStateAndMods(t)
	writeManifest(t, filepath.Join(modsDir, "advancedLinear"), `{
		"id":"advancedLinear","name":"Advanced Linear","version":"0.1.0","entry":"data/units",
		"localization":{"default":"en","supported":["en","es"],"path":"i18n"},
		"dependencies":[{"id":"linearAlgebra","version":">=0.1.0"}],"provides":[]}`)
	recorder := &eventRecorder{}
	ctx := &types.RuntimeContext{ActiveStateID: created.StateID, EnabledMods: []string{"core"}}
	c := NewModController(modsDir, store, ctx, noopLogger{}, recorder, nil)

	statuses, err := c.Refresh()
	if err != nil {
		t.Fatalf("refresh failed: %v", err)
	}
	if len(statuses) != 3 {
		t.Fatalf("expected three mod statuses, got %d", len(statuses))
	}
	if len(recorder.events) != 1 || recorder.events[0].Name != EventModsRefreshed {
		t.Fatalf("expected one mods.refreshed event, got %#v", recorder.events)
	}
	payload := recorder.events[0].Data
	if payload["total"] != 3 {
		t.Fatalf("expected total=3, got %v", payload["total"])
	}
	if payload["valid"] != 2 {
		t.Fatalf("expected valid=2, got %v", payload["valid"])
	}
	if payload["invalid"] != 1 {
		t.Fatalf("expected invalid=1, got %v", payload["invalid"])
	}
}

func TestModControllerDisableUnloadsRuntimeMod(t *testing.T) {
	modsDir, store, created := prepareStateAndMods(t)
	runtimeServices := testRuntimeServices{logger: noopLogger{}, eventBus: &eventRecorder{}, stateStore: store}
	lifecycle := mods.NewRuntimeLifecycle(mods.ManifestFactory{}, runtimeServices, noopLogger{}, &eventRecorder{})

	manifests, err := mods.Discover(modsDir)
	if err != nil {
		t.Fatalf("discover failed: %v", err)
	}
	if err := lifecycle.SyncEnabled(context.Background(), manifests, types.RuntimeContext{}, map[string]bool{"core": true, "linearAlgebra": true}, nil); err != nil {
		t.Fatalf("sync enabled failed: %v", err)
	}

	ctx := &types.RuntimeContext{ActiveStateID: created.StateID, EnabledMods: []string{"core", "linearAlgebra"}}
	c := NewModController(modsDir, store, ctx, noopLogger{}, &eventRecorder{}, lifecycle)

	if err := c.Disable("linearAlgebra"); err != nil {
		t.Fatalf("disable failed: %v", err)
	}

	if _, err := lifecycle.Registry().Get("linearAlgebra"); err == nil {
		t.Fatal("expected runtime mod linearAlgebra to be unloaded")
	}
}

func TestModControllerEnableRollsBackStateOnRuntimeInitFailure(t *testing.T) {
	modsDir, store, created := prepareStateAndMods(t)
	recorder := &eventRecorder{}
	factory := &failingFactory{failModID: "linearAlgebra", err: errors.New("boom")}
	runtimeServices := testRuntimeServices{logger: noopLogger{}, eventBus: recorder, stateStore: store}
	lifecycle := mods.NewRuntimeLifecycle(factory, runtimeServices, noopLogger{}, recorder)

	manifests, err := mods.Discover(modsDir)
	if err != nil {
		t.Fatalf("discover failed: %v", err)
	}
	if err := lifecycle.SyncEnabled(context.Background(), manifests, types.RuntimeContext{}, map[string]bool{"core": true}, nil); err != nil {
		t.Fatalf("preload core failed: %v", err)
	}

	ctx := &types.RuntimeContext{ActiveStateID: created.StateID, EnabledMods: []string{"core"}}
	c := NewModController(modsDir, store, ctx, noopLogger{}, recorder, lifecycle)

	if err := c.Enable("linearAlgebra"); err == nil {
		t.Fatal("expected enable to fail when runtime init fails")
	}

	reloaded, err := store.Load(created.StateID)
	if err != nil {
		t.Fatalf("load state failed: %v", err)
	}
	if hasValue(reloaded.EnabledMods, "linearAlgebra") {
		t.Fatalf("expected rollback to remove linearAlgebra, got %v", reloaded.EnabledMods)
	}
	if _, err := lifecycle.Registry().Get("linearAlgebra"); err == nil {
		t.Fatal("expected runtime registry not to include failed mod")
	}
}

type testRuntimeServices struct {
	logger     contracts.LoggingService
	eventBus   contracts.EventBusService
	stateStore contracts.StateStore
}

func (s testRuntimeServices) Logger() contracts.LoggingService {
	return s.logger
}

func (s testRuntimeServices) EventBus() contracts.EventBusService {
	return s.eventBus
}

func (s testRuntimeServices) StateStore() contracts.StateStore {
	return s.stateStore
}

type failingFactory struct {
	failModID string
	err       error
}

func (f *failingFactory) Build(manifest contracts.ModManifest) (contracts.ModRuntime, error) {
	if manifest.ID == f.failModID {
		return nil, f.err
	}
	factory := mods.ManifestFactory{}
	return factory.Build(manifest)
}

func hasValue(in []string, want string) bool {
	for _, v := range in {
		if v == want {
			return true
		}
	}
	return false
}
