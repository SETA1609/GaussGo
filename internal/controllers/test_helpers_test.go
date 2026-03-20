package controllers

import (
	"os"
	"path/filepath"
	"testing"

	"gaussgo/internal/contracts"
	"gaussgo/internal/state"
)

type eventRecorder struct {
	events []contracts.Event
}

func (r *eventRecorder) Emit(event contracts.Event) error {
	r.events = append(r.events, event)
	return nil
}

func (r *eventRecorder) Subscribe(string, contracts.EventHandler) (string, error) {
	return "", nil
}

func (r *eventRecorder) Unsubscribe(string) error {
	return nil
}

type noopLogger struct{}

func (noopLogger) Debug(string, map[string]any) {}
func (noopLogger) Info(string, map[string]any)  {}
func (noopLogger) Warn(string, map[string]any)  {}
func (noopLogger) Error(string, map[string]any) {}

func prepareStateAndMods(t *testing.T) (modsDir string, store *state.Store, created contracts.State) {
	t.Helper()

	root := t.TempDir()
	modsDir = filepath.Join(root, "mods")
	statesDir := filepath.Join(root, "states")

	writeManifest(t, filepath.Join(modsDir, "core"), `{
		"id":"core","name":"Core","version":"0.1.0","entry":"public",
		"localization":{"default":"en","supported":["en","es"],"path":"i18n"},
		"dependencies":[],"provides":[]}`)
	writeManifest(t, filepath.Join(modsDir, "linearAlgebra"), `{
		"id":"linearAlgebra","name":"Linear Algebra","version":"0.1.0","entry":"data/units",
		"localization":{"default":"en","supported":["en","es"],"path":"i18n"},
		"dependencies":[{"id":"core","version":">=0.1.0"}],"provides":[]}`)

	store = state.NewStore(statesDir)
	created, err := store.Create("tester")
	if err != nil {
		t.Fatalf("create state failed: %v", err)
	}
	if err := store.SetActive(created.StateID); err != nil {
		t.Fatalf("set active failed: %v", err)
	}

	return modsDir, store, created
}

func writeManifest(t *testing.T, dir string, content string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", dir, err)
	}
	if err := os.WriteFile(filepath.Join(dir, "manifest.json"), []byte(content), 0o644); err != nil {
		t.Fatalf("write manifest %s: %v", dir, err)
	}
}
