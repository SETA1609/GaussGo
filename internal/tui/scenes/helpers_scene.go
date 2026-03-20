package scenes

func buildHelpers(s RuntimeSnapshot) ViewModel {
	modID := s.CurrentModID
	if modID == "" {
		modID = "core"
	}
	lines := []string{
		tr(s, modID, "scene.helpers.help", "Helper operations from core runtime capabilities."),
	}
	if len(s.HelperCapabilities) > 0 {
		lines = append(lines, "Capabilities:")
		for _, capability := range s.HelperCapabilities {
			lines = append(lines, "- "+capability)
		}
	}
	if len(s.HelperResults) > 0 {
		lines = append(lines, "Recent helper results:")
		lines = append(lines, s.HelperResults...)
	}

	return ViewModel{
		Title:    tr(s, modID, "scene.helpers.title", "Helpers"),
		Subtitle: tr(s, modID, "scene.helpers.subtitle", "Helper operations from enabled mods"),
		Lines:    lines,
		Options: []Option{
			{Label: "Add 2 + 3", Action: Action{Type: ActionHelperAdd}},
			{Label: "Sub 7 - 4", Action: Action{Type: ActionHelperSub}},
			{Label: "Mul 6 * 5", Action: Action{Type: ActionHelperMul}},
			{Label: "Div 8 / 2", Action: Action{Type: ActionHelperDiv}},
			{Label: tr(s, modID, "common.back", "Back"), Action: Action{Type: ActionBack}},
		},
	}
}
