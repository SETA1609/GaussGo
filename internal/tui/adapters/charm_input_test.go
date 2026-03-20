package adapters

import (
	"testing"

	tea "charm.land/bubbletea/v2"

	"gaussgo/internal/tui/ports"
)

func TestCharmInputMapperMapsKeyAndWindowSize(t *testing.T) {
	mapper := NewCharmInputMapper()

	keyEvent, ok := mapper.Map(tea.KeyPressMsg(tea.Key{Text: "q", Code: 'q'}))
	if !ok {
		t.Fatal("expected key event to map")
	}
	if keyEvent.Type != ports.InputEventKey || keyEvent.Key != ports.Key("q") {
		t.Fatalf("unexpected key event: %#v", keyEvent)
	}

	sizeEvent, ok := mapper.Map(tea.WindowSizeMsg{Width: 120, Height: 40})
	if !ok {
		t.Fatal("expected window size event to map")
	}
	if sizeEvent.Type != ports.InputEventWindowSize || sizeEvent.Width != 120 || sizeEvent.Height != 40 {
		t.Fatalf("unexpected window event: %#v", sizeEvent)
	}
}
