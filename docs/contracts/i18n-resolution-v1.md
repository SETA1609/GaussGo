# i18n Resolution Contract v1

## Version
- `i18n-resolution/v1`

## Required Rules and Invariants
- Default locale is `en`.
- Active locale source of truth is `state.ui.locale`.
- Manifest must declare localization metadata (`default`, `supported`, `path`).
- Key lookup order:
  1. Active mod, active locale
  2. Core mod, active locale
  3. Active mod, default locale
  4. Core mod, default locale
  5. Key literal fallback
- Missing locale in state must normalize to `en`.

## Validation Rules
- Reject mod localization config when default locale is not in supported list.
- Reject localization path that does not exist during mod validation (or mark mod unavailable).
- Locale switch must persist immediately to active state.

## Example Lookup
Input:
- `modID=linearAlgebra`
- `locale=es`
- `key=menu.learning.randomQuiz`

Resolution sequence:
1. `mods/linearAlgebra/i18n/es.json`
2. `mods/core/i18n/es.json`
3. `mods/linearAlgebra/i18n/en.json`
4. `mods/core/i18n/en.json`
5. Return `menu.learning.randomQuiz`

## Backward Compatibility Notes
- New namespaces/keys are additive and safe.
- Changes to fallback order require a new contract version because UX text resolution changes globally.
