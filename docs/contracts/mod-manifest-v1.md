# Mod Manifest Contract v1

## Version
- `manifest/v1`

## Required Fields and Invariants
- Required top-level fields: `id`, `name`, `version`, `entry`, `localization`, `dependencies`, `provides`.
- `id` is stable, lowercase/alpha-numeric with optional internal dashes or camel-case segments.
- `id` must match the mod directory name.
- `version` follows semantic versioning (`MAJOR.MINOR.PATCH`).
- `localization` object is mandatory with:
  - `default` locale code
  - `supported` non-empty locale list
  - `path` directory relative to mod root
- `localization.default` must be present in `localization.supported`.
- `core` mod must exist in runtime and cannot depend on other mods.
- Dependency list entries contain `id` and semver range `version`.

## Validation Rules
- Reject manifests missing any required field.
- Reject invalid or empty `id`.
- Reject invalid semver in `version` or dependency ranges.
- Reject missing/invalid `localization.default`, `localization.supported`, `localization.path`.
- Reject manifests where dependency `id == manifest.id`.
- Reject unresolved dependencies at graph resolution time.

## Example Payload
```json
{
  "id": "linearAlgebra",
  "name": "Linear Algebra",
  "version": "0.1.0",
  "entry": "data/units",
  "localization": {
    "default": "en",
    "supported": ["en", "es"],
    "path": "i18n"
  },
  "dependencies": [
    { "id": "core", "version": ">=0.1.0" }
  ],
  "provides": ["learning.units", "exercise.compute"]
}
```

## Backward Compatibility Notes
- New optional fields are allowed in minor revisions and must be ignored by unknown-field tolerant parsers.
- Removing or changing meaning of required fields requires a new major contract version.
- Existing mods without `localization` are incompatible with `manifest/v1` and must be migrated.
