package tui

import (
	"context"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/NimbleMarkets/ntcharts/sparkline"

	"gaussgo/internal/bootstrap"
	"gaussgo/internal/contracts"
	"gaussgo/internal/i18n"
	"gaussgo/internal/state"
	"gaussgo/internal/tui/adapters"
	"gaussgo/internal/tui/components"
	"gaussgo/internal/tui/ports"
	"gaussgo/internal/tui/scenes"
)

type AppModel struct {
	runtime         bootstrap.Runtime
	statesDir       string
	modsDir         string
	width           int
	height          int
	selectedIdx     int
	snapshot        scenes.RuntimeSnapshot
	view            scenes.ViewModel
	spinner         ports.Spinner
	spring          ports.Spring
	zone            ports.Zone
	animPos         float64
	animVel         float64
	spark           sparkline.Model
	tick            int
	theme           components.Theme
	conceptByID     map[string]conceptFile
	inputMapper     ports.InputMapper
	renderer        ports.ViewPort
	layoutMetrics   ports.LayoutMetrics
	layoutEngine    ports.LayoutEngine
	sidepanelScroll int
}

func (m *AppModel) CurrentSceneID() string {
	return m.runtime.Context.CurrentSceneID
}

func (m *AppModel) ViewModel() scenes.ViewModel {
	return m.view
}

func (m *AppModel) RefreshView() {
	m.render()
}

func (m *AppModel) SetSelectedIndex(idx int) {
	if idx < 0 {
		m.selectedIdx = 0
		return
	}
	if idx >= len(m.view.Options) {
		if len(m.view.Options) == 0 {
			m.selectedIdx = 0
			return
		}
		m.selectedIdx = len(m.view.Options) - 1
		return
	}
	m.selectedIdx = idx
}

func NewAppModel(runtime bootstrap.Runtime, statesDir string) (*AppModel, error) {
	return NewAppModelWithModsDir(runtime, statesDir, filepath.Join(".", "mods"))
}

func NewAppModelWithModsDir(runtime bootstrap.Runtime, statesDir string, modsDir string) (*AppModel, error) {
	i18nResolver, err := i18n.NewResolver(modsDir, state.DefaultLocale)
	if err != nil {
		return nil, err
	}
	menuSchemas, err := scenes.LoadMenuSchemas(modsDir)
	if err != nil {
		return nil, err
	}

	model := &AppModel{
		runtime:       runtime,
		statesDir:     statesDir,
		modsDir:       modsDir,
		width:         100,
		height:        40,
		spinner:       adapters.NewCharmSpinner(),
		spring:        adapters.NewCharmSpring(),
		zone:          adapters.NewCharmZone(),
		spark:         sparkline.New(24, 4),
		theme:         components.DefaultTheme(),
		conceptByID:   map[string]conceptFile{},
		inputMapper:   adapters.NewCharmInputMapper(),
		renderer:      adapters.NewCharmRenderer(),
		layoutMetrics: adapters.NewCharmLipGlossMetrics(),
		layoutEngine:  adapters.NewCharmLayoutEngine(),
		snapshot: scenes.RuntimeSnapshot{
			CurrentLocale:        runtime.Context.CurrentLocale,
			I18n:                 i18nResolver,
			MenuSchemas:          menuSchemas,
			ActiveStateID:        runtime.Context.ActiveStateID,
			EnabledMods:          append([]string(nil), runtime.Context.EnabledMods...),
			CurrentModID:         runtime.Context.CurrentModID,
			CurrentUnitID:        runtime.Context.CurrentUnitID,
			Progress:             map[string]contracts.ModProgress{},
			StatsExpandedByMod:   map[string]bool{},
			StatsExpandedByUnit:  map[string]map[string]bool{},
			QuizMode:             "random",
			QuizSelectedConcepts: map[string]bool{},
			HelperResults:        nil,
			HelperCapabilities:   nil,
		},
	}

	if err := model.refreshFromState(); err != nil {
		return nil, err
	}
	if len(runtime.Context.CurrentModID) > 0 {
		if units, err := loadUnitIDs(modsDir, runtime.Context.CurrentModID); err == nil {
			model.snapshot.UnitIDs = units
		}
	}
	model.spark.PushAll([]float64{0.2, 0.3, 0.45, 0.5, 0.35, 0.6})
	model.spark.Draw()
	model.zone.Init()
	if err := model.refreshMods(); err != nil {
		model.snapshot.Flash = err.Error()
	}
	model.render()
	return model, nil
}

