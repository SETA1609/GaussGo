package state

import (
	"math"
	"strings"
	"time"

	"gaussgo/internal/apperrors"
	"gaussgo/internal/contracts"
)

func NormalizeAndValidate(in contracts.State) (contracts.State, []string, error) {
	// Normalization is intentionally centralized here so create/load/save flows
	// share exactly the same contract behavior.
	state := in
	warnings := make([]string, 0)

	if state.SchemaVersion == 0 {
		state.SchemaVersion = SchemaVersion
		warnings = append(warnings, "schemaVersion missing, defaulted to 1")
	}
	if state.SchemaVersion != SchemaVersion {
		return contracts.State{}, warnings, apperrors.New(apperrors.CodeValidation, apperrors.ErrorTypeInput, "unsupported schemaVersion")
	}

	if strings.TrimSpace(state.StateID) == "" {
		return contracts.State{}, warnings, apperrors.New(apperrors.CodeValidation, apperrors.ErrorTypeInput, "stateId is required")
	}
	if strings.TrimSpace(state.ProfileName) == "" {
		return contracts.State{}, warnings, apperrors.New(apperrors.CodeValidation, apperrors.ErrorTypeInput, "profileName is required")
	}

	if state.CreatedAt.IsZero() {
		state.CreatedAt = time.Now().UTC()
		warnings = append(warnings, "createdAt missing, defaulted to now")
	}
	if state.UpdatedAt.IsZero() {
		state.UpdatedAt = state.CreatedAt
		warnings = append(warnings, "updatedAt missing, defaulted to createdAt")
	}

	state.EnabledMods = ensureCoreEnabled(state.EnabledMods)
	state.UI = normalizeUI(state.UI, &warnings)
	state.Progress = normalizeProgress(state.Progress)

	for modID, p := range state.Progress {
		if p.ExerciseStats.Attempted < 0 || p.ExerciseStats.Correct < 0 {
			return contracts.State{}, warnings, apperrors.New(apperrors.CodeValidation, apperrors.ErrorTypeInput, "negative exercise stats for mod "+modID)
		}
		if p.ExerciseStats.Correct > p.ExerciseStats.Attempted {
			return contracts.State{}, warnings, apperrors.New(apperrors.CodeValidation, apperrors.ErrorTypeInput, "correct cannot exceed attempted for mod "+modID)
		}

		p.ReadConcepts = dedupeStrings(p.ReadConcepts)
		if p.ExerciseStats.Attempted == 0 {
			p.ExerciseStats.CorrectnessPercent = 0
		} else {
			pct := float64(p.ExerciseStats.Correct) / float64(p.ExerciseStats.Attempted) * 100
			p.ExerciseStats.CorrectnessPercent = math.Round(pct*100) / 100
		}
		state.Progress[modID] = p
	}

	return state, warnings, nil
}

func ensureCoreEnabled(enabled []string) []string {
	// Keep ordering stable while de-duplicating, then force core presence.
	set := make(map[string]struct{}, len(enabled)+1)
	out := make([]string, 0, len(enabled)+1)
	for _, m := range enabled {
		m = strings.TrimSpace(m)
		if m == "" {
			continue
		}
		if _, exists := set[m]; exists {
			continue
		}
		set[m] = struct{}{}
		out = append(out, m)
	}
	if _, exists := set["core"]; !exists {
		out = append([]string{"core"}, out...)
	}
	return out
}

func normalizeUI(ui contracts.UIState, warnings *[]string) contracts.UIState {
	// Locale fallback is a contract guarantee; unsupported values are tolerated
	// and normalized to keep runtime resilient.
	if ui.Preferences == nil {
		ui.Preferences = map[string]string{}
	}
	if strings.TrimSpace(ui.Locale) == "" {
		ui.Locale = DefaultLocale
		*warnings = append(*warnings, "ui.locale missing, defaulted to en")
	}
	if _, ok := SupportedLocales[ui.Locale]; !ok {
		ui.Locale = DefaultLocale
		*warnings = append(*warnings, "ui.locale unsupported, fallback to en")
	}
	return ui
}

func normalizeProgress(progress map[string]contracts.ModProgress) map[string]contracts.ModProgress {
	if progress == nil {
		return map[string]contracts.ModProgress{}
	}
	return cloneProgress(progress)
}

func dedupeStrings(in []string) []string {
	// We keep first-seen ordering to avoid surprising UI reorderings.
	seen := make(map[string]struct{}, len(in))
	out := make([]string, 0, len(in))
	for _, v := range in {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		if _, exists := seen[v]; exists {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	return out
}
