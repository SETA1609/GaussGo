package scenes

func buildHelpers(s RuntimeSnapshot) ViewModel {
	modID := s.CurrentModID
	if modID == "" {
		modID = "core"
	}
	return ViewModel{
		Title:    tr(s, modID, "scene.helpers.title", "Helpers"),
		Subtitle: tr(s, modID, "scene.helpers.subtitle", "Helper operations from enabled mods"),
		Lines: []string{
			tr(s, modID, "scene.helpers.help", "Phase 5 will populate helper actions."),
		},
		Options: []Option{
			{Label: tr(s, modID, "common.back", "Back"), Action: Action{Type: ActionBack}},
		},
	}
}
