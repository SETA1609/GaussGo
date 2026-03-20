# Phase 5 - Learning MVP

## Goal
Ship the first complete learning loop using `linearAlgebra` mod content and state-backed progress.

## Ownership
- Domain owner

## Phase 0 Contract Inputs
- `docs/contracts/state-schema-v1.md`
- `docs/contracts/runtime-context-v1.md`
- `docs/contracts/controller-interfaces-v1.md`
- `docs/contracts/mod-manifest-v1.md`

## Scope
- Load unit content from selected mod entry path.
- Render lesson text and concept sequence.
- Evaluate initial exercise types.
- Save progress into active state.
- Expose mod-provided helper functions for operation assistance.
- Generate random quizzes from already-read concepts.
- Provide quiz result review with step-by-step solution explanations.
- Support two quiz setup paths: random concept pool or user-selected concepts.
- Accept quiz answers via choice selection or answer-letter input.

## Deliverables
- `internal/content/unit_model.go`
- `internal/content/loader.go`
- `internal/content/registry.go`
- `internal/exercises/registry.go`
- `internal/exercises/multiple_choice.go`
- `internal/exercises/compute.go`
- `internal/exercises/random_quiz.go`
- `internal/tui/scenes/random_quiz_mode_select.go`
- `internal/helpers/registry.go`
- `internal/helpers/dispatcher.go`
- `internal/tui/scenes/lesson.go`
- `internal/tui/scenes/exercise.go`
- `internal/tui/scenes/quiz_results.go`

## Exercise MVP
- `multiple_choice`
- `compute`
- `random_quiz` (question set sampled from already-read concepts, via random mode or selected-concepts mode)

## Progress Rules
- Persist completion at unit and concept granularity.
- Mark concept as read once lesson view is opened/completed.
- Store exercise scores per concept.
- Update `updatedAt` on every persistence write.

## Learning Hub Feature Rules
- Statistics:
  - Group statistics by mod with expandable sections.
  - Compute per-mod correctness percent from `correct/attempted`.
  - Show total exercises attempted/completed.
  - Expanded section shows concept-level table with attempted, correct, and correctness percent.
- Lessons:
  - Concepts come from selected mod content (for example vectors/matrices in `linearAlgebra`).
- Random Quiz:
  - Source pool is only concepts in `progress.<modId>.readConcepts`.
  - Setup mode `Random` uses all eligible read concepts.
  - Setup mode `Selected Concepts` allows concept multi-select from eligible read concepts.
  - Questions are generated randomly from the final mode-resolved concept pool.
  - Answer input supports option selection and answer-letter typing.
  - Completion view shows aggregate results and per-question step-by-step solution.
  - If answer is wrong, include clarification of why it is wrong and how to solve correctly.
- Helpers:
  - Helper functions are discovered from enabled mods.
  - Helper invocation returns structured explanation/result payload.

## Tests
- Content loading from mod path with namespace IDs.
- Exercise evaluation correctness.
- Progress save/load continuity across restarts.
- Random quiz excludes unread concepts.
- Statistics calculations match persisted attempts/correct values.
- Helper registry discovers and executes mod helper functions.
- Quiz results include step-by-step solution rendering and wrong-answer clarification checks.
- Quiz engine tests verify both setup modes produce randomized problems.
- Quiz input tests verify both click/select and letter-entry paths.

## Exit Criteria
- User can choose `linearAlgebra`, read lessons, run random quiz from read concepts, review step-by-step quiz solutions with clarifications, call helper functions, and see persisted stats/progress on reload.
