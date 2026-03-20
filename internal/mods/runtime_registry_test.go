package mods

import (
	"context"
	"testing"

	"gaussgo/internal/contracts"
	"gaussgo/internal/types"
)

func TestRuntimeRegistryRegisterGetListDelete(t *testing.T) {
	r := NewRuntimeRegistry()
	mod := &noopRuntimeMod{manifest: contracts.ModManifest{ID: "core"}}

	if err := r.Register(mod); err != nil {
		t.Fatalf("register failed: %v", err)
	}

	got, err := r.Get("core")
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if got.ID() != "core" {
		t.Fatalf("expected core, got %s", got.ID())
	}

	listed := r.List()
	if len(listed) != 1 {
		t.Fatalf("expected list size 1, got %d", len(listed))
	}

	if err := r.Delete("core"); err != nil {
		t.Fatalf("delete failed: %v", err)
	}
	if _, err := r.Get("core"); err == nil {
		t.Fatal("expected not found after delete")
	}
}

func TestRuntimeRegistryRejectsDuplicate(t *testing.T) {
	r := NewRuntimeRegistry()
	mod := &noopRuntimeMod{manifest: contracts.ModManifest{ID: "core"}}
	if err := r.Register(mod); err != nil {
		t.Fatalf("register failed: %v", err)
	}
	if err := r.Register(mod); err == nil {
		t.Fatal("expected duplicate registration failure")
	}
}

func TestRuntimeLifecycleSyncEnabledLoadAndUnload(t *testing.T) {
	lifecycle := NewRuntimeLifecycle(ManifestFactory{}, nil, nil, nil)
	runtimeCtx := types.RuntimeContext{}

	manifests := []contracts.ModManifest{
		{ID: "core"},
		{ID: "linearAlgebra", Dependencies: []contracts.ModDependency{{ID: "core", Version: ">=0.1.0"}}},
	}

	enabled := map[string]bool{"core": true, "linearAlgebra": true}
	if err := lifecycle.SyncEnabled(context.Background(), manifests, runtimeCtx, enabled, nil); err != nil {
		t.Fatalf("sync enabled load failed: %v", err)
	}
	if len(lifecycle.Registry().List()) != 2 {
		t.Fatalf("expected 2 runtime mods, got %d", len(lifecycle.Registry().List()))
	}

	enabled = map[string]bool{"core": true}
	if err := lifecycle.SyncEnabled(context.Background(), manifests, runtimeCtx, enabled, nil); err != nil {
		t.Fatalf("sync enabled unload failed: %v", err)
	}
	if len(lifecycle.Registry().List()) != 1 {
		t.Fatalf("expected 1 runtime mod, got %d", len(lifecycle.Registry().List()))
	}
}

func TestManifestFactoryCoreProvidesBasicMathCapability(t *testing.T) {
	factory := ManifestFactory{}
	mod, err := factory.Build(contracts.ModManifest{ID: "core"})
	if err != nil {
		t.Fatalf("build failed: %v", err)
	}

	provider, ok := mod.(contracts.CapabilityProvider)
	if !ok {
		t.Fatal("expected core mod to expose capability provider")
	}
	catalog, ok := mod.(contracts.CapabilityCatalog)
	if !ok {
		t.Fatal("expected core mod to expose capability catalog")
	}
	capabilities := catalog.Capabilities()
	if len(capabilities) == 0 || capabilities[0] != contracts.CapabilityHelpersBasicMath {
		t.Fatalf("unexpected capability catalog: %v", capabilities)
	}
	capability, exists := provider.Capability(contracts.CapabilityHelpersBasicMath)
	if !exists {
		t.Fatal("expected basic math capability to exist")
	}

	mathHelpers, ok := capability.(contracts.BasicMathHelpers)
	if !ok {
		t.Fatal("expected capability to be BasicMathHelpers")
	}

	value, err := mathHelpers.Div(context.Background(), 8, 2)
	if err != nil {
		t.Fatalf("expected div success, got %v", err)
	}
	if value != 4 {
		t.Fatalf("expected div result 4, got %v", value)
	}
}
