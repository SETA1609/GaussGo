# Development Standards

## Go Version
- Source of truth: `go.mod` (`go 1.26`, toolchain `go1.26.1`).

## Formatting and Linting
- Use `go fmt ./...` for formatting.
- Use `go vet ./...` for static checks.
- Use `make lint` and `make check` in local workflow and CI.

## Testing
- Run `make test` for package tests.
- Run `make check` before opening or updating PRs.

## Repository Conventions
- Keep contracts under `docs/contracts` with explicit version suffixes.
- Keep compile-time contract stubs under `internal/contracts`.
- Avoid breaking contract semantics without introducing a new versioned contract document.
