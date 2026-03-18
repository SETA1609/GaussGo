# Controller Interfaces Contract v1

## Version
- `controllers/v1`

## Required Interfaces and Invariants
- `SceneController` owns scene transitions and one-step back behavior.
- `ModController` toggles mod enablement but must reject disabling `core`.
- `LocaleController` reads/writes locale and persists to active state.
- `LearningController` manages selected mod/unit/concept and quiz setup context.
- Every menu flow must expose a `Back` path to the immediate parent scene.

## Validation Rules
- `SceneController.Back()` returns error when stack is empty.
- `ModController.Disable("core")` returns contract error.
- `LocaleController.SetLocale(locale)` validates locale support before persistence.

## Example Signatures (Go)
```go
type SceneController interface {
    Current() string
    Navigate(sceneID string)
    Back() error
}

type ModController interface {
    Enable(modID string) error
    Disable(modID string) error
    Refresh() ([]ModStatus, error)
}

type LocaleController interface {
    Current() string
    SetLocale(locale string) error
}

type LearningController interface {
    SelectMod(modID string) error
    SelectUnit(unitID string) error
    SelectConcept(conceptID string) error
}
```

## Backward Compatibility Notes
- Adding methods to interfaces is breaking; prefer adding secondary optional interfaces.
- Error semantics should remain stable across v1 patch updates.
