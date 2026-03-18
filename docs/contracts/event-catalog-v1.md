# Event Catalog Contract v1

## Version
- `event-catalog/v1`

## Purpose
- Define canonical event names, minimal payload shapes, producer ownership, and expected consumers.
- Keep event naming stable across controllers, services, and TUI flows.

## Envelope
- Event envelope is defined by `docs/contracts/eventbus-service-v1.md`:
  - `name`
  - `timestamp`
  - `source`
  - `data`

## Naming Convention
- Use lowercase dot-separated names: `<domain>.<action>`.
- Use past-tense for completed state changes and present-progress form for started flows.

## Catalog

### `state.loaded`
- Producer: state load flow / app bootstrap
- Consumers: diagnostics UI, logger, analytics hooks
- Required `data`:
  - `stateId` (string)
  - `locale` (string)
  - `normalized` (bool)

### `state.saved`
- Producer: state store save operation
- Consumers: diagnostics UI, logger
- Required `data`:
  - `stateId` (string)
  - `updatedAt` (string, RFC3339)

### `locale.changed`
- Producer: locale controller
- Consumers: TUI re-render trigger, logger
- Required `data`:
  - `stateId` (string)
  - `previous` (string)
  - `current` (string)

### `mods.refreshed`
- Producer: mod controller refresh flow
- Consumers: mod settings scene, logger
- Required `data`:
  - `total` (number)
  - `valid` (number)
  - `invalid` (number)

### `mods.toggled`
- Producer: mod controller toggle flow
- Consumers: mod settings scene, logger
- Required `data`:
  - `modId` (string)
  - `enabled` (bool)

### `quiz.started`
- Producer: learning controller / quiz setup flow
- Consumers: quiz scene, logger
- Required `data`:
  - `modId` (string)
  - `mode` (string: `random` or `selected_concepts`)
  - `conceptCount` (number)

### `quiz.finished`
- Producer: quiz evaluator/results flow
- Consumers: results scene, progress writer, logger
- Required `data`:
  - `modId` (string)
  - `attempted` (number)
  - `correct` (number)
  - `correctnessPercent` (number)

## Validation Rules
- Event `name` must be in this catalog or in an approved extension namespace.
- Required payload keys must be present for each event.
- Unknown extra payload keys are allowed.

## Backward Compatibility Notes
- Adding new events is non-breaking.
- Removing/renaming catalog events is breaking and requires a major catalog revision.
- Removing required payload keys is breaking.
