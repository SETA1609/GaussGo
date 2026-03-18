package scenes

func buildRandomQuizSelect(s RuntimeSnapshot) ViewModel {
	modID := s.CurrentModID
	if modID == "" {
		modID = "core"
	}
	concepts := sortedReadConcepts(s.Progress, s.CurrentModID)
	lines := []string{
		tr(s, modID, "scene.random_quiz.help", "Choose quiz mode."),
		tr(s, modID, "scene.random_quiz.eligible", "Eligible concepts are read concepts only."),
	}

	for _, concept := range concepts {
		prefix := "[ ] "
		if s.QuizSelectedConcepts[concept] {
			prefix = "[x] "
		}
		lines = append(lines, prefix+concept)
	}

	options := []Option{
		{Label: tr(s, modID, "scene.random_quiz.random", "Random"), Action: Action{Type: ActionSetQuizMode, Value: "random"}},
		{Label: tr(s, modID, "scene.random_quiz.selected", "Selected Concepts"), Action: Action{Type: ActionSetQuizMode, Value: "selected_concepts"}},
	}
	for _, concept := range concepts {
		options = append(options, Option{Label: concept, Action: Action{Type: ActionToggleQuizConcept, Value: concept}, Disabled: s.QuizMode != "selected_concepts"})
	}
	options = append(options, Option{Label: tr(s, modID, "common.back", "Back"), Action: Action{Type: ActionBack}})

	return ViewModel{
		Title:    tr(s, modID, "scene.random_quiz.title", "Random Quiz Setup"),
		Subtitle: tr(s, modID, "scene.random_quiz.mode", "Mode: ") + s.QuizMode,
		Lines:    lines,
		Options:  options,
	}
}
