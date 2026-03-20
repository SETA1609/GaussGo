package components

import (
	"fmt"
	"strings"

	lipgloss "charm.land/lipgloss/v2"
)

type SidepanelData struct {
	AveragePercent float64
	Mods           []ModStats
	Width          int
	Height         int
	ScrollOffset   int
}

type ModStats struct {
	ID        string
	Percent   float64
	Expanded  bool
	Units     []UnitStats
}

type UnitStats struct {
	ID        string
	Percent   float64
	Expanded  bool
	Concepts  []ConceptStats
}

type ConceptStats struct {
	ID      string
	Percent float64
}

func RenderSidepanel(theme Theme, data SidepanelData) string {
	var lines []string

	// Header
	avgLine := theme.Title.Copy().
		Foreground(lipgloss.Color("86")).
		Render(fmt.Sprintf("AVG: %.1f%%", data.AveragePercent))
	lines = append(lines, avgLine, "")

	for _, mod := range data.Mods {
		prefix := "[+] "
		if mod.Expanded {
			prefix = "[-] "
		}
		
		modLine := fmt.Sprintf("%s%s (%.0f%%)", prefix, mod.ID, mod.Percent)
		lines = append(lines, theme.Header.Render(modLine))

		if mod.Expanded {
			for _, unit := range mod.Units {
				uPrefix := "  [+] "
				if unit.Expanded {
					uPrefix = "  [-] "
				}
				unitLine := fmt.Sprintf("%s%s (%.0f%%)", uPrefix, unit.ID, unit.Percent)
				lines = append(lines, theme.Subtitle.Render(unitLine))

				if unit.Expanded {
					for _, concept := range unit.Concepts {
						conceptLine := fmt.Sprintf("    %s: %.0f%%", concept.ID, concept.Percent)
						lines = append(lines, theme.Line.Render(conceptLine))
					}
				}
			}
		}
	}

	// Simple scroll implementation
	start := data.ScrollOffset
	if start < 0 {
		start = 0
	}
	if start >= len(lines) {
		start = len(lines) - 1
	}
	if start < 0 {
		start = 0
	}
	
	end := start + data.Height
	if end > len(lines) {
		end = len(lines)
	}

	visibleLines := lines[start:end]
	
	panelStyle := theme.Card.Copy().
		Width(data.Width).
		Height(data.Height).
		BorderForeground(lipgloss.Color("240"))

	return panelStyle.Render(strings.Join(visibleLines, "\n"))
}
