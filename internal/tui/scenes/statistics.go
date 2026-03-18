package scenes

import "fmt"

func buildStatistics(s RuntimeSnapshot) ViewModel {
	modID := s.CurrentModID
	if modID == "" {
		modID = "core"
	}
	lines := []string{tr(s, modID, "scene.statistics.help", "Per-mod progress overview")}
	options := make([]Option, 0, len(s.Progress)+1)

	for modID, progress := range s.Progress {
		lines = append(lines, fmt.Sprintf("%s: attempted=%d correct=%d percent=%.2f", modID, progress.ExerciseStats.Attempted, progress.ExerciseStats.Correct, progress.ExerciseStats.CorrectnessPercent))
		if s.StatsExpandedByMod[modID] {
			for _, concept := range progress.ReadConcepts {
				lines = append(lines, "  - "+concept)
			}
		}
		label := tr(s, modID, "scene.statistics.toggle", "Toggle ") + modID
		options = append(options, Option{Label: label, Action: Action{Type: ActionToggleStatsMod, Value: modID}})
	}

	if len(s.Progress) == 0 {
		lines = append(lines, tr(s, modID, "scene.statistics.none", "No stats yet."))
	}

	options = append(options, Option{Label: tr(s, modID, "common.back", "Back"), Action: Action{Type: ActionBack}})

	return ViewModel{
		Title:    tr(s, modID, "scene.statistics.title", "Statistics"),
		Subtitle: tr(s, modID, "scene.statistics.subtitle", "Expand modules for concept details"),
		Lines:    lines,
		Options:  options,
	}
}
