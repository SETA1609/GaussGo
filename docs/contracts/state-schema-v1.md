# State Schema Contract v1

## Version
- `state/v1`

## Required Fields and Invariants
- Required: `schemaVersion`, `stateId`, `profileName`, `createdAt`, `updatedAt`, `modsAtCreation`, `enabledMods`, `ui`, `progress`.
- `schemaVersion` is integer `1`.
- `stateId` is unique within `/states`.
- `enabledMods` must always contain `core`.
- `ui.locale` defaults to `en` when missing.
- `ui.lastScene` is optional.
- `ui.preferences` is optional free-form map.
- `progress.<modId>.readConcepts` stores concept IDs read at least once.
- `progress.<modId>.exerciseStats` stores `attempted`, `correct`, `correctnessPercent`.
- `correctnessPercent` is derived from `correct/attempted` when not explicitly stored.

## Validation Rules
- Reject missing required root fields.
- Reject non-ISO-8601 timestamps in `createdAt` and `updatedAt`.
- Normalize missing `core` in `enabledMods` by auto-inserting it.
- Normalize missing `ui` or `ui.locale` by setting locale to `en`.
- Clamp negative `attempted`/`correct` values to validation error.
- If `attempted == 0`, treat correctness percent as `0`.

## Example Payload
```json
{
  "schemaVersion": 1,
  "stateId": "default",
  "profileName": "Default Profile",
  "createdAt": "2026-03-17T00:00:00Z",
  "updatedAt": "2026-03-17T00:00:00Z",
  "modsAtCreation": [
    { "id": "core", "version": "0.1.0" },
    { "id": "linearAlgebra", "version": "0.1.0" }
  ],
  "enabledMods": ["core", "linearAlgebra"],
  "ui": {
    "locale": "en",
    "lastScene": "main_menu",
    "preferences": {}
  },
  "progress": {
    "linearAlgebra": {
      "readConcepts": ["vectors.dot-product"],
      "exerciseStats": {
        "attempted": 12,
        "correct": 9,
        "correctnessPercent": 75
      },
      "exerciseHistory": []
    }
  }
}
```

## Backward Compatibility Notes
- Additive fields are allowed and should be preserved when possible.
- Contract-breaking structural changes require a migration step and new schema version.
- Missing `ui.locale` from older states must be normalized to `en`.
