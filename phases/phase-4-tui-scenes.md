# Phase 4 - TUI Scenes

## Goal
Build Bubble Tea scenes on top of controllers, with thin view logic and strong UX flow.

## Ownership
- UX/TUI owner

## Phase 0 Contract Inputs
- `docs/contracts/controller-interfaces-v1.md`
- `docs/contracts/runtime-context-v1.md`
- `docs/contracts/i18n-resolution-v1.md`

## Scope
- Main menu scene.
- State creation and loading scenes.
- Mod settings scene with refresh action.
- Language selection scene.
- Mod and unit selection scenes.
- Learning Hub scene after Start Learning.

## Deliverables
- `internal/tui/app_model.go`
- `internal/tui/scenes/main_menu.go`
- `internal/tui/scenes/state_create.go`
- `internal/tui/scenes/state_load.go`
- `internal/tui/scenes/mod_settings.go`
- `internal/tui/scenes/language_select.go`
- `internal/tui/scenes/mod_select.go`
- `internal/tui/scenes/unit_select.go`
- `internal/tui/scenes/learning_hub.go`
- `internal/tui/scenes/statistics.go`
- `internal/tui/scenes/random_quiz_select.go`
- `internal/tui/scenes/helpers.go`

## Required UX Behaviors
- Main menu options:
  - Start Learning
  - Create New State
  - Load Existing State
  - Select Language
  - Mod Settings
  - Exit
- Every scene provides a `Back` action that returns exactly one level up in navigation.
- Mod Settings has `Refresh Mods` command.
- Disabled reason is visible for unavailable mods.
- Core row is visibly locked and non-toggleable.
- After choosing Start Learning and selecting a mod, show Learning Hub options:
  - Statistics
  - Lessons
  - Random Quiz
  - Helpers
  - Back to Main Menu
- Statistics scene groups data by mod and supports expanding each mod row.
- Expanded mod view shows a concept-level table (for example matrices, vectors) with per-concept metrics.
- Random Quiz concept picker allows multi-select only from concepts already read once.
- Random Quiz setup offers two options: `Random` and `Selected Concepts`.
- `Selected Concepts` opens concept multi-select picker; `Random` skips concept picking.
- Both setup modes still generate randomized quiz questions from the eligible pool.
- Exercise/quiz scene accepts answer by choice selection or by typing option letter.
- Quiz results scene presents final score and per-question step-by-step review.
- Wrong-answer review includes a clarification of the error and expected reasoning.
- Math/problem text rendering prefers LaTeX-style expressions when supported by the TUI rendering path.

## Localization Integration
- Scene labels must resolve through i18n keys.
- If key missing in active mod, fallback to core locale.
- If key missing in core, show key literal.

## Tests
- Scene-level update tests for menu actions.
- Refresh action updates mod table.
- Language selection updates labels after re-render.
- Statistics scene test coverage for mod expand/collapse and concept table rendering.
- Navigation tests verify `Back` behavior is one-step and consistent across scenes.
- Random quiz setup tests cover both `Random` and `Selected Concepts` branches.
- Quiz input tests cover both answer-letter input and direct choice selection.
- Random quiz result scene tests for step-by-step and wrong-answer clarification blocks.

## Exit Criteria
- User can navigate full menu tree, refresh mods, and change language from UI.

## Implementation Status
- Implemented in:
  - `internal/tui/app_model.go`
  - `internal/tui/runner.go`
  - `internal/tui/scenes/main_menu.go`
  - `internal/tui/scenes/state_create.go`
  - `internal/tui/scenes/state_load.go`
  - `internal/tui/scenes/mod_settings.go`
  - `internal/tui/scenes/language_select.go`
  - `internal/tui/scenes/mod_select.go`
  - `internal/tui/scenes/unit_select.go`
  - `internal/tui/scenes/learning_hub.go`
  - `internal/tui/scenes/statistics.go`
  - `internal/tui/scenes/random_quiz_select.go`
  - `internal/tui/scenes/helpers_scene.go`
  - `internal/tui/scenes/helpers.go`
  - `cmd/gaussgo/main.go` (startup now launches interactive scene runtime)
- Tests added:
  - `internal/tui/app_model_test.go` (navigation back behavior, refresh mods action, locale relabel check)
- Reusable UI primitives extracted for cross-interface reuse:
  - `internal/tui/components/theme.go`
  - `internal/tui/components/menu.go`
  - `internal/tui/components/layout.go`
  - `internal/tui/components/status.go`
- Data-driven scene schemas are now file-backed and loaded at runtime:
  - `mods/core/ui/scenes/main_menu.json`
  - `mods/core/ui/scenes/language_select.json`
  - `internal/tui/scenes/menu_schema.go` (filesystem loader + schema binding)
