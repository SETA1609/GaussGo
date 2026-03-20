package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gaussgo/internal/bootstrap"
	"gaussgo/internal/contracts"
	"gaussgo/internal/state"
	"gaussgo/internal/tui/scenes"
)

func TestAppModelLanguageSelectionRelabelsMainMenu(t *testing.T) {
	runtime, statesDir := setupRuntimeForTUI(t)
	model, err := NewAppModelWithModsDir(runtime, statesDir, filepath.Join(filepath.Dir(statesDir), "mods"))
	if err != nil {
		t.Fatalf("new app model failed: %v", err)
	}

	runtime.Controllers.Scene.Navigate(scenes.SceneLanguageSelect)
	model.RefreshView()
	model.SetSelectedIndex(1) // Espanol
	_ = model.Select()

	runtime.Controllers.Scene.Navigate(scenes.SceneMainMenu)
	model.RefreshView()
	text := model.RenderText()
	if !contains(text, "Menu Principal") {
		t.Fatalf("expected Spanish main menu label, got:\n%s", text)
	}
}

func TestAppModelBackNavigationReturnsOneLevel(t *testing.T) {
	runtime, statesDir := setupRuntimeForTUI(t)
	model, err := NewAppModelWithModsDir(runtime, statesDir, filepath.Join(filepath.Dir(statesDir), "mods"))
	if err != nil {
		t.Fatalf("new app model failed: %v", err)
	}

	runtime.Controllers.Scene.Navigate(scenes.SceneModSelect)
	runtime.Controllers.Scene.Navigate(scenes.SceneLearningHub)
	model.Back()

	if got := model.CurrentSceneID(); got != scenes.SceneModSelect {
		t.Fatalf("expected scene %s, got %s", scenes.SceneModSelect, got)
	}
}

func TestAppModelRefreshModsUpdatesStatuses(t *testing.T) {
	runtime, statesDir := setupRuntimeForTUI(t)
	model, err := NewAppModelWithModsDir(runtime, statesDir, filepath.Join(filepath.Dir(statesDir), "mods"))
	if err != nil {
		t.Fatalf("new app model failed: %v", err)
	}

	runtime.Controllers.Scene.Navigate(scenes.SceneModSettings)
	model.RefreshView()
	vm := model.ViewModel()
	refreshIdx := indexOfOption(vm.Options, "Refresh Mods")
	if refreshIdx < 0 {
		t.Fatalf("missing Refresh Mods option: %#v", vm.Options)
	}
	model.SetSelectedIndex(refreshIdx)
	_ = model.Select()

	text := model.RenderText()
	if !contains(text, "core") || !contains(text, "linearAlgebra") {
		t.Fatalf("expected refreshed mod rows, got:\n%s", text)
	}
}

