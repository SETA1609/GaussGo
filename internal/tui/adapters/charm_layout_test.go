package adapters

import (
	"testing"

	"gaussgo/internal/tui/ports"
)

func TestCharmLayoutEngineMethods(t *testing.T) {
	engine := NewCharmLayoutEngine()
	if got := engine.PlaceHorizontal(20, ports.AlignCenter, "x"); got == "" {
		t.Fatal("expected non-empty PlaceHorizontal output")
	}
	if got := engine.JoinHorizontal(ports.AlignTop, "a", "b"); got == "" {
		t.Fatal("expected non-empty JoinHorizontal output")
	}
	if got := engine.JoinVertical(ports.AlignLeft, "a", "b"); got == "" {
		t.Fatal("expected non-empty JoinVertical output")
	}
	if got := engine.Block("abc", 10, ports.AlignRight); got == "" {
		t.Fatal("expected non-empty Block output")
	}
	if got := engine.Place(20, 5, ports.AlignCenter, ports.AlignTop, "x"); got == "" {
		t.Fatal("expected non-empty Place output")
	}
}
