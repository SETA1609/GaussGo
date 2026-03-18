package mods

import (
	"testing"

	"gaussgo/internal/contracts"
)

func TestResolveLoadOrderCoreFirst(t *testing.T) {
	manifests := []contracts.ModManifest{
		{
			ID:      "linearAlgebra",
			Version: "0.1.0",
			Dependencies: []contracts.ModDependency{{
				ID:      "core",
				Version: ">=0.1.0",
			}},
		},
		{ID: "core", Version: "0.1.0"},
	}

	order, err := ResolveLoadOrder(manifests)
	if err != nil {
		t.Fatalf("resolve failed: %v", err)
	}
	if len(order) != 2 {
		t.Fatalf("expected 2 items, got %d", len(order))
	}
	if order[0] != "core" {
		t.Fatalf("expected core first, got %s", order[0])
	}
}

func TestResolveLoadOrderMissingDependency(t *testing.T) {
	manifests := []contracts.ModManifest{
		{ID: "core", Version: "0.1.0"},
		{
			ID:      "linearAlgebra",
			Version: "0.1.0",
			Dependencies: []contracts.ModDependency{{
				ID:      "missing",
				Version: ">=0.1.0",
			}},
		},
	}

	if _, err := ResolveLoadOrder(manifests); err == nil {
		t.Fatal("expected missing dependency error")
	}
}

func TestResolveLoadOrderCycle(t *testing.T) {
	manifests := []contracts.ModManifest{
		{ID: "core", Version: "0.1.0"},
		{
			ID:      "a",
			Version: "0.1.0",
			Dependencies: []contracts.ModDependency{{
				ID:      "b",
				Version: ">=0.1.0",
			}},
		},
		{
			ID:      "b",
			Version: "0.1.0",
			Dependencies: []contracts.ModDependency{{
				ID:      "a",
				Version: ">=0.1.0",
			}},
		},
	}

	if _, err := ResolveLoadOrder(manifests); err == nil {
		t.Fatal("expected cycle error")
	}
}

func TestEnforceCorePresence(t *testing.T) {
	manifests := []contracts.ModManifest{{ID: "linearAlgebra", Version: "0.1.0"}}
	if err := EnforceCorePresence(manifests); err == nil {
		t.Fatal("expected core required error")
	}
}
