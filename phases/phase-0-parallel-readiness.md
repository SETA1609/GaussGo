# Phase 0 - Parallel Readiness

## Goal
Define all shared dependencies, contracts, and repo conventions so Phases 1-5 can proceed in parallel with minimal blocking.

## Ownership
- Platform/framework owner

## Why This Phase Exists
Without a contract-first foundation, parallel work causes drift across models, interfaces, naming, and runtime assumptions. Phase 0 removes those blockers first.

## Scope
- Shared domain schemas and interface contracts.
- Shared package/module boundaries.
- Build/test/tooling standards.
- Integration points and ownership matrix.
- Non-functional conventions (errors, logging, naming, docs).

## Deliverables

### 1) Contract Pack (Authoritative)
- `docs/contracts/mod-manifest-v1.md`
- `docs/contracts/state-schema-v1.md`
- `docs/contracts/runtime-context-v1.md`
- `docs/contracts/controller-interfaces-v1.md`
- `docs/contracts/i18n-resolution-v1.md`
- `docs/contracts/logging-service-v1.md`
- `docs/contracts/eventbus-service-v1.md`
- `docs/contracts/event-catalog-v1.md`

Each contract must include:
- Version identifier
- Required fields and invariants
- Validation rules
- Example payloads
- Backward compatibility notes

### 2) Interface Stubs (Compile Targets)
- `internal/contracts/mods.go`
- `internal/contracts/state.go`
- `internal/contracts/controllers.go`
- `internal/contracts/i18n.go`
- `internal/contracts/logging.go`
- `internal/contracts/events.go`

These stubs are minimal compile-time interfaces used by all phases.

### 3) Shared Types and Errors
- `internal/types/runtime_context.go` (must include `CurrentLocale`)
- `internal/types/ids.go` (namespaced ID helpers)
- `internal/apperrors/codes.go` (shared error codes)
- `internal/apperrors/error_types.go`

### 4) Tooling Baseline
- `Makefile` targets:
  - `make test`
  - `make lint`
  - `make check`
  - `make run`
- CI workflow skeleton for test/lint/check.
- Common Go version pin and formatting/lint rules.

### 5) Repo Structure Baseline
- Add `docs/` tree for contracts and ADR-like notes.
- Confirm `phases/` is the execution plan source.
- Add ownership note per phase in each phase file.

### 6) Integration Map
Create `docs/integration-map.md` including:
- Producer/consumer map per phase
- Expected handoff artifacts
- Merge order and risk points
- Contract change procedure

## Mandatory Invariants (Unlocked for All Phases)
- `core` mod is mandatory and never disableable.
- Manifest requires localization object (`default`, `supported`, `path`).
- Default locale is `en`.
- Locale selection persists to `state.ui.locale`.
- Mod refresh rescans and updates statuses, but does not auto-enable new mods.
- Logging is structured and supports at least `debug`, `info`, `warn`, `error` levels.
- Event bus supports emit and subscribe/unsubscribe semantics.

## Parallelization Strategy

### Workstreams after Phase 0
- Workstream A: Phase 1 (registry + mods discovery)
- Workstream B: Phase 2 (state core)
- Workstream C: Phase 3 (controllers + bootstrap)
- Workstream D: Phase 4 (tui scenes)
- Workstream E: Phase 5 (learning mvp)

### Dependency Edges
- Phase 4 depends on controller interfaces and runtime context contract from Phase 0.
- Phase 5 depends on content/exercise contracts from Phase 0.
- Phase 6 depends on outputs from at least one implemented vertical flow.

## Ownership Matrix (Suggested)
- Contracts + shared types: platform/framework owner
- Mods + registry: systems owner
- State + persistence: data/runtime owner
- Controllers + bootstrap: app orchestration owner
- TUI scenes: UX/TUI owner
- Learning evaluators/content: domain owner

## Done Criteria
- All contract docs merged and versioned.
- Shared interfaces compile without implementation details.
- Tooling/CI baseline is green.
- All phase docs reference Phase 0 contracts as inputs.
- Team can start Phases 1-5 in parallel without ambiguous assumptions.

## Exit Criteria
Phase 0 is complete when no phase needs to invent its own schema/interface to continue implementation.
