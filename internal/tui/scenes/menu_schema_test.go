package scenes

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadMenuSchemasFromFiles(t *testing.T) {
	root := t.TempDir()
	coreScenes := filepath.Join(root, "core", "ui", "scenes")
	if err := os.MkdirAll(coreScenes, 0o755); err != nil {
		t.Fatalf("mkdir scenes: %v", err)
	}

	content := `{"titleKey":"scene.main_menu.title","title":"Main Menu","subtitleKey":"scene.main_menu.subtitle","subtitle":"Choose","lines":[],"items":[]}`
	if err := os.WriteFile(filepath.Join(coreScenes, "main_menu.json"), []byte(content), 0o644); err != nil {
		t.Fatalf("write schema: %v", err)
	}

	schemas, err := LoadMenuSchemas(root)
	if err != nil {
		t.Fatalf("load menu schemas failed: %v", err)
	}
	schema, ok := schemas[SceneMainMenu]
	if !ok {
		t.Fatalf("expected %s schema loaded", SceneMainMenu)
	}
	if schema.Title != "Main Menu" {
		t.Fatalf("expected title Main Menu, got %q", schema.Title)
	}
}
