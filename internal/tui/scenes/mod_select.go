package scenes

func buildModSelect(s RuntimeSnapshot) ViewModel {
	modID := "core"
	options := make([]Option, 0, len(s.ModStatuses)+1)
	for _, status := range s.ModStatuses {
		if status.Enabled {
			options = append(options, Option{Label: status.ModID, Action: Action{Type: ActionSelectMod, Value: status.ModID}})
		}
	}
	options = append(options, Option{Label: tr(s, modID, "common.back", "Back"), Action: Action{Type: ActionBack}})

	return ViewModel{
		Title:    tr(s, modID, "scene.mod_select.title", "Select Mod"),
		Subtitle: tr(s, modID, "scene.mod_select.subtitle", "Choose module to study"),
		Lines: []string{
			tr(s, modID, "scene.mod_select.help", "Enabled modules only."),
		},
		Options: options,
	}
}
