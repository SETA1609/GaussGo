package mods

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDiscoverFindsAndValidates(t *testing.T) {
	root := t.TempDir()

	coreDir := filepath.Join(root, "core")
	if err := os.MkdirAll(coreDir, 0o755); err != nil {
		t.Fatalf("mkdir core: %v", err)
	}
	if err := os.WriteFile(filepath.Join(coreDir, "manifest.json"), []byte(`{
		"id":"core",
		"name":"Core",
		"version":"0.1.0",
		"entry":"public",
		"localization":{"default":"en","supported":["en"],"path":"i18n"},
		"dependencies":[],
		"provides":[]
	}`), 0o644); err != nil {
		t.Fatalf("write core manifest: %v", err)
	}

	laDir := filepath.Join(root, "linearAlgebra")
	if err := os.MkdirAll(laDir, 0o755); err != nil {
		t.Fatalf("mkdir linearAlgebra: %v", err)
	}
	if err := os.WriteFile(filepath.Join(laDir, "manifest.json"), []byte(`{
		"id":"linearAlgebra",
		"name":"Linear Algebra",
		"version":"0.1.0",
		"entry":"data/units",
		"localization":{"default":"en","supported":["en","es"],"path":"i18n"},
		"dependencies":[{"id":"core","version":">=0.1.0"}],
		"provides":["learning.units"]
	}`), 0o644); err != nil {
		t.Fatalf("write linearAlgebra manifest: %v", err)
	}

	manifests, err := Discover(root)
	if err != nil {
		t.Fatalf("discover failed: %v", err)
	}
	if len(manifests) != 2 {
		t.Fatalf("expected 2 manifests, got %d", len(manifests))
	}
}

func TestDiscoverFailsOnIDDirectoryMismatch(t *testing.T) {
	root := t.TempDir()

	coreDir := filepath.Join(root, "core")
	if err := os.MkdirAll(coreDir, 0o755); err != nil {
		t.Fatalf("mkdir core: %v", err)
	}
	if err := os.WriteFile(filepath.Join(coreDir, "manifest.json"), []byte(`{
		"id":"core",
		"name":"Core",
		"version":"0.1.0",
		"entry":"public",
		"localization":{"default":"en","supported":["en"],"path":"i18n"},
		"dependencies":[],
		"provides":[]
	}`), 0o644); err != nil {
		t.Fatalf("write core manifest: %v", err)
	}

	wrongDir := filepath.Join(root, "linearAlgebra")
	if err := os.MkdirAll(wrongDir, 0o755); err != nil {
		t.Fatalf("mkdir linearAlgebra: %v", err)
	}
	if err := os.WriteFile(filepath.Join(wrongDir, "manifest.json"), []byte(`{
		"id":"different",
		"name":"Linear Algebra",
		"version":"0.1.0",
		"entry":"data/units",
		"localization":{"default":"en","supported":["en"],"path":"i18n"},
		"dependencies":[{"id":"core","version":">=0.1.0"}],
		"provides":[]
	}`), 0o644); err != nil {
		t.Fatalf("write mismatch manifest: %v", err)
	}

	if _, err := Discover(root); err == nil {
		t.Fatal("expected mismatch id/directory error")
	}
}

func TestDiscoverFailsWhenCoreMissing(t *testing.T) {
	root := t.TempDir()
	laDir := filepath.Join(root, "linearAlgebra")
	if err := os.MkdirAll(laDir, 0o755); err != nil {
		t.Fatalf("mkdir linearAlgebra: %v", err)
	}
	if err := os.WriteFile(filepath.Join(laDir, "manifest.json"), []byte(`{
		"id":"linearAlgebra",
		"name":"Linear Algebra",
		"version":"0.1.0",
		"entry":"data/units",
		"localization":{"default":"en","supported":["en","es"],"path":"i18n"},
		"dependencies":[],
		"provides":["learning.units"]
	}`), 0o644); err != nil {
		t.Fatalf("write linearAlgebra manifest: %v", err)
	}

	if _, err := Discover(root); err == nil {
		t.Fatal("expected missing core error")
	}
}

func TestDiscoverFailsWhenModsDirMissing(t *testing.T) {
	if _, err := Discover(filepath.Join(t.TempDir(), "missing")); err == nil {
		t.Fatal("expected missing directory error")
	}
}
