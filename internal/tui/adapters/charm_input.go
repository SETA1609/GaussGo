package adapters

import (
	tea "charm.land/bubbletea/v2"

	"gaussgo/internal/tui/ports"
)

type CharmInputMapper struct{}

func NewCharmInputMapper() ports.InputMapper {
	return CharmInputMapper{}
}

func (CharmInputMapper) Map(msg any) (ports.InputEvent, bool) {
	switch typed := msg.(type) {
	case tea.KeyMsg:
		return ports.InputEvent{Type: ports.InputEventKey, Key: ports.Key(typed.String())}, true
	case tea.WindowSizeMsg:
		return ports.InputEvent{Type: ports.InputEventWindowSize, Width: typed.Width, Height: typed.Height}, true
	default:
		return ports.InputEvent{}, false
	}
}
