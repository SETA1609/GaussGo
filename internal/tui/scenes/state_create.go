package scenes

func buildStateCreate(s RuntimeSnapshot) ViewModel {
	modID := "core"
	return ViewModel{
		Title:    tr(s, modID, "scene.state_create.title", "Create State"),
		Subtitle: tr(s, modID, "scene.state_create.subtitle", "Quick profile setup"),
		Lines: []string{
			tr(s, modID, "scene.state_create.help", "Create a default profile state."),
		},
		Options: []Option{
			{Label: tr(s, modID, "scene.state_create.create_default", "Create 'student'"), Action: Action{Type: ActionCreateState, Value: "student"}},
			{Label: tr(s, modID, "common.back", "Back"), Action: Action{Type: ActionBack}},
		},
	}
}
