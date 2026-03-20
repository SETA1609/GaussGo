package state

import (
	"testing"

	"gaussgo/internal/contracts"
)

func TestNormalizeAndValidateSchemaDefaultsAndWarnings(t *testing.T) {
	in := contracts.State{
		StateID:     "s1",
		ProfileName: "P",
		EnabledMods: []string{"linearAlgebra"},
		UI:          contracts.UIState{},
	}

	out, warnings, err := NormalizeAndValidate(in)
	if err != nil {
		t.Fatalf("normalize failed: %v", err)
	}
	if out.SchemaVersion != 1 {
		t.Fatalf("expected schema version 1, got %d", out.SchemaVersion)
	}
	if len(warnings) == 0 {
		t.Fatal("expected normalization warnings")
	}
}

func TestNormalizeAndValidateSchemaMismatch(t *testing.T) {
	_, _, err := NormalizeAndValidate(contracts.State{
		SchemaVersion: 99,
		StateID:       "s1",
		ProfileName:   "P",
		EnabledMods:   []string{"core"},
		UI:            contracts.UIState{Locale: "en"},
	})
	if err == nil {
		t.Fatal("expected schema mismatch error")
	}
}

func TestNormalizeAndValidateStatsValidationAndDedupe(t *testing.T) {
	_, _, err := NormalizeAndValidate(contracts.State{
		SchemaVersion: 1,
		StateID:       "s1",
		ProfileName:   "P",
		EnabledMods:   []string{"core"},
		UI:            contracts.UIState{Locale: "en"},
		Progress: map[string]contracts.ModProgress{
			"linearAlgebra": {
				ExerciseStats: contracts.ExerciseStats{Attempted: 1, Correct: 2},
			},
		},
	})
	if err == nil {
		t.Fatal("expected correct > attempted error")
	}

	out, _, err := NormalizeAndValidate(contracts.State{
		SchemaVersion: 1,
		StateID:       "s2",
		ProfileName:   "P",
		EnabledMods:   []string{"core"},
		UI:            contracts.UIState{Locale: "en"},
		Progress: map[string]contracts.ModProgress{
			"linearAlgebra": {
				ReadConcepts: []string{"a", "a", " ", "b"},
				ExerciseStats: contracts.ExerciseStats{
					Attempted: 4,
					Correct:   3,
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("normalize failed: %v", err)
	}
	p := out.Progress["linearAlgebra"]
	if len(p.ReadConcepts) != 2 {
		t.Fatalf("expected deduped read concepts, got %+v", p.ReadConcepts)
	}
}

func TestNormalizeAndValidateMissingIdentity(t *testing.T) {
	_, _, err := NormalizeAndValidate(contracts.State{SchemaVersion: 1, ProfileName: "P", UI: contracts.UIState{Locale: "en"}})
	if err == nil {
		t.Fatal("expected missing stateId error")
	}

	_, _, err = NormalizeAndValidate(contracts.State{SchemaVersion: 1, StateID: "s1", ProfileName: "  ", UI: contracts.UIState{Locale: "en"}})
	if err == nil {
		t.Fatal("expected missing profileName error")
	}
}
