# Phase 2 - State Core

## Goal
Implement persistent state management with strict core enforcement and locale persistence.

## Ownership
- Data/runtime owner

## Phase 0 Contract Inputs
- `docs/contracts/state-schema-v1.md`
- `docs/contracts/mod-manifest-v1.md`
- `docs/contracts/i18n-resolution-v1.md`
- `docs/contracts/runtime-context-v1.md`

## Scope
- State model and schema validation.
- Create/load/save state operations.
- Active state pointer handling.
- Automatic normalization rules on load.

## Deliverables
- `internal/state/model.go`
- `internal/state/store.go`
- `internal/state/validate.go`
- `internal/state/active.go`

## State Schema Requirements
- `schemaVersion`
- `stateId`
- `profileName`
- `createdAt`, `updatedAt`
- `modsAtCreation`
- `enabledMods`
- `progress`
- `ui` object:
  - `locale` (default `en`)
  - `lastScene` (optional)
  - `preferences` (optional map)

Progress must support learning hub features:
- `progress.<modId>.readConcepts`: concept ids read at least once
- `progress.<modId>.exerciseStats`:
  - `attempted`
  - `correct`
  - `correctnessPercent` (derived or cached)
- `progress.<modId>.exerciseHistory` (optional lightweight history)

## Core and Locale Rules
- If `core` missing in `enabledMods`, add it automatically.
- `ui.locale` missing -> set `en`.
- If locale not supported globally, fallback to `en` and record warning.
- Random quiz eligibility is based on `readConcepts` only.

## Atomic Persistence
- Save through temp file + fsync + rename.
- Preserve file permissions suitable for local profile storage.

## Tests
- Create state initializes `core` and `ui.locale=en`.
- Load state normalizes missing `core` and missing locale.
- Active pointer read/write behavior.
- State roundtrip save/load integrity.
- Read concept tracking persists and reloads correctly.
- Mod exercise statistics persist and compute correctly.

## Exit Criteria
- Multiple states can be created and loaded.
- Current active state is tracked in `/states/active.txt`.

## Implementation Status
- Implemented in:
  - `internal/state/model.go`
  - `internal/state/validate.go`
  - `internal/state/store.go`
  - `internal/state/active.go`
- Tests added:
  - `internal/state/store_test.go`
  - `internal/state/active_test.go`
  - `internal/state/test_helpers_test.go`
