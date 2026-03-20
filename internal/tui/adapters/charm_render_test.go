package adapters

import (
	"testing"

	"gaussgo/internal/tui/ports"
)

func TestCharmRendererReturnsBubbleTeaView(t *testing.T) {
	renderer := NewCharmRenderer()
	rendered := renderer.Render("hello", ports.ViewConfig{AltScreen: true, MouseMode: MouseModeCellMotion})
	if rendered == nil {
		t.Fatal("expected non-nil rendered view")
	}
}
