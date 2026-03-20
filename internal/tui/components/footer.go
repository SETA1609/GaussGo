package components

import (
	"gaussgo/internal/tui/ports"
)

type FooterData struct {
	Hint   string
	User   string
	State  string
	Locale string
	Width  int
}

func RenderFooter(theme Theme, data FooterData, layout ports.LayoutEngine) string {
	if layout == nil {
		layout = noOpLayoutEngine{}
	}
	hintStyle := theme.Hint.Copy().Italic(true)
	statusStyle := theme.StatusPill.Copy()

	hint := hintStyle.Render(data.Hint)

	statusInfos := []string{
		"User: " + data.User,
		"State: " + data.State,
		"Locale: " + data.Locale,
	}

	var statusPart string
	for _, info := range statusInfos {
		statusPart += statusStyle.Render(info) + " "
	}

	left := layout.Block(hint, data.Width/2, ports.AlignLeft)
	right := layout.Block(statusPart, data.Width/2, ports.AlignRight)

	return layout.JoinHorizontal(ports.AlignBottom, left, right)
}
