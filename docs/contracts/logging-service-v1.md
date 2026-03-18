# Logging Service Contract v1

## Version
- `logging-service/v1`

## Required Interface and Invariants
- Logging service is available to bootstrap, controllers, and runtime services.
- Minimum levels: `debug`, `info`, `warn`, `error`.
- Log entries are structured with message and context fields.
- Logging must never crash runtime flow; logger errors are swallowed or surfaced as warnings.

## Validation Rules
- Reject unknown levels in strict mode; map to `info` in permissive mode.
- Context values must be serializable primitives/maps/slices.

## Example Interface (Go)
```go
type LoggingService interface {
    Debug(msg string, fields map[string]any)
    Info(msg string, fields map[string]any)
    Warn(msg string, fields map[string]any)
    Error(msg string, fields map[string]any)
}
```

## Example Usage
```text
level=info msg="state loaded" stateId=default locale=en
level=warn msg="locale fallback" requested=fr fallback=en
```

## Backward Compatibility Notes
- Adding optional helper methods is non-breaking.
- Removing or renaming level methods requires a major contract revision.
