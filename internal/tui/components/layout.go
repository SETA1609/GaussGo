package components

import (
	"strings"

	lipgloss "charm.land/lipgloss/v2"
)

type ScreenData struct {
	Title      string
	Subtitle   string
	Header     string
	Chart      string
	Flash      string
	Lines      []string
	Menu       []MenuItem
	Selected   int
	Hint       string
	WrapInCard bool
	MaxWidth   int
	MaxHeight  int
	FullWidth  int
	FullHeight int
	CenterX    bool
	CenterY    bool
}

func RenderScreen(theme Theme, data ScreenData) string {
	var b strings.Builder

	b.WriteString(theme.Title.Render(data.Title))
	b.WriteString("\n")
	if strings.TrimSpace(data.Subtitle) != "" {
		b.WriteString(theme.Subtitle.Render(data.Subtitle))
		b.WriteString("\n")
	}
	if strings.TrimSpace(data.Header) != "" {
		b.WriteString(theme.Header.Render(data.Header))
		b.WriteString("\n")
	}
	if strings.TrimSpace(data.Chart) != "" {
		b.WriteString(theme.Subtitle.Render("activity"))
		b.WriteString("\n")
		b.WriteString(data.Chart)
		b.WriteString("\n")
	}
	if strings.TrimSpace(data.Flash) != "" {
		b.WriteString(theme.Flash.Render(data.Flash))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	for _, line := range data.Lines {
		b.WriteString(theme.Line.Render("- " + line))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(RenderMenu(theme, data.Menu, data.Selected))

	if strings.TrimSpace(data.Hint) != "" {
		b.WriteString("\n")
		b.WriteString(theme.Hint.Render(data.Hint))
	}

	out := b.String()
	if data.WrapInCard {
		card := theme.Card
		if data.MaxWidth > 0 {
			card = card.MaxWidth(data.MaxWidth)
		}
		if data.MaxHeight > 0 {
			card = card.MaxHeight(data.MaxHeight)
		}
		out = card.Render(out)
	}

	if data.CenterX || data.CenterY {
		align := lipgloss.Center
		if !data.CenterX {
			align = lipgloss.Left
		}
		vAlign := lipgloss.Center
		if !data.CenterY {
			vAlign = lipgloss.Top
		}
		width := data.FullWidth
		height := data.FullHeight
		if width <= 0 {
			width = data.MaxWidth
		}
		if width <= 0 {
			width = lipgloss.Width(out)
		}
		if height <= 0 {
			height = data.MaxHeight
		}
		if height <= 0 {
			height = lipgloss.Height(out)
		}
		out = lipgloss.Place(width, height, align, vAlign, out)
	}
	return out
}
