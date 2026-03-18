# Git Style

## Commit Principles
- Make atomic commits focused on one concern.
- Use imperative, concise commit messages.
- Prefer split commits for docs/contracts vs refactors vs features.

## Commit Message Pattern
- `<type>: <short purpose>`
- Common types: `feat`, `fix`, `refactor`, `docs`, `chore`, `test`.

## Branching
- Use descriptive branch names for scoped work.
- Keep branches up to date and avoid unnecessary merge noise.

## Safety
- Never commit secrets or local credential files.
- Avoid destructive git commands unless explicitly requested.
- Do not amend or force-push shared history without clear intent.

## Co-Author Policy
- Always include co-author line for `SETA1609` in commit messages.
- Use: `Co-authored-by: SETA1609 <SETA1609@users.noreply.github.com>`.
