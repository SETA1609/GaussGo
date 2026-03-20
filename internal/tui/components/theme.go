package components

import lipgloss "charm.land/lipgloss/v2"

type Theme struct {
	Title      lipgloss.Style
	Subtitle   lipgloss.Style
	Header     lipgloss.Style
	Line       lipgloss.Style
	Cursor     lipgloss.Style
	Selected   lipgloss.Style
	Disabled   lipgloss.Style
	Flash      lipgloss.Style
	Hint       lipgloss.Style
	Card       lipgloss.Style
	StatusPill lipgloss.Style
}

func DefaultTheme() Theme {
	return Theme{
		Title:      lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")),
		Subtitle:   lipgloss.NewStyle().Foreground(lipgloss.Color("241")),
		Header:     lipgloss.NewStyle().Foreground(lipgloss.Color("111")),
		Line:       lipgloss.NewStyle().Foreground(lipgloss.Color("252")),
		Cursor:     lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86")),
		Selected:   lipgloss.NewStyle().Foreground(lipgloss.Color("230")).Background(lipgloss.Color("60")).Padding(0, 1),
		Disabled:   lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
		Flash:      lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Italic(true),
		Hint:       lipgloss.NewStyle().Foreground(lipgloss.Color("244")),
		Card:       lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("62")).Padding(0, 1),
		StatusPill: lipgloss.NewStyle().Foreground(lipgloss.Color("230")).Background(lipgloss.Color("61")).Padding(0, 1),
	}
}
