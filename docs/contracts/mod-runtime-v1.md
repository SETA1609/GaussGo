# Mod Runtime Contract v1

## Version
- `mod-runtime/v1`

## Purpose
Define how executable mod behavior is exposed at runtime, separate from static manifest metadata.

## Core Ownership Rule
- The runtime mod contract types live in `core` runtime packages (`internal/contracts` in current repo).
- `core` is the provider of shared runtime APIs and lifecycle semantics.
- Every non-core mod must declare dependency on `core` in `manifest.json`.
- `core` itself cannot depend on other mods.

## Required Runtime Concepts
- `ModManifest` remains static metadata (`id`, `version`, `dependencies`, `provides`, etc.).
- `ModRuntime` is an initialized mod instance with executable behavior.
- `RuntimeModRegistry` stores initialized `ModRuntime` instances by `modID`.
- Capability access is typed-first (preferred), with optional generic lookup for extensibility.

## Required Interfaces and Invariants
```go
type ModRuntime interface {
    ID() string
    Manifest() ModManifest
    Init(ctx RuntimeContext) error
}

type RuntimeModRegistry interface {
    Register(mod ModRuntime) error
    Get(modID string) (ModRuntime, error)
    List() map[string]ModRuntime
}

type CapabilityProvider interface {
    Has(capability string) bool
    Capability(capability string) (any, bool)
}
```

Invariants:
- `ModRuntime.ID()` must equal `ModRuntime.Manifest().ID`.
- Runtime registration keys are unique by `modID`.
- Initialization order follows resolved load order.
- A mod cannot be initialized unless all declared dependencies are already initialized.
- On init failure, startup is aborted with diagnostics and no partial "ready" state is exposed.

## Bootstrap Lifecycle Rules
- Discover and validate manifests first.
- Resolve deterministic load order.
- Build runtime instances (factory/constructor stage).
- Initialize in load order.
- Register successfully initialized instances in `RuntimeModRegistry`.
- Expose registry as read-mostly dependency for controllers/services.

## Capability Exposure Rules
- Prefer explicit typed capability interfaces for critical flows (for example helper execution).
- Generic `Capability(string)` access is allowed for optional/experimental features.
- Capability names should be namespaced (`helpers.compute`, `content.units`, etc.).
- If two enabled mods claim the same exclusive capability key, bootstrap must fail with a conflict error unless multi-provider behavior is explicitly defined.

## Validation Rules
- Reject runtime registration when `modID` is empty.
- Reject duplicate runtime registration for same `modID`.
- Reject capability resolution requests for unknown capabilities with typed contract errors.
- Reject enabling mods that do not depend on `core` (except `core` itself).

## Suggested Diagnostic Events
- `mod.runtime.init.started`
- `mod.runtime.init.succeeded`
- `mod.runtime.init.failed`
- `mod.runtime.capability.conflict`

## Example Flow
1. `core` + `linearAlgebra` manifests are discovered.
2. Resolver returns load order: `core`, `linearAlgebra`.
3. Bootstrap builds `core` runtime instance, then `linearAlgebra` runtime instance.
4. `core` initializes first; contract APIs are available.
5. `linearAlgebra` initializes and registers helper/content capabilities.
6. Learning/Helper flows query runtime registry to call mod behavior.

## Backward Compatibility Notes
- Adding optional capability keys is non-breaking.
- Adding methods to `ModRuntime` is breaking; use optional extension interfaces instead.
- Changing lifecycle semantics (for example allowing partial init success) is breaking.
