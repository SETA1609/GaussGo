package adapters

import (
	"testing"

	bubblespinner "charm.land/bubbles/v2/spinner"
)

func TestCharmSpinnerInitAndUpdate(t *testing.T) {
	spinner := NewCharmSpinner()
	if spinner.InitCmd() == nil {
		t.Fatal("expected init command from spinner")
	}

	if _, handled := spinner.Update(struct{}{}); handled {
		t.Fatal("expected unknown message to be unhandled")
	}

	if _, handled := spinner.Update(bubblespinner.TickMsg{}); !handled {
		t.Fatal("expected tick message to be handled")
	}

	_ = spinner.View()
}
