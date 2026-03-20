package scenes

func buildUnitSelect(s RuntimeSnapshot) ViewModel {
	modID := s.CurrentModID
	if modID == "" {
		modID = "core"
	}

	options := make([]Option, 0, len(s.UnitIDs)+1)
	for _, unitID := range s.UnitIDs {
		options = append(options, Option{Label: unitID, Action: Action{Type: ActionSelectUnit, Value: unitID}})
	}
	options = append(options, Option{Label: tr(s, "core", "common.back", "Back"), Action: Action{Type: ActionBack}})

	return ViewModel{
		Title:    tr(s, "core", "scene.unit_select.title", "Select Unit"),
		Subtitle: tr(s, "core", "scene.unit_select.subtitle", "Select a unit from current mod"),
		Lines: []string{
			tr(s, "core", "scene.unit_select.current_mod", "Current mod: ") + modID,
		},
		Options: options,
	}
}
