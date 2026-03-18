# Planning and Docs

## Planning Source of Truth
- `PLAN.md` is the top-level planning index.
- `phases/` files define execution details per phase.
- `docs/contracts/` holds versioned, authoritative contracts.

## Documentation Updates
- Update planning docs when behavior/contracts change.
- Keep examples aligned with actual schema and interfaces.
- Prefer additive changes; version breaking contract changes.

## Cross-References
- When adding new services or flows, update:
  - Phase 0 deliverables/inputs
  - Relevant phase scopes/deliverables/tests
  - `docs/integration-map.md`

## Writing Style
- Be concise and specific.
- Use stable terminology across docs (mod, concept, state, locale, quiz mode).
- Avoid ambiguous requirements; state invariants explicitly.
