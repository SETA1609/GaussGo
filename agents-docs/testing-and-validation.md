# Testing and Validation

## Required Local Checks
- Run `make lint` before pushing.
- Run `make test` for changed packages.
- Run `make check` before opening/updating PRs.

## Test Scope
- Cover behavior changes with focused unit tests.
- Prefer deterministic tests over timing-dependent tests.
- Validate contract-boundary behavior (state, controllers, i18n, events).

## Validation Notes
- Keep fixtures minimal and readable.
- Add regression tests for fixed bugs.
- If tests are deferred, document why and add follow-up task in phase plan.
