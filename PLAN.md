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

## Gap Tracker (Included in Plan)
The following gaps are explicitly included and must be implemented:
1. Add `LocaleController` for runtime locale switching and persistence.
2. Add `internal/i18n` package to the implementation map.
3. Extend runtime context with `CurrentLocale`.
4. Clarify `Refresh Mods` behavior (rescan + status refresh + no auto-enable unless user confirms).
5. Add and document `internal/bootstrap` package responsibilities.
6. Provide explicit `ui` state object example (`locale`, `lastScene`, preferences).
7. Keep heading levels and structure consistent across phase docs.
8. Add contracts and wiring plan for `LoggingService`.
9. Add contracts and wiring plan for `EventBusService` (emit + subscribe).

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

## Phase Documents
- Phase 0: `phases/phase-0-parallel-readiness.md`
- Phase 1: `phases/phase-1-foundation.md`
- Phase 2: `phases/phase-2-state-core.md`
- Phase 3: `phases/phase-3-controllers-runtime.md`
- Phase 4: `phases/phase-4-tui-scenes.md`
- Phase 5: `phases/phase-5-learning-mvp.md`
- Phase 6: `phases/phase-6-hardening.md`

## Suggested Execution Order
1. Complete Phase 0 first to define all shared contracts and tooling.
2. After Phase 0, run Phases 1-5 in parallel where possible.
3. Start Phase 6 after at least one vertical slice from Phases 1-5 is working.
4. Final integration pass resolves cross-phase merge points.

## Acceptance Summary
- New valid mod folder appears after refresh without code changes.
- `core` always present and always enabled.
- User can create/load states and switch language.
- Language persists per state.
- User can view per-mod statistics (correctness % + exercises done).
- User can expand a mod in Statistics and view concept-level stats in table form.
- User can read lessons from mod concepts (for example vectors/matrices from `linearAlgebra`).
- User can run random quizzes from already-read concepts only.
- User can run quiz in `Random` mode or `Selected Concepts` mode (both generate random problems from their pool).
- User can answer by selecting a choice or typing the answer letter.
- User can review step-by-step quiz solutions with clarifications for incorrect answers.
- User can call helper functions provided by mods.
- User can learn from at least one unit and save progress/state.
- Runtime supports centralized structured logging and event emission/subscription for orchestration.

## Immediate Next Task
Start with `phases/phase-0-parallel-readiness.md` to unlock parallel execution of the remaining phases.
