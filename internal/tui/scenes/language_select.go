package scenes

func buildLanguageSelect(s RuntimeSnapshot) ViewModel {
	schema, ok := s.MenuSchemas[SceneLanguageSelect]
	if !ok {
		schema = MenuSchema{
			TitleKey:    "scene.language.title",
			Title:       "Select Language",
			SubtitleKey: "scene.language.subtitle",
			Subtitle:    "Language applies to current state",
			Lines: []menuLineSchema{
				{Key: "scene.language.current", Fallback: "Current: ", Source: "currentLocale"},
			},
			Items: []menuItemSchema{
				{LabelKey: "language.english", Label: "English", Action: ActionSetLocale, Value: "en", DisabledIf: "locale_en"},
				{LabelKey: "language.spanish", Label: "Espanol", Action: ActionSetLocale, Value: "es", DisabledIf: "locale_es"},
				{LabelKey: "common.back", Label: "Back", Action: ActionBack, Value: ""},
			},
		}
	}

	return buildFromSchema(s, "core", schema)
}
