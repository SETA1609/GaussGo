# Phase 3 - Controllers and Runtime Orchestration

## Goal
Introduce runtime controllers and context management separate from TUI rendering.

## Ownership
- App orchestration owner

## Phase 0 Contract Inputs
- `docs/contracts/controller-interfaces-v1.md`
- `docs/contracts/runtime-context-v1.md`
- `docs/contracts/state-schema-v1.md`
- `docs/contracts/i18n-resolution-v1.md`
- `docs/contracts/logging-service-v1.md`
- `docs/contracts/eventbus-service-v1.md`

## Scope
- App lifecycle orchestration through controllers.
- Scene navigation and history.
- Learning session state selection (mod/unit/concept).
- Mod enable/disable logic.
- Locale selection controller.
- Controller-level logging and event emission.
- Runtime mod instance lifecycle orchestration (build/init/register).

## Deliverables
- `internal/controllers/app_controller.go`
- `internal/controllers/scene_controller.go`
- `internal/controllers/learning_controller.go`
- `internal/controllers/mod_controller.go`
- `internal/controllers/locale_controller.go`
- `internal/controllers/events.go`
- `internal/mods/runtime_registry.go`
- `internal/mods/runtime_factory.go` (or equivalent constructor wiring)
- Runtime context model with locale.

## Runtime Context
Must include:
- `ActiveStateID`
- `EnabledMods`
- `CurrentSceneID`
- `CurrentModID`
- `CurrentUnitID`
- `CurrentConceptID`
- `CurrentLocale`

## Controller Rules
- `ModController` rejects disabling `core`.
- `LocaleController` supports `en` and `es` initially, extensible later.
- `LocaleController` persists selected locale into active state.
- `SceneController` owns navigation stack and back behavior.
- Controllers emit domain events for state transitions and user actions.
- Controllers use `LoggingService` for structured operational logs.

## Bootstrap Responsibilities
- Boot order: mods -> state -> registries -> controllers -> tui.
- Bootstrap builds immutable snapshots for read-mostly data.
- Bootstrap returns structured startup diagnostics for UI display.

## Tests
- Controller behavior tests for scene transitions.
- Core toggle rejection tests.
- Locale change updates context and state.
- Bootstrap ordering and failure tests.
- Runtime mod initialization order and failure-path tests.
- Runtime capability lookup tests from controller call sites.

## Exit Criteria
- Non-interactive runtime can execute: load state, set locale, select mod, and emit context.

## Implementation Status
- Implemented in:
  - `internal/controllers/app_controller.go`
  - `internal/controllers/scene_controller.go`
  - `internal/controllers/learning_controller.go`
  - `internal/controllers/mod_controller.go`
  - `internal/controllers/locale_controller.go`
  - `internal/controllers/events.go`
  - `internal/bootstrap/bootstrap.go` (runtime bootstrap wiring)
  - `cmd/gaussgo/main.go` (startup orchestration entrypoint)
- Tests added:
  - `internal/controllers/app_controller_test.go`
  - `internal/controllers/scene_controller_test.go`
  - `internal/controllers/learning_controller_test.go`
  - `internal/controllers/mod_controller_test.go`
  - `internal/controllers/locale_controller_test.go`
  - `internal/bootstrap/bootstrap_test.go` (runtime bootstrap scenarios)
