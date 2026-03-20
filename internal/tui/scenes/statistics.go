package scenes

import "fmt"

func buildStatistics(s RuntimeSnapshot) ViewModel {
	modID := s.CurrentModID
	if modID == "" {
		modID = "core"
	}
	lines := []string{tr(s, modID, "scene.statistics.help", "Progress overview")}
	options := make([]Option, 0)

	for mID, progress := range s.Progress {
		if mID == "core" {
			continue
		}
		
		modLabel := fmt.Sprintf("%s: %.2f%%", mID, progress.ExerciseStats.CorrectnessPercent)
		lines = append(lines, modLabel)
		
		toggleModLabel := fmt.Sprintf("Toggle Mod: %s", mID)
		options = append(options, Option{Label: toggleModLabel, Action: Action{Type: ActionToggleStatsMod, Value: mID}})

		if s.StatsExpandedByMod[mID] {
			// Real list from file system would be better, but we can only see what's in progress or load from disk
			// For now, let's use the same logic as AppModel's sidepanel calculation (abstracted or repeated)
			// Wait, the scene builder just creates the "Main Content" part.
			// The sidepanel is rendered separately in AppModel.RenderText.
			// So SceneStatistics just needs to provide the "Controls" for the sidepanel.
			
			// If a mod is expanded, show units
			lines = append(lines, "  Units:")
			for unitID := range s.StatsExpandedByUnit[mID] {
				unitLabel := fmt.Sprintf("    %s", unitID)
				lines = append(lines, unitLabel)
				
				toggleUnitLabel := fmt.Sprintf("    Toggle Unit: %s", unitID)
				options = append(options, Option{Label: toggleUnitLabel, Action: Action{Type: ActionToggleStatsUnit, Value: mID + ":" + unitID}})
			}
		}
	}

	options = append(options, Option{Label: tr(s, modID, "common.back", "Back"), Action: Action{Type: ActionBack}})

	return ViewModel{
		Title:    tr(s, modID, "scene.statistics.title", "Statistics & Progress"),
		Subtitle: "Use options below to expand modules and units",
		Lines:    lines,
		Options:  options,
	}
}
