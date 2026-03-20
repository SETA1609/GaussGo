package adapters

import (
	bubblespinner "charm.land/bubbles/v2/spinner"

	"gaussgo/internal/tui/ports"
)

type CharmSpinner struct {
	spinner bubblespinner.Model
}

func NewCharmSpinner() ports.Spinner {
	return &CharmSpinner{spinner: bubblespinner.New()}
}

func (s *CharmSpinner) InitCmd() any {
	return s.spinner.Tick
}

func (s *CharmSpinner) Update(msg any) (cmd any, handled bool) {
	tickMsg, ok := msg.(bubblespinner.TickMsg)
	if !ok {
		return nil, false
	}
	var teaCmd any
	s.spinner, teaCmd = s.spinner.Update(tickMsg)
	return teaCmd, true
}

func (s *CharmSpinner) View() string {
	return s.spinner.View()
}
