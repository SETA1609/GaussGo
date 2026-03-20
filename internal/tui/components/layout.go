package components

import (
	"strings"

	lipgloss "charm.land/lipgloss/v2"

	"gaussgo/internal/tui/ports"
)

type ScreenData struct {
	Title        string
	Subtitle     string
	Header       string
	Chart        string
	Flash        string
	Lines        []string
	Menu         []MenuItem
	Selected     int
	Hint         string
	WrapInCard   bool
	MaxWidth     int
	MaxHeight    int
	FullWidth    int
	FullHeight   int
	CenterX      bool
	CenterY      bool
	MenuRenderer ports.MenuRenderer
	LayoutEngine ports.LayoutEngine
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
	b.WriteString(RenderMenuWithRenderer(theme, data.Menu, data.Selected, data.MenuRenderer))

	if strings.TrimSpace(data.Hint) != "" {
		b.WriteString("\n")
		b.WriteString(theme.Hint.Render(data.Hint))
	}

	out := b.String()
	engine := data.LayoutEngine
	if engine == nil {
		engine = noOpLayoutEngine{}
	}
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
		hAlign := ports.AlignCenter
		if align == lipgloss.Left {
			hAlign = ports.AlignLeft
		}
		vAlignPort := ports.AlignCenter
		if vAlign == lipgloss.Top {
			vAlignPort = ports.AlignTop
		}
		out = engine.Place(width, height, hAlign, vAlignPort, out)
	}
	return out
}

type noOpLayoutEngine struct{}

func (noOpLayoutEngine) Place(_ int, _ int, _ ports.Align, _ ports.Align, content string) string {
	return content
}

func (noOpLayoutEngine) PlaceHorizontal(_ int, _ ports.Align, content string) string {
	return content
}

func (noOpLayoutEngine) JoinHorizontal(_ ports.Align, parts ...string) string {
	out := ""
	for _, part := range parts {
		out += part
	}
	return out
}

func (noOpLayoutEngine) JoinVertical(_ ports.Align, parts ...string) string {
	if len(parts) == 0 {
		return ""
	}
	out := parts[0]
	for i := 1; i < len(parts); i++ {
		out += "\n" + parts[i]
	}
	return out
}

func (noOpLayoutEngine) Block(content string, _ int, _ ports.Align) string {
	return content
}

func RenderMainLayout(theme Theme, navbar string, footer string, content string, sidepanel string, width int, height int, layout ports.LayoutEngine) string {
	if layout == nil {
		layout = noOpLayoutEngine{}
	}
	// Calculate middle section height
	navbarHeight := lipgloss.Height(navbar)
	footerHeight := lipgloss.Height(footer)
	middleHeight := height - navbarHeight - footerHeight
	if middleHeight < 0 {
		middleHeight = 0
	}

	// Calculate content and sidepanel widths
	sidepanelWidth := lipgloss.Width(sidepanel)
	contentWidth := width - sidepanelWidth
	if contentWidth < 0 {
		contentWidth = 0
	}

	// Ensure content and sidepanel have the correct height
	contentStyle := lipgloss.NewStyle().Height(middleHeight)
	sidepanelStyle := lipgloss.NewStyle().Height(middleHeight)

	middle := layout.JoinHorizontal(ports.AlignTop,
		contentStyle.Render(layout.Block(content, contentWidth, ports.AlignLeft)),
		sidepanelStyle.Render(layout.Block(sidepanel, sidepanelWidth, ports.AlignLeft)),
	)

	return layout.JoinVertical(ports.AlignLeft,
		navbar,
		middle,
		footer,
	)
}
