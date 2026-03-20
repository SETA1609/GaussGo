package adapters

import lipgloss "charm.land/lipgloss/v2"

type CharmLipGlossMetrics struct{}

func NewCharmLipGlossMetrics() CharmLipGlossMetrics {
	return CharmLipGlossMetrics{}
}

func (CharmLipGlossMetrics) Height(content string) int {
	return lipgloss.Height(content)
}

func (CharmLipGlossMetrics) Width(content string) int {
	return lipgloss.Width(content)
}
