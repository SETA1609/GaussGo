package adapters

import "testing"

func TestCharmZoneMarkAndScan(t *testing.T) {
	z := NewCharmZone()
	z.Init()

	marked := z.Mark("id-1", "hello")
	if marked == "" {
		t.Fatal("expected marked content")
	}

	scanned := z.Scan(marked)
	if scanned == "" {
		t.Fatal("expected scanned content")
	}
}
