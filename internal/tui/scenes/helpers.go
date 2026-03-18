package scenes

import (
	"sort"

	"gaussgo/internal/contracts"
)

const (
	SceneMainMenu       = "main_menu"
	SceneStateCreate    = "state_create"
	SceneStateLoad      = "state_load"
	SceneLanguageSelect = "language_select"
	SceneModSettings    = "mod_settings"
	SceneModSelect      = "mod_select"
	SceneUnitSelect     = "unit_select"
	SceneLearningHub    = "learning_hub"
	SceneConceptSlides  = "concept_slides"
	SceneStatistics     = "statistics"
	SceneRandomQuiz     = "random_quiz_select"
	SceneHelpers        = "helpers"
)

const (
	ActionNone              = "none"
	ActionNavigate          = "navigate"
	ActionBack              = "back"
	ActionCreateState       = "create_state"
	ActionLoadState         = "load_state"
	ActionSetLocale         = "set_locale"
	ActionRefreshMods       = "refresh_mods"
	ActionToggleMod         = "toggle_mod"
	ActionSelectMod         = "select_mod"
	ActionSelectUnit        = "select_unit"
	ActionConceptNext       = "concept_next"
	ActionConceptPrev       = "concept_prev"
	ActionToggleStatsMod    = "toggle_stats_mod"
	ActionSetQuizMode       = "set_quiz_mode"
	ActionToggleQuizConcept = "toggle_quiz_concept"
	ActionExit              = "exit"
)

type Action struct {
	Type  string
	Value string
}

type Option struct {
	Label    string
	Action   Action
	Disabled bool
}

type ViewModel struct {
	Title    string
	Subtitle string
	Lines    []string
	Options  []Option
}

type RuntimeSnapshot struct {
	CurrentLocale        string
	I18n                 contracts.I18nResolver
	MenuSchemas          map[string]MenuSchema
	ActiveStateID        string
	EnabledMods          []string
	CurrentModID         string
	CurrentUnitID        string
	UnitIDs              []string
	ConceptIDs           []string
	CurrentConceptIndex  int
	CurrentConceptTitle  string
	CurrentConceptBody   string
	ModStatuses          []contracts.ModStatus
	AvailableStates      []string
	Progress             map[string]contracts.ModProgress
	StatsExpandedByMod   map[string]bool
	QuizMode             string
	QuizSelectedConcepts map[string]bool
	Flash                string
}

func Build(sceneID string, s RuntimeSnapshot) ViewModel {
	switch sceneID {
	case SceneMainMenu:
		return buildMainMenu(s)
	case SceneStateCreate:
		return buildStateCreate(s)
	case SceneStateLoad:
		return buildStateLoad(s)
	case SceneLanguageSelect:
		return buildLanguageSelect(s)
	case SceneModSettings:
		return buildModSettings(s)
	case SceneModSelect:
		return buildModSelect(s)
	case SceneUnitSelect:
		return buildUnitSelect(s)
	case SceneLearningHub:
		return buildLearningHub(s)
	case SceneConceptSlides:
		return buildConceptSlides(s)
	case SceneStatistics:
		return buildStatistics(s)
	case SceneRandomQuiz:
		return buildRandomQuizSelect(s)
	case SceneHelpers:
		return buildHelpers(s)
	default:
		return ViewModel{
			Title:   "GaussGo",
			Lines:   []string{"Unknown scene: " + sceneID},
			Options: []Option{{Label: "Back", Action: Action{Type: ActionBack}}},
		}
	}
}

func tr(s RuntimeSnapshot, modID string, key string, fallback string) string {
	if s.I18n == nil {
		return fallback
	}
	resolved := s.I18n.Resolve(modID, s.CurrentLocale, key)
	if resolved == key {
		return fallback
	}
	return resolved
}

func sortedReadConcepts(progress map[string]contracts.ModProgress, modID string) []string {
	p, ok := progress[modID]
	if !ok {
		return nil
	}
	out := append([]string(nil), p.ReadConcepts...)
	sort.Strings(out)
	return out
}
