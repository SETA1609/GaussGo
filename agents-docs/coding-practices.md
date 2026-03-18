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

## Ports and Adapters
- Follow ports-and-adapters (hexagonal) style for boundaries that touch external systems.
- Domain/application code depends on ports (interfaces), not concrete adapters.
- Place adapters for filesystem, network, persistence, logging, and event bus behind service interfaces.
- Avoid wrapping every stdlib/helper import; abstract only volatile or boundary-facing dependencies.

## Error Handling
- Use shared error codes and error types from `internal/apperrors`.
- Do not swallow errors silently unless contract says to degrade gracefully.
- Preserve root cause when wrapping errors.
