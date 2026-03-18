# Phase 6 - Hardening

## Goal
Stabilize runtime behavior, diagnostics, compatibility handling, and container execution quality.

## Ownership
- Platform/framework owner

## Phase 0 Contract Inputs
- `docs/contracts/mod-manifest-v1.md`
- `docs/contracts/state-schema-v1.md`
- `docs/contracts/runtime-context-v1.md`
- `docs/contracts/i18n-resolution-v1.md`

## Scope
- Startup diagnostics and error UX.
- Compatibility checks between state and current mod versions.
- Migration placeholders for state schema evolution.
- Container verification and runtime docs alignment.

## Deliverables
- Compatibility validator for `modsAtCreation` vs discovered mods.
- Human-readable warning model for UI and logs.
- State migration scaffold (`schemaVersion` dispatcher).
- Updated runbook snippets in docs.

## Compatibility Rules
- Warn when mod missing from current environment but present in state.
- Warn on version drift; do not block if compatible policy says safe.
- Block startup only on core contract failures.

## Diagnostics
- Structured startup report:
  - discovered mods
  - invalid mods
  - dependency issues
  - locale fallback notices
  - active state normalization events

## Tests
- Compatibility warning generation.
- Migration fallback behavior.
- Container smoke run.

## Exit Criteria
- Framework provides actionable errors, compatibility notices, and reliable startup in local and container contexts.
