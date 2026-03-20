package scenes

func buildLearningHub(s RuntimeSnapshot) ViewModel {
	modID := s.CurrentModID
	if modID == "" {
		modID = "core"
	}
	return ViewModel{
		Title:    tr(s, modID, "scene.learning_hub.title", "Learning Hub"),
		Subtitle: tr(s, modID, "scene.learning_hub.subtitle", "Choose activity"),
		Lines: []string{
			tr(s, modID, "scene.learning_hub.mod", "Mod: ") + s.CurrentModID,
		},
		Options: []Option{
			{Label: tr(s, modID, "scene.learning_hub.statistics", "Statistics"), Action: Action{Type: ActionNavigate, Value: SceneStatistics}},
			{Label: tr(s, modID, "scene.learning_hub.lessons", "Lessons"), Action: Action{Type: ActionNavigate, Value: SceneUnitSelect}, Disabled: len(s.UnitIDs) == 0},
			{Label: tr(s, modID, "scene.learning_hub.random_quiz", "Random Quiz"), Action: Action{Type: ActionNavigate, Value: SceneRandomQuiz}},
			{Label: tr(s, modID, "scene.learning_hub.helpers", "Helpers"), Action: Action{Type: ActionNavigate, Value: SceneHelpers}},
			{Label: tr(s, modID, "scene.learning_hub.back_main", "Back to Main Menu"), Action: Action{Type: ActionNavigate, Value: SceneMainMenu}},
		},
	}
}
