package scenes

func buildStateLoad(s RuntimeSnapshot) ViewModel {
	modID := "core"
	options := make([]Option, 0, len(s.AvailableStates)+1)
	for _, id := range s.AvailableStates {
		options = append(options, Option{Label: id, Action: Action{Type: ActionLoadState, Value: id}})
	}
	options = append(options, Option{Label: tr(s, modID, "common.back", "Back"), Action: Action{Type: ActionBack}})

	lines := []string{tr(s, modID, "scene.state_load.help", "Select a state to load.")}
	if len(s.AvailableStates) == 0 {
		lines = append(lines, tr(s, modID, "scene.state_load.none", "No states found."))
	}

	return ViewModel{
		Title:    tr(s, modID, "scene.state_load.title", "Load State"),
		Subtitle: tr(s, modID, "scene.state_load.subtitle", "Switch active profile"),
		Lines:    lines,
		Options:  options,
	}
}
