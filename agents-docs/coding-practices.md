# Coding Practices

## General
- Keep implementations small, composable, and testable.
- Prefer clear names over short names.
- Keep package boundaries aligned with phase contracts.
- Avoid hidden global state unless explicitly required.

## Go Style
- Follow `go fmt` output.
- Keep exported APIs minimal and stable.
- Return errors with enough context to debug quickly.
- Prefer simple data structures and explicit control flow.

## Architecture
- Keep TUI scenes thin; controllers own runtime decisions.
- Use contracts in `internal/contracts` as integration boundaries.
- Keep runtime context and state schema aligned with docs/contracts.

## Error Handling
- Use shared error codes and error types from `internal/apperrors`.
- Do not swallow errors silently unless contract says to degrade gracefully.
- Preserve root cause when wrapping errors.
