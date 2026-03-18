package bootstrap

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"gaussgo/internal/apperrors"
	"gaussgo/internal/contracts"
	"gaussgo/internal/state"
)

func TestBootstrapValidModsDir(t *testing.T) {
	root := t.TempDir()
	writeManifest(t, filepath.Join(root, "core"), `{
		"id":"core","name":"Core","version":"0.1.0","entry":"public",
		"localization":{"default":"en","supported":["en"],"path":"i18n"},
		"dependencies":[],"provides":[]}`)
	writeManifest(t, filepath.Join(root, "linearAlgebra"), `{
		"id":"linearAlgebra","name":"Linear Algebra","version":"0.1.0","entry":"data/units",
		"localization":{"default":"en","supported":["en"],"path":"i18n"},
		"dependencies":[{"id":"core","version":">=0.1.0"}],"provides":[]}`)

	d, err := Bootstrap(root, DefaultServices())
	if err != nil {
		t.Fatalf("bootstrap failed: %v", err)
	}
	if len(d.DiscoveredMods) != 2 {
		t.Fatalf("expected 2 discovered mods, got %d", len(d.DiscoveredMods))
	}
}

func TestBootstrapMissingCore(t *testing.T) {
	root := t.TempDir()
	writeManifest(t, filepath.Join(root, "linearAlgebra"), `{
		"id":"linearAlgebra","name":"Linear Algebra","version":"0.1.0","entry":"data/units",
		"localization":{"default":"en","supported":["en"],"path":"i18n"},
		"dependencies":[],"provides":[]}`)

	if _, err := Bootstrap(root, DefaultServices()); err == nil {
		t.Fatal("expected bootstrap error when core missing")
	}
}

func TestBootstrapInvalidManifest(t *testing.T) {
	root := t.TempDir()
	writeManifest(t, filepath.Join(root, "core"), `{
		"id":"core","name":"Core","version":"0.1.0","entry":"public",
		"localization":{"default":"en","supported":["en"],"path":"i18n"},
		"dependencies":[],"provides":[]}`)
	writeManifest(t, filepath.Join(root, "linearAlgebra"), `{"id":"linearAlgebra"}`)

	if _, err := Bootstrap(root, DefaultServices()); err == nil {
		t.Fatal("expected bootstrap error for invalid manifest")
	}
}

func TestBootstrapRuntimeLoadsActiveState(t *testing.T) {
	original := newStoreForBootstrapRuntime
	t.Cleanup(func() { newStoreForBootstrapRuntime = original })

	modsRoot := t.TempDir()
	writeManifest(t, filepath.Join(modsRoot, "core"), `{
		"id":"core","name":"Core","version":"0.1.0","entry":"public",
		"localization":{"default":"en","supported":["en","es"],"path":"i18n"},
		"dependencies":[],"provides":[]}`)

	statesRoot := t.TempDir()
	store := state.NewStore(statesRoot)
	created, err := store.Create("tester")
	if err != nil {
		t.Fatalf("create state failed: %v", err)
	}
	if err := store.SetActive(created.StateID); err != nil {
		t.Fatalf("set active failed: %v", err)
	}

	runtime, err := BootstrapRuntime(modsRoot, statesRoot, "main_menu", DefaultServices())
	if err != nil {
		t.Fatalf("bootstrap runtime failed: %v", err)
	}
	if runtime.Context.ActiveStateID != created.StateID {
		t.Fatalf("expected active state %s, got %s", created.StateID, runtime.Context.ActiveStateID)
	}
	if runtime.Controllers.Scene.Current() != "main_menu" {
		t.Fatalf("expected scene main_menu, got %s", runtime.Controllers.Scene.Current())
	}
}

func TestBootstrapRuntimeFailsWithoutActiveState(t *testing.T) {
	original := newStoreForBootstrapRuntime
	t.Cleanup(func() { newStoreForBootstrapRuntime = original })

	modsRoot := t.TempDir()
	writeManifest(t, filepath.Join(modsRoot, "core"), `{
		"id":"core","name":"Core","version":"0.1.0","entry":"public",
		"localization":{"default":"en","supported":["en","es"],"path":"i18n"},
		"dependencies":[],"provides":[]}`)

	statesRoot := t.TempDir()
	if _, err := BootstrapRuntime(modsRoot, statesRoot, "main_menu", DefaultServices()); err == nil {
		t.Fatal("expected bootstrap runtime error when active state is missing")
	} else {
		var appErr apperrors.Error
		if !errors.As(err, &appErr) {
			t.Fatalf("expected apperrors.Error, got %T", err)
		}
		if appErr.Code != apperrors.CodeNotFound {
			t.Fatalf("expected not_found code, got %s", appErr.Code)
		}
	}
}

func TestBootstrapRuntimePreservesAppErrorClassification(t *testing.T) {
	original := newStoreForBootstrapRuntime
	t.Cleanup(func() { newStoreForBootstrapRuntime = original })

	modsRoot := t.TempDir()
	writeManifest(t, filepath.Join(modsRoot, "core"), `{
		"id":"core","name":"Core","version":"0.1.0","entry":"public",
		"localization":{"default":"en","supported":["en","es"],"path":"i18n"},
		"dependencies":[],"provides":[]}`)

	store := &failingStateStore{err: apperrors.New(apperrors.CodeValidation, apperrors.ErrorTypeInput, "parse state file")}

	stateDir := t.TempDir()
	if err := os.MkdirAll(stateDir, 0o755); err != nil {
		t.Fatalf("mkdir state dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(stateDir, "active.txt"), []byte("default\n"), 0o600); err != nil {
		t.Fatalf("write active pointer: %v", err)
	}

	services := DefaultServices()
	newStoreForBootstrapRuntime = func(string) contracts.StateStore { return store }

	_, err := BootstrapRuntime(modsRoot, stateDir, "main_menu", services)
	if err == nil {
		t.Fatal("expected bootstrap runtime error")
	}

	var appErr apperrors.Error
	if !errors.As(err, &appErr) {
		t.Fatalf("expected apperrors.Error, got %T", err)
	}
	if appErr.Code != apperrors.CodeValidation {
		t.Fatalf("expected validation code preserved, got %s", appErr.Code)
	}
}

type failingStateStore struct {
	err error
}

func (s *failingStateStore) Create(string) (contracts.State, error) { return contracts.State{}, s.err }
func (s *failingStateStore) Load(string) (contracts.State, error)   { return contracts.State{}, s.err }
func (s *failingStateStore) Save(contracts.State) error             { return s.err }
func (s *failingStateStore) SetActive(string) error                 { return s.err }
func (s *failingStateStore) GetActive() (string, error)             { return "default", nil }

func writeManifest(t *testing.T, dir string, content string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", dir, err)
	}
	if err := os.WriteFile(filepath.Join(dir, "manifest.json"), []byte(content), 0o644); err != nil {
		t.Fatalf("write manifest %s: %v", dir, err)
	}
}
