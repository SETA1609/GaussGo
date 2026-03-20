package components

import (
	"gaussgo/internal/tui/ports"
)

func RenderNavbar(theme Theme, width int, layout ports.LayoutEngine) string {
	if layout == nil {
		layout = noOpLayoutEngine{}
	}
	title := theme.Title.
		Copy().
		Bold(true).
		Padding(0, 2).
		Render("GAUSS GO")

	return layout.PlaceHorizontal(width, ports.AlignCenter, title)
}
