package mods

import (
	"testing"

	"gaussgo/internal/contracts"
)

func TestBuildStatusesCoreAlwaysEnabled(t *testing.T) {
	manifests := []contracts.ModManifest{{ID: "core", Version: "0.1.0"}}
	statuses, err := BuildStatuses(manifests, map[string]bool{"core": false})
	if err != nil {
		t.Fatalf("build statuses failed: %v", err)
	}
	if len(statuses) != 1 || !statuses[0].Enabled {
		t.Fatalf("expected core enabled, got %+v", statuses)
	}
}

func TestBuildStatusesDependencyMissingDisablesMod(t *testing.T) {
	manifests := []contracts.ModManifest{
		{ID: "core", Version: "0.1.0"},
		{ID: "linearAlgebra", Version: "0.1.0", Dependencies: []contracts.ModDependency{{ID: "core", Version: ">=0.1.0"}}},
	}
	statuses, err := BuildStatuses(manifests, map[string]bool{})
	if err != nil {
		t.Fatalf("build statuses failed: %v", err)
	}
	if len(statuses) != 2 {
		t.Fatalf("expected 2 statuses, got %d", len(statuses))
	}
	if statuses[1].Available {
		t.Fatal("expected dependent mod unavailable when dependency missing from enabled map")
	}
}