func (m *AppModel) Init() tea.Cmd {
	spinnerCmd, _ := m.spinner.InitCmd().(tea.Cmd)
	return tea.Batch(spinnerCmd, tickCmd())
}

// Update handles Bubble Tea messages and updates the application model accordingly.
// It manages terminal window resizing, animations, and user key inputs.
func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if event, ok := m.inputMapper.Map(msg); ok {
		if quit := m.applyInputEvent(event); quit {
			return m, tea.Quit
		}
		return m, nil
	}

	switch msg.(type) {
	case tickMsg:
		m.tick++
		next := 0.5 + 0.45*math.Sin(float64(m.tick)/4.0)
		if next < 0 {
			next = 0
		}
		m.spark.Push(next)
		m.spark.Draw()
		return m, tickCmd()
	}

	if cmd, handled := m.spinner.Update(msg); handled {
		teaCmd, ok := cmd.(tea.Cmd)
		if !ok {
			return m, nil
		}
		return m, teaCmd
	}

	return m, nil
}

func (m *AppModel) View() tea.View {
	rendered := m.renderer.Render(m.RenderText(), ports.ViewConfig{AltScreen: true, MouseMode: adapters.MouseModeCellMotion})
	v, ok := rendered.(tea.View)
	if !ok {
		fallback := tea.NewView(m.RenderText())
		fallback.AltScreen = true
		fallback.MouseMode = tea.MouseModeCellMotion
		return fallback
	}
	return v
}

func (m *AppModel) MoveUp() {
	if m.selectedIdx > 0 {
		m.selectedIdx--
	}
	m.render()
}

func (m *AppModel) MoveDown() {
	if m.selectedIdx+1 < len(m.view.Options) {
		m.selectedIdx++
	}
	m.render()
}

func (m *AppModel) Select() bool {
	if len(m.view.Options) == 0 {
		return false
	}
	selected := m.view.Options[m.selectedIdx]
	if selected.Disabled {
		return false
	}
	return m.applyAction(selected.Action)
}

func (m *AppModel) Back() bool {
	return m.applyAction(scenes.Action{Type: scenes.ActionBack})
}

// RenderText generates the final string representation of the TUI.
// It calculates responsive sizing for the main "card" based on the terminal
// window dimensions and ensures the content is correctly aligned.
func (m *AppModel) RenderText() string {
	target := float64(m.selectedIdx)
	m.animPos, m.animVel = m.spring.Update(m.animPos, m.animVel, target)

	navbar := components.RenderNavbar(m.theme, m.width, m.layoutEngine)

	footer := components.RenderFooter(m.theme, components.FooterData{
		Hint:   "Controls: j/down, k/up, enter, esc/back, [, ], q",
		User:   m.snapshot.ActiveStateID, // Using state ID as user for now
		State:  m.runtime.Context.CurrentSceneID,
		Locale: m.snapshot.CurrentLocale,
		Width:  m.width,
	}, m.layoutEngine)

	menu := make([]components.MenuItem, 0, len(m.view.Options))
	for idx, option := range m.view.Options {
		id := "option-" + strconv.Itoa(idx)
		menu = append(menu, components.MenuItem{ID: id, Label: option.Label, Disabled: option.Disabled})
	}

	sidepanelWidth := 30
	if m.width < 80 {
		sidepanelWidth = 0 // Hide sidepanel on small screens
	}
	contentWidth := m.width - sidepanelWidth

	// Traditional screen data for the main content
	content := components.RenderScreen(m.theme, components.ScreenData{
		Title:        m.view.Title,
		Subtitle:     m.view.Subtitle,
		Chart:        m.spark.View(),
		Flash:        m.snapshot.Flash,
		Lines:        m.view.Lines,
		Menu:         menu,
		Selected:     m.selectedIdx,
		WrapInCard:   true,
		MaxWidth:     contentWidth - 4,
		MaxHeight:    m.height - 4,
		FullWidth:    contentWidth,
		FullHeight:   m.height - m.layoutMetrics.Height(navbar) - m.layoutMetrics.Height(footer),
		CenterX:      true,
		CenterY:      true,
		MenuRenderer: m.zone,
		LayoutEngine: m.layoutEngine,
	})

	var sidepanel string
	if sidepanelWidth > 0 {
		sidepanel = components.RenderSidepanel(m.theme, m.calculateSidepanelData(sidepanelWidth, m.height-6))
	}

	out := components.RenderMainLayout(m.theme, navbar, footer, content, sidepanel, m.width, m.height, m.layoutEngine)
	return m.zone.Scan(out)
}

