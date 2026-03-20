package contracts

type SceneController interface {
	Current() string
	Navigate(sceneID string)
	Back() error
}

// AppController manages the high-level application state, including loading,
// creating, and activating learner profiles (states).
type AppController interface {
	// LoadActiveState loads the current active state from the store.
	// If no state is active, it attempts to resolve the default active state.
	LoadActiveState() (State, error)

	// CreateAndActivateState creates a new state with the given profile name,
	// sets it as active, and returns the newly loaded state. This method
	// ensures that state creation events are emitted and logged.
	CreateAndActivateState(profileName string) (State, error)

	// ActivateState sets the given stateID as the active state and reloads
	// the runtime context to reflect the new state.
	ActivateState(stateID string) error
}

type LearningController interface {
	SelectMod(modID string) error
	SelectUnit(unitID string) error
	SelectConcept(conceptID string) error
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
