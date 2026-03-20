package components

import (
	"strings"

	"gaussgo/internal/tui/ports"
)

type MenuItem struct {
	ID       string
	Label    string
	Disabled bool
}

type defaultMenuRenderer struct{}

func (defaultMenuRenderer) Mark(_ string, content string) string {
	return content
}

func RenderMenu(theme Theme, items []MenuItem, selected int) string {
	return RenderMenuWithRenderer(theme, items, selected, defaultMenuRenderer{})
}

func RenderMenuWithRenderer(theme Theme, items []MenuItem, selected int, renderer ports.MenuRenderer) string {
	if len(items) == 0 {
		return ""
	}
	if renderer == nil {
		renderer = defaultMenuRenderer{}
	}

	var b strings.Builder
	for i, item := range items {
		cursor := "  "
		if i == selected {
			cursor = theme.Cursor.Render("> ")
		}

		label := item.Label
		if item.Disabled {
			label += " (locked)"
			label = theme.Disabled.Render(label)
		} else if i == selected {
			label = theme.Selected.Render(label)
		}

		zoneID := item.ID
		if strings.TrimSpace(zoneID) == "" {
			zoneID = "item-" + normalizeID(item.Label)
		}

		b.WriteString(cursor)
		b.WriteString(renderer.Mark(zoneID, label))
		b.WriteString("\n")
	}

	return b.String()
}

func normalizeID(in string) string {
	in = strings.ToLower(strings.TrimSpace(in))
	in = strings.ReplaceAll(in, " ", "-")
	in = strings.ReplaceAll(in, "/", "-")
	in = strings.ReplaceAll(in, ".", "-")
	if in == "" {
		return "menu-item"
	}
	return in
}
