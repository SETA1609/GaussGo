package scenes

import "strconv"

func buildConceptSlides(s RuntimeSnapshot) ViewModel {
	modID := s.CurrentModID
	if modID == "" {
		modID = "core"
	}

	current := s.CurrentConceptIndex + 1
	total := len(s.ConceptIDs)
	if total == 0 {
		current = 0
	}

	lines := []string{
		tr(s, modID, "scene.concepts.path", "Path: ") + s.CurrentModID + " / " + s.CurrentUnitID,
		tr(s, modID, "scene.concepts.position", "Slide ") + strconv.Itoa(current) + "/" + strconv.Itoa(total),
	}
	if s.CurrentConceptTitle != "" {
		lines = append(lines, tr(s, modID, "scene.concepts.title", "Concept: ")+s.CurrentConceptTitle)
	}
	if s.CurrentConceptBody != "" {
		lines = append(lines, s.CurrentConceptBody)
	}

	return ViewModel{
		Title:    tr(s, modID, "scene.concepts.title_main", "Concept Slides"),
		Subtitle: tr(s, modID, "scene.concepts.subtitle", "Use next/back to move through concepts"),
		Lines:    lines,
		Options: []Option{
			{Label: tr(s, modID, "scene.concepts.prev", "Previous"), Action: Action{Type: ActionConceptPrev}, Disabled: s.CurrentConceptIndex <= 0},
			{Label: tr(s, modID, "scene.concepts.next", "Next"), Action: Action{Type: ActionConceptNext}, Disabled: s.CurrentConceptIndex+1 >= total},
			{Label: tr(s, modID, "common.back", "Back"), Action: Action{Type: ActionBack}},
		},
	}
}
