# Integration Map

## Producer/Consumer by Phase
- Phase 0 produces contracts, shared interfaces, shared types, and tooling baseline consumed by Phases 1-6.
- Phase 1 produces mod discovery/resolution and bootstrap contracts consumed by Phases 3, 4, and 6.
- Phase 2 produces state model/store consumed by Phases 3, 4, 5, and 6.
- Phase 3 produces controller orchestration consumed by Phase 4.
- Phase 4 produces TUI scene flow consumed by integration testing in Phase 6.
- Phase 5 produces learning/content execution consumed by Phase 6 hardening.

## Handoff Artifacts
- Phase 0 -> all: `docs/contracts/*`, `internal/contracts/*`, shared `internal/types/*`, shared `internal/errors/*`.
- Phase 1 -> 3/4/6: mod registry API, refresh behavior, bootstrap diagnostics.
- Phase 2 -> 3/4/5/6: state persistence API, normalization behavior.
- Phase 3 -> 4: stable controller interfaces and runtime context integration.
- Phase 5 -> 6: quiz/stats/progress behaviors and validation scenarios.

## Merge Order and Risk Points
- Merge Phase 0 first; treat contracts as locked inputs.
- Merge Phase 1 and 2 next; risk is drift between discovered mods and state normalization.
- Merge Phase 3 after Phase 1/2 APIs are stable; risk is controller mismatch.
- Merge Phase 4 and 5 in parallel once interfaces are stable; risk is UI assumptions diverging from controller/state behavior.
- Merge Phase 6 last; risk is late discovery of contract inconsistencies.

## Contract Change Procedure
- Propose change in contract doc with explicit versioning impact.
- Update corresponding compile-time stubs under `internal/contracts`.
- Run `make check` and affected tests.
- Announce downstream phase impact in PR description.
- If breaking, publish a new contract version document instead of mutating old semantics silently.
