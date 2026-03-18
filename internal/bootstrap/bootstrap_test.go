package bootstrap

import (
	"os"
	"path/filepath"
	"testing"
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

func writeManifest(t *testing.T, dir string, content string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", dir, err)
	}
	if err := os.WriteFile(filepath.Join(dir, "manifest.json"), []byte(content), 0o644); err != nil {
		t.Fatalf("write manifest %s: %v", dir, err)
	}
}