func (m *AppModel) calculateSidepanelData(width int, height int) components.SidepanelData {
	var mods []components.ModStats
	var totalPercent uint64
	var modCount int

	for modID, progress := range m.snapshot.Progress {
		if modID == "core" {
			continue
		}

		modCount++
		totalPercent += uint64(progress.ExerciseStats.CorrectnessPercent)

		mod := components.ModStats{
			ID:       modID,
			Percent:  progress.ExerciseStats.CorrectnessPercent,
			Expanded: m.snapshot.StatsExpandedByMod[modID],
		}

		if mod.Expanded {
			// In a real app, we'd load unit list and concept stats.
			// Since we don't have per-unit stats in ModProgress yet, we'll
			// mock them based on ReadConcepts for this UX demonstration.
			units, _ := loadUnitIDs(m.modsDir, modID)
			for _, unitID := range units {
				unit := components.UnitStats{
					ID:       unitID,
					Percent:  0, // Would be calculated if data existed
					Expanded: m.snapshot.StatsExpandedByUnit[modID][unitID],
				}

				if unit.Expanded {
					concepts, _, _ := loadConceptsForUnit(m.modsDir, modID, unitID)
					for _, conceptID := range concepts {
						fullID := unitID + "." + conceptID // Simplified namespacing for now, or use types.JoinNamespacedID
						concept := components.ConceptStats{
							ID:      conceptID,
							Percent: 0,
						}
						for _, read := range progress.ReadConcepts {
							if read == fullID {
								concept.Percent = 100
								break
							}
						}
						unit.Concepts = append(unit.Concepts, concept)
					}
				}
				mod.Units = append(mod.Units, unit)
			}
		}
		mods = append(mods, mod)
	}

	avg := 0.0
	if modCount > 0 {
		avg = float64(totalPercent) / float64(modCount)
	}

	return components.SidepanelData{
		AveragePercent: avg,
		Mods:           mods,
		Width:          width,
		Height:         height,
		ScrollOffset:   m.sidepanelScroll,
	}
}

type tickMsg struct{}

