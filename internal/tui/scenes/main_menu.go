package scenes

func buildMainMenu(s RuntimeSnapshot) ViewModel {
	schema, ok := s.MenuSchemas[SceneMainMenu]
	if !ok {
		schema = MenuSchema{
			TitleKey:    "scene.main_menu.title",
			Title:       "Main Menu",
			SubtitleKey: "scene.main_menu.subtitle",
			Subtitle:    "Choose an action",
			Lines: []menuLineSchema{
				{Key: "scene.main_menu.active_state", Fallback: "Active state: ", Source: "activeStateID"},
				{Key: "scene.main_menu.locale", Fallback: "Locale: ", Source: "currentLocale"},
			},
			Items: []menuItemSchema{
				{LabelKey: "menu.start_learning", Label: "Start Learning", Action: ActionNavigate, Value: SceneModSelect},
				{LabelKey: "menu.create_state", Label: "Create New State", Action: ActionNavigate, Value: SceneStateCreate},
				{LabelKey: "menu.load_state", Label: "Load Existing State", Action: ActionNavigate, Value: SceneStateLoad},
				{LabelKey: "menu.select_language", Label: "Select Language", Action: ActionNavigate, Value: SceneLanguageSelect},
				{LabelKey: "menu.mod_settings", Label: "Mod Settings", Action: ActionNavigate, Value: SceneModSettings},
				{LabelKey: "menu.exit", Label: "Exit", Action: ActionExit, Value: ""},
			},
		}
	}

	return buildFromSchema(s, "core", schema)
}
