package adapters

import (
	lipgloss "charm.land/lipgloss/v2"

	"gaussgo/internal/tui/ports"
)

type CharmLayoutEngine struct{}

func NewCharmLayoutEngine() ports.LayoutEngine {
	return CharmLayoutEngine{}
}

func (CharmLayoutEngine) Place(width int, height int, hAlign ports.Align, vAlign ports.Align, content string) string {
	return lipgloss.Place(width, height, mapAlign(hAlign), mapAlign(vAlign), content)
}

func (CharmLayoutEngine) PlaceHorizontal(width int, hAlign ports.Align, content string) string {
	return lipgloss.PlaceHorizontal(width, mapAlign(hAlign), content)
}

func (CharmLayoutEngine) JoinHorizontal(vAlign ports.Align, parts ...string) string {
	return lipgloss.JoinHorizontal(mapAlign(vAlign), parts...)
}

func (CharmLayoutEngine) JoinVertical(hAlign ports.Align, parts ...string) string {
	return lipgloss.JoinVertical(mapAlign(hAlign), parts...)
}

func (CharmLayoutEngine) Block(content string, width int, hAlign ports.Align) string {
	return lipgloss.NewStyle().Width(width).Align(mapAlign(hAlign)).Render(content)
}

func mapAlign(align ports.Align) lipgloss.Position {
	switch align {
	case ports.AlignLeft:
		return lipgloss.Left
	case ports.AlignCenter:
		return lipgloss.Center
	case ports.AlignRight:
		return lipgloss.Right
	case ports.AlignTop:
		return lipgloss.Top
	case ports.AlignBottom:
		return lipgloss.Bottom
	default:
		return lipgloss.Left
	}
}
