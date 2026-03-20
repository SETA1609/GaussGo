# GaussGo Plan Index

## Purpose
This file is the planning index for GaussGo.

The runtime/framework is the Go program itself. Mods integrate with the framework through stable contracts, with `core` as the required baseline mod.

## Principles
- Local-first CLI/TUI runtime and local state files.
- Data-driven JSON/YAML content.
- Drop-in mods from `/mods/<id>`.
- Deterministic startup validation before entering UI.
- `core` is always required and always enabled.
- Default locale is English (`en`) and locale is stored per state.

## Architecture Decisions
- Registry-first architecture with a reusable base registry.
- Controllers own runtime state transitions (`scene`, `learning`, `mods`, `locale`).
- TUI scenes remain thin and call controllers.
- Each mod declares localization in `manifest.json`.
- Cross-cutting runtime services include `LoggingService` and `EventBusService`.

## Current Baseline
- Contracts and compile-time interfaces are in place.
- Phases 0-6 are implemented at framework level.
- Startup and controllers are wired with state, mods, i18n, logging, and event bus.
- TUI navigation and learning-hub shells are available.

## Mod and State Rules (Canonical)
- `core` manifest must exist or startup fails.
- `core` cannot be disabled in Mod Settings.
- Main menu includes `Select Language`.
- Selected language is saved to `state.ui.locale`.
- If `state.ui.locale` is missing, load defaults to `en`.
- Mod Settings includes `Refresh Mods` action.
- Start Learning opens a mod-aware Learning Hub with: `Statistics`, `Lessons`, `Random Quiz`, `Helpers`, `Back to Main Menu`.
- Every menu/scene includes a `Back` action that returns one step up the navigation tree.
- Statistics are grouped by mod, with expandable concept-level tables per mod.
- Each concept row includes at least attempted, correct, and correctness percent metrics.
- Random Quiz only includes concepts marked as read at least once in state.
- Random Quiz has two setup modes: `Random` and `Selected Concepts`.
- `Selected Concepts` mode lets the user pick the concept subset before quiz generation.
- Random Quiz questions are generated randomly from the selected eligible concept pool.
- Multiple-choice answers can be selected from options or entered by letter shortcut.
- Quiz completion shows final results and a step-by-step solution review.
- Wrong answers include a clarification explaining the mistake and the correct reasoning.
- Quiz/problem content should support LaTeX-style math formatting when renderer support exists, with plain-text fallback.
- Helpers are callable operations/functions provided by enabled mods.
- Runtime emits domain events for major actions (state load/save, locale change, mod refresh/toggle, quiz start/finish).
- Runtime logs structured events through `LoggingService` with level and context fields.

## Open Planning Items
1. Introduce a runtime mod type (`ModRuntime`) that carries both manifest metadata and executable behavior.
2. Add a runtime mod registry for initialized mod instances, separate from manifest/status discovery.
3. Define lifecycle hooks for mod instances (`Init`, optional `Shutdown`) and enforce deterministic initialization by resolved load order.
4. Add capability exposure rules so controllers/scenes can access mod features safely (for example helpers) via typed interfaces instead of loose maps.
5. Wire bootstrap to build, initialize, and register runtime mod instances after manifest resolution.
6. Add docs/contracts for runtime mod capabilities and bootstrap lifecycle (`docs/contracts/mod-runtime-v1.md`) where shared runtime types are owned by `core`, then align phase docs and README.
7. Add tests for runtime mod registration, duplicate capability conflicts, initialization failures, and capability lookup behavior.
8. Introduce thin TUI adapter ports around Bubble Tea/Charm boundaries so scene logic depends on internal interfaces, not vendor types.

## Phase Documents
- Phase 0: `phases/phase-0-parallel-readiness.md`
- Phase 1: `phases/phase-1-foundation.md`
- Phase 2: `phases/phase-2-state-core.md`
- Phase 3: `phases/phase-3-controllers-runtime.md`
- Phase 4: `phases/phase-4-tui-scenes.md`
- Phase 5: `phases/phase-5-learning-mvp.md`
- Phase 6: `phases/phase-6-hardening.md`

## Suggested Next Execution Order
1. Define `ModRuntime` and capability contracts in `internal/contracts/mods.go` and contract docs.
2. Implement runtime mod registry under `internal/mods` and bootstrap lifecycle wiring.
3. Add TUI port interfaces (program runner, input mapping, renderer primitives) and implement Bubble Tea/Charm adapters in `internal/tui/adapters`.
4. Migrate helper access paths to runtime capabilities and keep manifest `provides` as declarative metadata.
5. Expand tests and update integration docs (`docs/integration-map.md`) and README.

## Acceptance Summary
- Runtime supports initialized mod instances registered by mod ID.
- Runtime exposes mod capabilities/functions through a typed access layer.
- Bootstrap guarantees deterministic mod initialization order and clear startup diagnostics on failures.
- Helper/feature calls resolve through runtime capabilities of enabled mods.
- Contract docs and tests cover runtime mod lifecycle and capability lookup.

## Immediate Next Task
Create a small design slice for runtime mod instances:
- Add `ModRuntime` + capability contracts.
- Implement runtime mod registry.
- Wire bootstrap initialization and add focused tests.
