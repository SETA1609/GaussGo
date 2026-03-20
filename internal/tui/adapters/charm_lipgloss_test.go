package adapters

import "testing"

func TestCharmLipGlossMetrics(t *testing.T) {
	m := NewCharmLipGlossMetrics()
	if got := m.Width("abc"); got <= 0 {
		t.Fatalf("expected positive width, got %d", got)
	}
	if got := m.Height("a\nb"); got != 2 {
		t.Fatalf("expected height 2, got %d", got)
	}
}
