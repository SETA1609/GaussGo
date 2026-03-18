package mods

import (
	"os"
	"path/filepath"
	"testing"

	"gaussgo/internal/contracts"
)

func TestRegistryRejectsDuplicateID(t *testing.T) {
	r := NewRegistry()
	m := contracts.ModManifest{ID: "core", Name: "Core", Version: "0.1.0", Entry: "public", Localization: contracts.ModLocalization{Default: "en", Supported: []string{"en"}, Path: "i18n"}}
	if err := r.Register(m); err != nil {
		t.Fatalf("register failed: %v", err)
	}
	if err := r.Register(m); err == nil {
		t.Fatal("expected duplicate registration error")
	}
}

func TestRepositoryRefreshIdempotent(t *testing.T) {
	root := t.TempDir()

	coreDir := filepath.Join(root, "core")
	if err := os.MkdirAll(coreDir, 0o755); err != nil {
		t.Fatalf("mkdir core: %v", err)
	}
	if err := os.WriteFile(filepath.Join(coreDir, "manifest.json"), []byte(`{
		"id":"core","name":"Core","version":"0.1.0","entry":"public",
		"localization":{"default":"en","supported":["en"],"path":"i18n"},
		"dependencies":[],"provides":[]}`), 0o644); err != nil {
		t.Fatalf("write core manifest: %v", err)
	}

	repo := NewRepository(root, map[string]bool{"core": true})
	first, err := repo.Refresh()
	if err != nil {
		t.Fatalf("first refresh failed: %v", err)
	}
	second, err := repo.Refresh()
	if err != nil {
		t.Fatalf("second refresh failed: %v", err)
	}
	if len(first) != len(second) {
		t.Fatalf("expected same status count, got %d and %d", len(first), len(second))
	}
}
