# Phase 1 - Foundation

## Goal
Establish framework primitives: base registry, mod discovery, manifest validation, dependency resolution, and bootstrap wiring contracts.

## Ownership
- Systems owner

## Phase 0 Contract Inputs
- `docs/contracts/mod-manifest-v1.md`
- `docs/contracts/runtime-context-v1.md`
- `docs/contracts/controller-interfaces-v1.md`
- `docs/contracts/i18n-resolution-v1.md`
- `docs/contracts/logging-service-v1.md`
- `docs/contracts/eventbus-service-v1.md`

## Scope
- Generic base registry implementation.
- Mod manifest structs and parser.
- Manifest validation (including required localization object).
- Dependency graph resolution and deterministic load order.
- `internal/bootstrap` responsibilities defined for startup orchestration.

## Deliverables
- `internal/registry/base.go`
- `internal/registry/errors.go`
- `internal/mods/manifest.go`
- `internal/mods/discover.go`
- `internal/mods/validate.go`
- `internal/mods/resolve.go`
- `internal/mods/status.go`
- `internal/mods/registry.go`
- `internal/bootstrap/bootstrap.go`
- `internal/bootstrap/services.go` (service wiring for logger and event bus)

## Required Rules
- Fail fast if `/mods/core/manifest.json` is missing or invalid.
- `core` must be first in valid topological order.
- Manifest requires `localization` with:
  - `default`
  - `supported`
  - `path`
- Validate `localization.default` is contained in `localization.supported`.

## Refresh Mods Contract
- Refresh operation rescans `/mods` and re-runs validation/resolution.
- Refresh updates status list in-memory.
- Refresh does not auto-enable newly discovered mods; user must enable manually in Mod Settings.

## Tests
- Registry CRUD and duplicate registration behavior.
- Manifest validation error cases (missing fields, invalid semver, invalid localization).
- Dependency resolver cycle detection and missing dependency errors.
- Core presence enforcement.

## Exit Criteria
- Runtime can print discovered mods with statuses and load order.
- Refresh function returns deterministic updated mod snapshot.

## Implementation Status
- Implemented in:
  - `internal/registry/base.go`
  - `internal/registry/errors.go`
  - `internal/mods/manifest.go`
  - `internal/mods/discover.go`
  - `internal/mods/validate.go`
  - `internal/mods/resolve.go`
  - `internal/mods/status.go`
  - `internal/mods/registry.go`
  - `internal/bootstrap/bootstrap.go`
  - `internal/bootstrap/services.go`
- Tests added:
  - `internal/registry/base_test.go`
  - `internal/mods/discover_test.go`
  - `internal/mods/validate_test.go`
  - `internal/mods/resolve_test.go`
  - `internal/mods/status_test.go`
  - `internal/mods/registry_test.go`
  - `internal/bootstrap/services_test.go`
  - `internal/bootstrap/bootstrap_test.go`

### Bug-Report Follow-up
- Duplicate dependency IDs are now rejected in manifest validation.
- Repository refresh status cache is guarded with mutex for concurrency safety.
- Event type duplication removed in favor of `contracts.Event`.
- Structured errors moved to `internal/apperrors` and integrated in mods/bootstrap/event bus paths.
