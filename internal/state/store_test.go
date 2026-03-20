package state

import (
	"os"
	"path/filepath"
	"testing"

	"gaussgo/internal/contracts"
)

func TestCreateInitializesCoreAndLocale(t *testing.T) {
	store := NewStore(t.TempDir())

	created, err := store.Create("Test User")
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	if created.StateID == "" {
		t.Fatal("expected non-empty state id")
	}
	if len(created.EnabledMods) == 0 || created.EnabledMods[0] != "core" {
		t.Fatalf("expected core enabled, got %+v", created.EnabledMods)
	}
	if created.UI.Locale != "en" {
		t.Fatalf("expected locale en, got %s", created.UI.Locale)
	}
}

func TestLoadMissingStateID(t *testing.T) {
	store := NewStore(t.TempDir())
	if _, err := store.Load("missing"); err == nil {
		t.Fatal("expected missing state load error")
	}
}

func TestLoadMalformedJSON(t *testing.T) {
	root := t.TempDir()
	store := NewStore(root)
	if err := os.WriteFile(filepath.Join(root, "state-bad.json"), []byte("{"), 0o644); err != nil {
		t.Fatalf("write malformed file: %v", err)
	}
	if _, err := store.Load("bad"); err == nil {
		t.Fatal("expected malformed json error")
	}
}

func TestLoadUnsupportedSchema(t *testing.T) {
	root := t.TempDir()
	store := NewStore(root)
	if err := os.WriteFile(filepath.Join(root, "state-bad.json"), []byte(`{
		"schemaVersion": 99,
		"stateId": "bad",
		"profileName": "Bad",
		"createdAt": "2026-03-18T00:00:00Z",
		"updatedAt": "2026-03-18T00:00:00Z",
		"modsAtCreation": [],
		"enabledMods": ["core"],
		"ui": {"locale": "en"},
		"progress": {}
	}`), 0o644); err != nil {
		t.Fatalf("write unsupported schema file: %v", err)
	}
	if _, err := store.Load("bad"); err == nil {
		t.Fatal("expected unsupported schema error")
	}
}

func TestLoadNormalizesCoreAndLocale(t *testing.T) {
	store := NewStore(t.TempDir())

	state := contracts.State{
		SchemaVersion: 1,
		StateID:       "n1",
		ProfileName:   "N",
		CreatedAt:     mustTime("2026-03-18T00:00:00Z"),
		UpdatedAt:     mustTime("2026-03-18T00:00:00Z"),
		EnabledMods:   []string{"linearAlgebra"},
		UI: contracts.UIState{
			Locale: "",
		},
		Progress: map[string]contracts.ModProgress{},
	}

	if err := store.Save(state); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	loaded, err := store.Load("n1")
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	if loaded.UI.Locale != "en" {
		t.Fatalf("expected locale normalized to en, got %s", loaded.UI.Locale)
	}
	if loaded.EnabledMods[0] != "core" {
		t.Fatalf("expected core first after normalization, got %+v", loaded.EnabledMods)
	}
}

func TestStateRoundtripPersistsProgressAndStats(t *testing.T) {
	store := NewStore(t.TempDir())

	state := contracts.State{
		SchemaVersion: 1,
		StateID:       "default",
		ProfileName:   "Default",
		CreatedAt:     mustTime("2026-03-18T00:00:00Z"),
		UpdatedAt:     mustTime("2026-03-18T00:00:00Z"),
		EnabledMods:   []string{"core", "linearAlgebra"},
		UI: contracts.UIState{
			Locale:      "en",
			Preferences: map[string]string{},
		},
		Progress: map[string]contracts.ModProgress{
			"linearAlgebra": {
				ReadConcepts: []string{"vectors.dot-product"},
				ExerciseStats: contracts.ExerciseStats{
					Attempted: 4,
					Correct:   3,
				},
				ExerciseHistory: []string{"q1", "q2"},
			},
		},
	}

	if err := store.Save(state); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	loaded, err := store.Load("default")
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	p := loaded.Progress["linearAlgebra"]
	if len(p.ReadConcepts) != 1 || p.ReadConcepts[0] != "vectors.dot-product" {
		t.Fatalf("read concepts not preserved: %+v", p.ReadConcepts)
	}
	if p.ExerciseStats.CorrectnessPercent != 75 {
		t.Fatalf("expected 75%% correctness, got %v", p.ExerciseStats.CorrectnessPercent)
	}
}