func setupRuntimeForTUI(t *testing.T) (bootstrap.Runtime, string) {
	t.Helper()
	root := t.TempDir()
	modsDir := filepath.Join(root, "mods")
	statesDir := filepath.Join(root, "states")

	writeManifest(t, filepath.Join(modsDir, "core"), `{
		"id":"core","name":"Core","version":"0.1.0","entry":"public",
		"localization":{"default":"en","supported":["en","es"],"path":"i18n"},
		"dependencies":[],"provides":[]}`)
	writeManifest(t, filepath.Join(modsDir, "linearAlgebra"), `{
		"id":"linearAlgebra","name":"Linear Algebra","version":"0.1.0","entry":"data/units",
		"localization":{"default":"en","supported":["en","es"],"path":"i18n"},
		"dependencies":[{"id":"core","version":">=0.1.0"}],"provides":[]}`)

	if err := os.MkdirAll(filepath.Join(modsDir, "core", "ui", "scenes"), 0o755); err != nil {
		t.Fatalf("mkdir core scene dir failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(modsDir, "core", "ui", "scenes", "main_menu.json"), []byte(`{
	  "titleKey":"scene.main_menu.title",
	  "title":"Main Menu",
	  "subtitleKey":"scene.main_menu.subtitle",
	  "subtitle":"Choose an action",
	  "lines":[
		{"key":"scene.main_menu.active_state","fallback":"Active state: ","source":"activeStateID"},
		{"key":"scene.main_menu.locale","fallback":"Locale: ","source":"currentLocale"}
	  ],
	  "items":[
		{"labelKey":"menu.start_learning","label":"Start Learning","action":"navigate","value":"mod_select"},
		{"labelKey":"menu.create_state","label":"Create New State","action":"navigate","value":"state_create"},
		{"labelKey":"menu.load_state","label":"Load Existing State","action":"navigate","value":"state_load"},
		{"labelKey":"menu.select_language","label":"Select Language","action":"navigate","value":"language_select"},
		{"labelKey":"menu.mod_settings","label":"Mod Settings","action":"navigate","value":"mod_settings"},
		{"labelKey":"menu.exit","label":"Exit","action":"exit","value":""}
	  ]
	}`), 0o644); err != nil {
		t.Fatalf("write main_menu schema failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(modsDir, "core", "ui", "scenes", "language_select.json"), []byte(`{
	  "titleKey":"scene.language.title",
	  "title":"Select Language",
	  "subtitleKey":"scene.language.subtitle",
	  "subtitle":"Language applies to current state",
	  "lines":[{"key":"scene.language.current","fallback":"Current: ","source":"currentLocale"}],
	  "items":[
		{"labelKey":"language.english","label":"English","action":"set_locale","value":"en","disabledIf":"locale_en"},
		{"labelKey":"language.spanish","label":"Espanol","action":"set_locale","value":"es","disabledIf":"locale_es"},
		{"labelKey":"common.back","label":"Back","action":"back","value":""}
	  ]
	}`), 0o644); err != nil {
		t.Fatalf("write language_select schema failed: %v", err)
	}

	if err := os.MkdirAll(filepath.Join(modsDir, "core", "i18n"), 0o755); err != nil {
		t.Fatalf("mkdir core i18n dir failed: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(modsDir, "core", "i18n", "common"), 0o755); err != nil {
		t.Fatalf("mkdir core i18n common dir failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(modsDir, "core", "i18n", "common", "en.json"), []byte(`{
	  "scene.main_menu.title":"Main Menu",
	  "scene.main_menu.subtitle":"Choose an action",
	  "scene.main_menu.active_state":"Active state: ",
	  "scene.main_menu.locale":"Locale: ",
	  "menu.start_learning":"Start Learning",
	  "menu.create_state":"Create New State",
	  "menu.load_state":"Load Existing State",
	  "menu.select_language":"Select Language",
	  "menu.mod_settings":"Mod Settings",
	  "menu.exit":"Exit",
	  "scene.language.title":"Select Language",
	  "scene.language.subtitle":"Language applies to current state",
	  "scene.language.current":"Current: ",
	  "language.english":"English",
	  "language.spanish":"Espanol",
	  "common.back":"Back"
	}`), 0o644); err != nil {
		t.Fatalf("write core en i18n failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(modsDir, "core", "i18n", "common", "es.json"), []byte(`{
	  "scene.main_menu.title":"Menu Principal",
	  "scene.main_menu.subtitle":"Elige una accion",
	  "scene.main_menu.active_state":"Estado activo: ",
	  "scene.main_menu.locale":"Idioma: ",
	  "menu.start_learning":"Empezar a Aprender",
	  "menu.create_state":"Crear Nuevo Estado",
	  "menu.load_state":"Cargar Estado Existente",
	  "menu.select_language":"Seleccionar Idioma",
	  "menu.mod_settings":"Configuracion de Mods",
	  "menu.exit":"Salir",
	  "scene.language.title":"Seleccionar Idioma",
	  "scene.language.subtitle":"El idioma se aplica al estado actual",
	  "scene.language.current":"Actual: ",
	  "language.english":"Ingles",
	  "language.spanish":"Espanol",
	  "common.back":"Volver"
	}`), 0o644); err != nil {
		t.Fatalf("write core es i18n failed: %v", err)
	}

	store := state.NewStore(statesDir)
	created, err := store.Create("tester")
	if err != nil {
		t.Fatalf("create state failed: %v", err)
	}
	if err := store.SetActive(created.StateID); err != nil {
		t.Fatalf("set active failed: %v", err)
	}

	runtime, err := bootstrap.BootstrapRuntime(modsDir, statesDir, scenes.SceneMainMenu, bootstrap.DefaultServices())
	if err != nil {
		t.Fatalf("bootstrap runtime failed: %v", err)
	}
	return runtime, statesDir
}

func writeManifest(t *testing.T, dir string, content string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", dir, err)
	}
	if err := os.WriteFile(filepath.Join(dir, "manifest.json"), []byte(content), 0o644); err != nil {
		t.Fatalf("write manifest %s: %v", dir, err)
	}
}

func contains(s string, needle string) bool {
	return strings.Contains(s, needle)
}

func indexOfOption(options []scenes.Option, label string) int {
	for i, option := range options {
		if option.Label == label {
			return i
		}
	}
	return -1
}

func TestAppModelCreateStateActivatesNewState(t *testing.T) {
	runtime, statesDir := setupRuntimeForTUI(t)
	model, err := NewAppModelWithModsDir(runtime, statesDir, filepath.Join(filepath.Dir(statesDir), "mods"))
	if err != nil {
		t.Fatalf("new app model failed: %v", err)
	}
	originalStateID := runtime.Context.ActiveStateID

	runtime.Controllers.Scene.Navigate(scenes.SceneStateCreate)
	model.RefreshView()
	model.SetSelectedIndex(0) // "Create 'student'"
	_ = model.Select()

	if runtime.Context.ActiveStateID == originalStateID {
		t.Fatalf("expected ActiveStateID to change, still %s", originalStateID)
	}
	if runtime.Context.ActiveStateID == "" {
		t.Fatal("expected non-empty ActiveStateID after create")
	}
}

func TestAppModelLoadStateActivatesChosenState(t *testing.T) {
	runtime, statesDir := setupRuntimeForTUI(t)
	// create a second state to load
	store := state.NewStore(statesDir)
	second, err := store.Create("second")
	if err != nil {
		t.Fatalf("create second state: %v", err)
	}

	model, err := NewAppModelWithModsDir(runtime, statesDir, filepath.Join(filepath.Dir(statesDir), "mods"))
	if err != nil {
		t.Fatalf("new app model: %v", err)
	}

	runtime.Controllers.Scene.Navigate(scenes.SceneStateLoad)
	model.RefreshView()
	idx := indexOfOption(model.ViewModel().Options, second.StateID)
	if idx < 0 {
		t.Fatalf("second state not listed in load scene")
	}
	model.SetSelectedIndex(idx)
	_ = model.Select()

	if runtime.Context.ActiveStateID != second.StateID {
		t.Fatalf("expected ActiveStateID %s, got %s", second.StateID, runtime.Context.ActiveStateID)
	}
}

func TestAppModelStatisticsExpandCollapse(t *testing.T) {
	runtime, statesDir := setupRuntimeForTUI(t)
	model, err := NewAppModelWithModsDir(runtime, statesDir, filepath.Join(filepath.Dir(statesDir), "mods"))
	if err != nil {
		t.Fatalf("new app model: %v", err)
	}

	// seed progress with a read concept
	model.snapshot.Progress = map[string]contracts.ModProgress{
		"linearAlgebra": {
			ReadConcepts:  []string{"vectors.dot-product"},
			ExerciseStats: contracts.ExerciseStats{Attempted: 5, Correct: 4},
		},
	}
	// Initial state for nested toggle
	model.snapshot.StatsExpandedByUnit["linearAlgebra"] = map[string]bool{"vectors": false}

	// Create necessary unit files for the loader to find
	unitsDir := filepath.Join(filepath.Dir(statesDir), "mods", "linearAlgebra", "data", "units")
	if err := os.MkdirAll(unitsDir, 0o755); err != nil {
		t.Fatalf("mkdir units dir: %v", err)
	}
	unitData := `{
		"id": "vectors",
		"title": "Vectors",
		"concepts": [
			{"id": "dot-product", "title": "Dot Product", "explanation": "..."}
		]
	}`
	if err := os.WriteFile(filepath.Join(unitsDir, "vectors.json"), []byte(unitData), 0o644); err != nil {
		t.Fatalf("write unit file: %v", err)
	}

	runtime.Controllers.Scene.Navigate(scenes.SceneStatistics)
	model.RefreshView()

	// collapsed mod: unit row absent
	if strings.Contains(model.RenderText(), "Toggle Unit: vectors") {
		t.Fatal("expected unit toggle hidden before mod expand")
	}

	idx := indexOfOption(model.ViewModel().Options, "Toggle Mod: linearAlgebra")
	if idx < 0 {
		t.Fatalf("Toggle Mod: linearAlgebra option not found in %#v", model.ViewModel().Options)
	}
	model.SetSelectedIndex(idx)
	_ = model.Select()

	// expanded mod: unit row present
	text := model.RenderText()
	if !strings.Contains(text, "Toggle Unit: vectors") {
		t.Fatalf("expected unit toggle visible after mod expand, got:\n%s", text)
	}

	// toggle unit
	idx = indexOfOption(model.ViewModel().Options, "    Toggle Unit: vectors")
	if idx < 0 {
		t.Fatal("Toggle Unit: vectors option not found")
	}
	model.SetSelectedIndex(idx)
	_ = model.Select()

	// expanded unit: concept row present (in sidepanel/main lines)
	if !strings.Contains(model.RenderText(), "dot-product: 100%") {
		t.Fatal("expected concept row visible after unit expand")
	}
}

func TestAppModelQuizModeSelectedConceptsEnablesToggle(t *testing.T) {
	runtime, statesDir := setupRuntimeForTUI(t)
	model, _ := NewAppModelWithModsDir(runtime, statesDir, filepath.Join(filepath.Dir(statesDir), "mods"))
	model.snapshot.CurrentModID = "linearAlgebra"
	model.snapshot.Progress = map[string]contracts.ModProgress{
		"linearAlgebra": {ReadConcepts: []string{"dot-product", "cross-product"}},
	}
	model.snapshot.QuizMode = "random"

	runtime.Controllers.Scene.Navigate(scenes.SceneRandomQuiz)
	model.RefreshView()

	// in random mode, concept toggles are disabled
	for _, opt := range model.ViewModel().Options {
		if opt.Action.Type == scenes.ActionToggleQuizConcept && !opt.Disabled {
			t.Fatalf("expected concept toggle disabled in random mode, got enabled: %s", opt.Label)
		}
	}

	// switch to selected_concepts
	idx := indexOfOption(model.ViewModel().Options, "Selected Concepts")
	if idx < 0 {
		t.Fatal("Selected Concepts option not found")
	}
	model.SetSelectedIndex(idx)
	_ = model.Select()

	for _, opt := range model.ViewModel().Options {
		if opt.Action.Type == scenes.ActionToggleQuizConcept && opt.Disabled {
			t.Fatalf("expected concept toggle enabled in selected_concepts mode, got disabled: %s", opt.Label)
		}
	}
}

func TestAppModelHelpersRunsCoreMathCapability(t *testing.T) {
	runtime, statesDir := setupRuntimeForTUI(t)
	model, err := NewAppModelWithModsDir(runtime, statesDir, filepath.Join(filepath.Dir(statesDir), "mods"))
	if err != nil {
		t.Fatalf("new app model: %v", err)
	}

	runtime.Controllers.Scene.Navigate(scenes.SceneHelpers)
	model.RefreshView()
	idx := indexOfOption(model.ViewModel().Options, "Add 2 + 3")
	if idx < 0 {
		t.Fatal("missing helper add option")
	}
	model.SetSelectedIndex(idx)
	_ = model.Select()

	if len(model.snapshot.HelperResults) == 0 {
		t.Fatal("expected helper results to be populated")
	}
	if !strings.Contains(model.snapshot.HelperResults[0], "add(2,3) = 5") {
		t.Fatalf("unexpected helper result: %s", model.snapshot.HelperResults[0])
	}

	model.RefreshView()
	if len(model.snapshot.HelperCapabilities) == 0 {
		t.Fatal("expected helper capabilities to be listed")
	}
	found := false
	for _, capability := range model.snapshot.HelperCapabilities {
		if capability == string(contracts.CapabilityHelpersBasicMath) {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected %s capability in list, got %v", contracts.CapabilityHelpersBasicMath, model.snapshot.HelperCapabilities)
	}
}
