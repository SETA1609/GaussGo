# Runtime Context Contract v1

## Version
- `runtime-context/v1`

## Required Fields and Invariants
- Required fields:
  - `ActiveStateID`
  - `EnabledMods`
  - `CurrentSceneID`
  - `CurrentModID`
  - `CurrentUnitID`
  - `CurrentConceptID`
  - `CurrentLocale`
- `CurrentLocale` defaults to `en`.
- `EnabledMods` always includes `core`.
- IDs use namespaced form where applicable (`mod.unit.concept` style or equivalent helper output).

## Validation Rules
- Context must not be initialized without `ActiveStateID` after state load.
- `CurrentLocale` must be one of discovered supported locales or fallback to `en`.
- `CurrentModID` must exist in enabled mod set when not empty.

## Example Payload
```json
{
  "ActiveStateID": "default",
  "EnabledMods": ["core", "linearAlgebra"],
  "CurrentSceneID": "learning_hub",
  "CurrentModID": "linearAlgebra",
  "CurrentUnitID": "vectors",
  "CurrentConceptID": "vectors.dot-product",
  "CurrentLocale": "en"
}
```

## Backward Compatibility Notes
- New optional fields can be appended without breaking v1 consumers.
- Renaming any required field requires a major contract revision.
