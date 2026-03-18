package contracts

import "time"

type ModVersionRef struct {
	ID      string
	Version string
}

type ExerciseStats struct {
	Attempted          int
	Correct            int
	CorrectnessPercent float64
}

type ModProgress struct {
	ReadConcepts    []string
	ExerciseStats   ExerciseStats
	ExerciseHistory []string
}

type UIState struct {
	Locale      string
	LastScene   string
	Preferences map[string]string
}

type State struct {
	SchemaVersion  int
	StateID        string
	ProfileName    string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	ModsAtCreation []ModVersionRef
	EnabledMods    []string
	UI             UIState
	Progress       map[string]ModProgress
}

type StateStore interface {
	Create(profileName string) (State, error)
	Load(stateID string) (State, error)
	Save(state State) error
	SetActive(stateID string) error
	GetActive() (string, error)
}