func tickCmd() tea.Cmd {
	return tea.Tick(80*time.Millisecond, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (m *AppModel) applyAction(action scenes.Action) bool {
	m.snapshot.Flash = ""
	switch action.Type {
	case scenes.ActionNone:
		return false
	case scenes.ActionNavigate:
		m.runtime.Controllers.Scene.Navigate(action.Value)
		m.selectedIdx = 0
	case scenes.ActionBack:
		if err := m.runtime.Controllers.Scene.Back(); err != nil {
			m.snapshot.Flash = err.Error()
		}
		m.selectedIdx = 0
	case scenes.ActionCreateState:
		st, err := m.runtime.Controllers.App.CreateAndActivateState(action.Value)
		if err != nil {
			m.snapshot.Flash = err.Error()
			break
		}
		_ = m.refreshFromState()
		m.snapshot.Flash = "state created: " + st.StateID
	case scenes.ActionLoadState:
		if err := m.runtime.Controllers.App.ActivateState(action.Value); err != nil {
			m.snapshot.Flash = err.Error()
			break
		}
		_ = m.refreshFromState()
		m.snapshot.Flash = "state loaded: " + action.Value
	case scenes.ActionSetLocale:
		if err := m.runtime.Controllers.Locale.SetLocale(action.Value); err != nil {
			m.snapshot.Flash = err.Error()
			break
		}
		m.snapshot.CurrentLocale = action.Value
		_ = m.refreshFromState()
	case scenes.ActionRefreshMods:
		if err := m.refreshMods(); err != nil {
			m.snapshot.Flash = err.Error()
		} else {
			m.snapshot.Flash = "mods refreshed"
		}
	case scenes.ActionToggleMod:
		if isEnabled(m.snapshot.ModStatuses, action.Value) {
			if err := m.runtime.Controllers.Mod.Disable(action.Value); err != nil {
				m.snapshot.Flash = err.Error()
				break
			}
		} else {
			if err := m.runtime.Controllers.Mod.Enable(action.Value); err != nil {
				m.snapshot.Flash = err.Error()
				break
			}
		}
		if err := m.refreshMods(); err != nil {
			m.snapshot.Flash = err.Error()
		}
		_ = m.refreshFromState()
	case scenes.ActionSelectMod:
		if err := m.runtime.Controllers.Learning.SelectMod(action.Value); err != nil {
			m.snapshot.Flash = err.Error()
			break
		}
		m.snapshot.CurrentModID = action.Value
		units, unitErr := loadUnitIDs(m.modsDir, action.Value)
		if unitErr != nil {
			m.snapshot.UnitIDs = nil
			m.snapshot.Flash = unitErr.Error()
		} else {
			m.snapshot.UnitIDs = units
		}
		m.snapshot.CurrentUnitID = ""
		m.snapshot.ConceptIDs = nil
		m.snapshot.CurrentConceptIndex = 0
		m.snapshot.CurrentConceptTitle = ""
		m.snapshot.CurrentConceptBody = ""
		m.runtime.Controllers.Scene.Navigate(scenes.SceneLearningHub)
		m.selectedIdx = 0
	case scenes.ActionSelectUnit:
		if err := m.runtime.Controllers.Learning.SelectUnit(action.Value); err != nil {
			m.snapshot.Flash = err.Error()
			break
		}
		m.snapshot.CurrentUnitID = action.Value
		conceptIDs, conceptByID, loadErr := loadConceptsForUnit(m.modsDir, m.snapshot.CurrentModID, action.Value)
		if loadErr != nil {
			m.snapshot.Flash = loadErr.Error()
			break
		}
		m.snapshot.ConceptIDs = conceptIDs
		m.snapshot.CurrentConceptIndex = 0
		m.conceptByID = conceptByID
		m.updateCurrentConceptSlide()
		m.runtime.Controllers.Scene.Navigate(scenes.SceneConceptSlides)
		m.selectedIdx = 0
	case scenes.ActionConceptNext:
		if m.snapshot.CurrentConceptIndex+1 < len(m.snapshot.ConceptIDs) {
			m.snapshot.CurrentConceptIndex++
			m.updateCurrentConceptSlide()
		}
	case scenes.ActionConceptPrev:
		if m.snapshot.CurrentConceptIndex > 0 {
			m.snapshot.CurrentConceptIndex--
			m.updateCurrentConceptSlide()
		}
	case scenes.ActionToggleStatsMod:
		m.snapshot.StatsExpandedByMod[action.Value] = !m.snapshot.StatsExpandedByMod[action.Value]
	case scenes.ActionToggleStatsUnit:
		parts := strings.Split(action.Value, ":")
		if len(parts) == 2 {
			modID, unitID := parts[0], parts[1]
			if m.snapshot.StatsExpandedByUnit[modID] == nil {
				m.snapshot.StatsExpandedByUnit[modID] = map[string]bool{}
			}
			m.snapshot.StatsExpandedByUnit[modID][unitID] = !m.snapshot.StatsExpandedByUnit[modID][unitID]
		}
	case scenes.ActionSetQuizMode:
		m.snapshot.QuizMode = action.Value
	case scenes.ActionToggleQuizConcept:
		if m.snapshot.QuizMode == "selected_concepts" {
			m.snapshot.QuizSelectedConcepts[action.Value] = !m.snapshot.QuizSelectedConcepts[action.Value]
		}
	case scenes.ActionExit:
		return true
	case scenes.ActionHelperAdd:
		m.runMathHelper("add", 2, 3)
	case scenes.ActionHelperSub:
		m.runMathHelper("sub", 7, 4)
	case scenes.ActionHelperMul:
		m.runMathHelper("mul", 6, 5)
	case scenes.ActionHelperDiv:
		m.runMathHelper("div", 8, 2)
	}

	m.render()
	return false
}

func (m *AppModel) runMathHelper(op string, a float64, b float64) {
	if m.runtime.ModRegistry == nil {
		m.snapshot.Flash = "runtime mod registry unavailable"
		return
	}

	coreMod, err := m.runtime.ModRegistry.Get("core")
	if err != nil {
		m.snapshot.Flash = "core runtime mod unavailable"
		return
	}

	provider, ok := coreMod.(contracts.CapabilityProvider)
	if !ok {
		m.snapshot.Flash = "core mod has no capability provider"
		return
	}

	raw, exists := provider.Capability(contracts.CapabilityHelpersBasicMath)
	if !exists {
		m.snapshot.Flash = "basic math helper capability missing"
		return
	}

	mathHelpers, ok := raw.(contracts.BasicMathHelpers)
	if !ok {
		m.snapshot.Flash = "invalid basic math helper capability type"
		return
	}

	ctx := context.Background()
	var value float64
	switch op {
	case "add":
		value, err = mathHelpers.Add(ctx, a, b)
	case "sub":
		value, err = mathHelpers.Sub(ctx, a, b)
	case "mul":
		value, err = mathHelpers.Mul(ctx, a, b)
	case "div":
		value, err = mathHelpers.Div(ctx, a, b)
	default:
		m.snapshot.Flash = "unknown helper operation"
		return
	}
	if err != nil {
		m.snapshot.Flash = err.Error()
		return
	}

	line := op + "(" + strconv.FormatFloat(a, 'f', -1, 64) + "," + strconv.FormatFloat(b, 'f', -1, 64) + ") = " + strconv.FormatFloat(value, 'f', -1, 64)
	m.snapshot.HelperResults = append([]string{line}, m.snapshot.HelperResults...)
	if len(m.snapshot.HelperResults) > 5 {
		m.snapshot.HelperResults = m.snapshot.HelperResults[:5]
	}
	m.snapshot.Flash = "helper executed"
}

func (m *AppModel) render() {
	m.snapshot.CurrentLocale = m.runtime.Context.CurrentLocale
	m.snapshot.ActiveStateID = m.runtime.Context.ActiveStateID
	m.snapshot.EnabledMods = append([]string(nil), m.runtime.Context.EnabledMods...)
	m.snapshot.CurrentModID = m.runtime.Context.CurrentModID
	m.snapshot.CurrentUnitID = m.runtime.Context.CurrentUnitID
	m.snapshot.HelperCapabilities = m.collectHelperCapabilities()
	m.snapshot.AvailableStates = listStateIDs(m.statesDir)
	m.view = scenes.Build(m.runtime.Context.CurrentSceneID, m.snapshot)
	if m.selectedIdx >= len(m.view.Options) {
		m.selectedIdx = 0
	}
}

func (m *AppModel) collectHelperCapabilities() []string {
	if m.runtime.ModRegistry == nil {
		return nil
	}

	registered := m.runtime.ModRegistry.List()
	capabilitySet := map[string]struct{}{}
	for _, runtimeMod := range registered {
		if catalog, ok := runtimeMod.(contracts.CapabilityCatalog); ok {
			for _, capability := range catalog.Capabilities() {
				capabilitySet[string(capability)] = struct{}{}
			}
			continue
		}

		if provider, ok := runtimeMod.(contracts.CapabilityProvider); ok {
			if provider.HasCapability(contracts.CapabilityHelpersBasicMath) {
				capabilitySet[string(contracts.CapabilityHelpersBasicMath)] = struct{}{}
			}
		}
	}

	capabilities := make([]string, 0, len(capabilitySet))
	for capability := range capabilitySet {
		capabilities = append(capabilities, capability)
	}
	sort.Strings(capabilities)
	return capabilities
}

func (m *AppModel) refreshMods() error {
	statuses, err := m.runtime.Controllers.Mod.Refresh()
	if err != nil {
		return err
	}
	m.snapshot.ModStatuses = statuses
	return nil
}

func (m *AppModel) refreshFromState() error {
	st, err := m.runtime.Controllers.App.LoadActiveState()
	if err != nil {
		return err
	}
	m.snapshot.Progress = st.Progress
	m.snapshot.CurrentLocale = st.UI.Locale
	m.snapshot.ActiveStateID = st.StateID
	m.snapshot.EnabledMods = append([]string(nil), st.EnabledMods...)
	return nil
}

func isEnabled(statuses []contracts.ModStatus, modID string) bool {
	for _, status := range statuses {
		if status.ModID == modID {
			return status.Enabled
		}
	}
	return false
}

func listStateIDs(statesDir string) []string {
	entries, err := os.ReadDir(statesDir)
	if err != nil {
		return nil
	}

	ids := make([]string, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasPrefix(name, "state-") || !strings.HasSuffix(name, ".json") {
			continue
		}
		base := strings.TrimSuffix(strings.TrimPrefix(name, "state-"), ".json")
		if strings.TrimSpace(base) != "" {
			ids = append(ids, base)
		}
	}
	sort.Strings(ids)
	return ids
}

func (m *AppModel) updateCurrentConceptSlide() {
	if len(m.snapshot.ConceptIDs) == 0 {
		m.snapshot.CurrentConceptTitle = ""
		m.snapshot.CurrentConceptBody = ""
		return
	}
	idx := m.snapshot.CurrentConceptIndex
	if idx < 0 {
		idx = 0
	}
	if idx >= len(m.snapshot.ConceptIDs) {
		idx = len(m.snapshot.ConceptIDs) - 1
	}
	m.snapshot.CurrentConceptIndex = idx
	conceptID := m.snapshot.ConceptIDs[idx]
	concept := m.conceptByID[conceptID]
	if strings.TrimSpace(concept.Title) == "" {
		m.snapshot.CurrentConceptTitle = conceptID
	} else {
		m.snapshot.CurrentConceptTitle = concept.Title
	}
	m.snapshot.CurrentConceptBody = concept.Explanation
}

func (m *AppModel) applyInputEvent(event ports.InputEvent) bool {
	switch event.Type {
	case ports.InputEventWindowSize:
		m.width = event.Width
		m.height = event.Height
		m.spark.Resize(maxInt(12, minInt(m.width-20, 40)), 4)
		m.spark.Draw()
	case ports.InputEventKey:
		switch event.Key {
		case ports.KeyUp, ports.KeyK:
			m.MoveUp()
		case ports.KeyDown, ports.KeyJ:
			m.MoveDown()
		case ports.KeyEnter:
			if quit := m.Select(); quit {
				return true
			}
		case ports.KeyEsc, ports.KeyBackspace:
			m.Back()
		case ports.KeyQ, ports.KeyCtrlC:
			return true
		case ports.KeyLeftBracket:
			if m.sidepanelScroll > 0 {
				m.sidepanelScroll--
			}
		case ports.KeyRightBracket:
			m.sidepanelScroll++
		}
	}

	return false
}

func maxInt(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func minInt(a int, b int) int {
	if a < b {
		return a
	}
	return b
}
