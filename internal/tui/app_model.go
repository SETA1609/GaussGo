package tui

import (
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	bubblespinner "charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"github.com/NimbleMarkets/ntcharts/sparkline"
	"github.com/charmbracelet/harmonica"
	zone "github.com/lrstanley/bubblezone/v2"

	"gaussgo/internal/bootstrap"
	"gaussgo/internal/contracts"
	"gaussgo/internal/i18n"
	"gaussgo/internal/state"
	"gaussgo/internal/tui/components"
	"gaussgo/internal/tui/scenes"
)

type AppModel struct {
	runtime     bootstrap.Runtime
	statesDir   string
	modsDir     string
	width       int
	height      int
	selectedIdx int
	snapshot    scenes.RuntimeSnapshot
	view        scenes.ViewModel
	spinner     bubblespinner.Model
	spring      harmonica.Spring
	animPos     float64
	animVel     float64
	spark       sparkline.Model
	tick        int
	theme       components.Theme
	conceptByID map[string]conceptFile
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
		runtime:     runtime,
		statesDir:   statesDir,
		modsDir:     modsDir,
		width:       100,
		height:      40,
		spinner:     bubblespinner.New(),
		spring:      harmonica.NewSpring(harmonica.FPS(60), 8.0, 0.7),
		spark:       sparkline.New(24, 4),
		theme:       components.DefaultTheme(),
		conceptByID: map[string]conceptFile{},
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
			QuizMode:             "random",
			QuizSelectedConcepts: map[string]bool{},
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
	zone.NewGlobal()
	if err := model.refreshMods(); err != nil {
		model.snapshot.Flash = err.Error()
	}
	model.render()
	return model, nil
}

func (m *AppModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, tickCmd())
}

func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		m.tick++
		next := 0.5 + 0.45*math.Sin(float64(m.tick)/4.0)
		if next < 0 {
			next = 0
		}
		m.spark.Push(next)
		m.spark.Draw()
		return m, tickCmd()
	case bubblespinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.MoveUp()
		case "down", "j":
			m.MoveDown()
		case "enter":
			if quit := m.Select(); quit {
				return m, tea.Quit
			}
		case "esc", "backspace":
			m.Back()
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.spark.Resize(maxInt(12, minInt(m.width-20, 40)), 4)
		m.spark.Draw()
	}

	return m, nil
}

func (m *AppModel) View() tea.View {
	v := tea.NewView(m.RenderText())
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
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

func (m *AppModel) RenderText() string {
	target := float64(m.selectedIdx)
	m.animPos, m.animVel = m.spring.Update(m.animPos, m.animVel, target)

	header := m.spinner.View() + "  scene=" + m.runtime.Context.CurrentSceneID + "  locale=" + m.snapshot.CurrentLocale
	header = components.RenderStatusPill(m.theme, header)

	menu := make([]components.MenuItem, 0, len(m.view.Options))
	for idx, option := range m.view.Options {
		id := "option-" + strconv.Itoa(idx)
		menu = append(menu, components.MenuItem{ID: id, Label: option.Label, Disabled: option.Disabled})
	}
	maxWidth := int(float64(m.width) * 0.8)
	if m.width < 60 {
		maxWidth = m.width - 2
	}
	if maxWidth < 40 && m.width > 40 {
		maxWidth = 40
	}

	maxHeight := int(float64(m.height) * 0.8)
	if m.height < 25 {
		maxHeight = m.height - 1
	}
	if maxHeight < 10 && m.height > 10 {
		maxHeight = 10
	}

	centerY := true
	if maxHeight >= m.height-2 {
		centerY = false
	}

	out := components.RenderScreen(m.theme, components.ScreenData{
		Title:      m.view.Title,
		Subtitle:   m.view.Subtitle,
		Header:     header,
		Chart:      m.spark.View(),
		Flash:      m.snapshot.Flash,
		Lines:      m.view.Lines,
		Menu:       menu,
		Selected:   m.selectedIdx,
		Hint:       "Controls: j/down, k/up, enter, esc/backspace, q",
		WrapInCard: true,
		MaxWidth:   maxWidth,
		MaxHeight:  maxHeight,
		FullWidth:  m.width,
		FullHeight: m.height,
		CenterX:    true,
		CenterY:    centerY,
	})

	return zone.Scan(out)
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
	case scenes.ActionSetQuizMode:
		m.snapshot.QuizMode = action.Value
	case scenes.ActionToggleQuizConcept:
		if m.snapshot.QuizMode == "selected_concepts" {
			m.snapshot.QuizSelectedConcepts[action.Value] = !m.snapshot.QuizSelectedConcepts[action.Value]
		}
	case scenes.ActionExit:
		return true
	}

	m.render()
	return false
}

func (m *AppModel) render() {
	m.snapshot.CurrentLocale = m.runtime.Context.CurrentLocale
	m.snapshot.ActiveStateID = m.runtime.Context.ActiveStateID
	m.snapshot.EnabledMods = append([]string(nil), m.runtime.Context.EnabledMods...)
	m.snapshot.CurrentModID = m.runtime.Context.CurrentModID
	m.snapshot.CurrentUnitID = m.runtime.Context.CurrentUnitID
	m.snapshot.AvailableStates = listStateIDs(m.statesDir)
	m.view = scenes.Build(m.runtime.Context.CurrentSceneID, m.snapshot)
	if m.selectedIdx >= len(m.view.Options) {
		m.selectedIdx = 0
	}
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
