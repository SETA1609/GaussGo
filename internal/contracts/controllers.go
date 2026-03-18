package contracts

type SceneController interface {
	Current() string
	Navigate(sceneID string)
	Back() error
}

type AppController interface {
	LoadActiveState() (State, error)
	CreateAndActivateState(profileName string) (State, error)
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
