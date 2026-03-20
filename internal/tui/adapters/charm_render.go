package adapters

import (
	tea "charm.land/bubbletea/v2"

	"gaussgo/internal/tui/ports"
)

const MouseModeCellMotion = "cell_motion"

type CharmRenderer struct{}

func NewCharmRenderer() ports.ViewPort {
	return CharmRenderer{}
}

func (CharmRenderer) Render(text string, config ports.ViewConfig) any {
	v := tea.NewView(text)
	v.AltScreen = config.AltScreen
	if config.MouseMode == MouseModeCellMotion {
		v.MouseMode = tea.MouseModeCellMotion
	}
	return v
}
