package state

import (
	"time"

	"gaussgo/internal/contracts"
)

const (
	SchemaVersion = 1
	DefaultLocale = "en"
)

var SupportedLocales = map[string]struct{}{
	"en": {},
	"es": {},
}

func SupportedLocalesSet() map[string]struct{} {
	out := make(map[string]struct{}, len(SupportedLocales))
	for locale := range SupportedLocales {
		out[locale] = struct{}{}
	}
	return out
}

type fileState struct {
	SchemaVersion  int                              `json:"schemaVersion"`
	StateID        string                           `json:"stateId"`
	ProfileName    string                           `json:"profileName"`
	CreatedAt      string                           `json:"createdAt"`
	UpdatedAt      string                           `json:"updatedAt"`
	ModsAtCreation []contracts.ModVersionRef        `json:"modsAtCreation"`
	EnabledMods    []string                         `json:"enabledMods"`
	UI             fileUIState                      `json:"ui"`
	Progress       map[string]contracts.ModProgress `json:"progress"`
}

type fileUIState struct {
	Locale      string            `json:"locale"`
	LastScene   string            `json:"lastScene,omitempty"`
	Preferences map[string]string `json:"preferences,omitempty"`
}

func toFileState(in contracts.State) fileState {
	// Persist timestamps with nanosecond precision to avoid accidental
	// equality across fast consecutive saves.
	return fileState{
		SchemaVersion:  in.SchemaVersion,
		StateID:        in.StateID,
		ProfileName:    in.ProfileName,
		CreatedAt:      in.CreatedAt.UTC().Format(time.RFC3339Nano),
		UpdatedAt:      in.UpdatedAt.UTC().Format(time.RFC3339Nano),
		ModsAtCreation: append([]contracts.ModVersionRef(nil), in.ModsAtCreation...),
		EnabledMods:    append([]string(nil), in.EnabledMods...),
		UI: fileUIState{
			Locale:      in.UI.Locale,
			LastScene:   in.UI.LastScene,
			Preferences: cloneStringMap(in.UI.Preferences),
		},
		Progress: cloneProgress(in.Progress),
	}
}

func fromFileState(in fileState) (contracts.State, error) {
	createdAt, err := parseStateTime(in.CreatedAt)
	if err != nil {
		return contracts.State{}, err
	}
	updatedAt, err := parseStateTime(in.UpdatedAt)
	if err != nil {
		return contracts.State{}, err
	}

	return contracts.State{
		SchemaVersion:  in.SchemaVersion,
		StateID:        in.StateID,
		ProfileName:    in.ProfileName,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
		ModsAtCreation: append([]contracts.ModVersionRef(nil), in.ModsAtCreation...),
		EnabledMods:    append([]string(nil), in.EnabledMods...),
		UI: contracts.UIState{
			Locale:      in.UI.Locale,
			LastScene:   in.UI.LastScene,
			Preferences: cloneStringMap(in.UI.Preferences),
		},
		Progress: cloneProgress(in.Progress),
	}, nil
}

func parseStateTime(in string) (time.Time, error) {
	// Backward-compatible parsing: prefer RFC3339Nano but accept RFC3339.
	if t, err := time.Parse(time.RFC3339Nano, in); err == nil {
		return t, nil
	}
	return time.Parse(time.RFC3339, in)
}

func cloneStringMap(in map[string]string) map[string]string {
	if in == nil {
		return map[string]string{}
	}
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func cloneProgress(in map[string]contracts.ModProgress) map[string]contracts.ModProgress {
	if in == nil {
		return map[string]contracts.ModProgress{}
	}
	out := make(map[string]contracts.ModProgress, len(in))
	for modID, p := range in {
		out[modID] = contracts.ModProgress{
			ReadConcepts:    append([]string(nil), p.ReadConcepts...),
			ExerciseStats:   p.ExerciseStats,
			ExerciseHistory: append([]string(nil), p.ExerciseHistory...),
		}
	}
	return out
}
