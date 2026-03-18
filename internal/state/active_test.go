package state

import (
	"testing"
	"time"
)

func TestActivePointerReadWrite(t *testing.T) {
	dir := t.TempDir()

	if err := SetActive(dir, "default"); err != nil {
		t.Fatalf("set active failed: %v", err)
	}

	stateID, err := GetActive(dir)
	if err != nil {
		t.Fatalf("get active failed: %v", err)
	}
	if stateID != "default" {
		t.Fatalf("expected default, got %s", stateID)
	}
}

func TestStoreSetAndGetActive(t *testing.T) {
	store := NewStore(t.TempDir())
	if err := store.SetActive("abc"); err != nil {
		t.Fatalf("store set active failed: %v", err)
	}
	g, err := store.GetActive()
	if err != nil {
		t.Fatalf("store get active failed: %v", err)
	}
	if g != "abc" {
		t.Fatalf("expected abc, got %s", g)
	}
}

func TestGetActiveMissingPointer(t *testing.T) {
	if _, err := GetActive(t.TempDir()); err == nil {
		t.Fatal("expected missing active pointer error")
	}
}

func TestSetActiveEmptyStateID(t *testing.T) {
	if err := SetActive(t.TempDir(), "   "); err == nil {
		t.Fatal("expected empty stateID validation error")
	}
}

func TestSaveUpdatesTimestamp(t *testing.T) {
	store := NewStore(t.TempDir())
	state, err := store.Create("Time Test")
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	first, err := store.Load(state.StateID)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	time.Sleep(10 * time.Millisecond)
	if err := store.Save(first); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	second, err := store.Load(state.StateID)
	if err != nil {
		t.Fatalf("reload failed: %v", err)
	}
	if !second.UpdatedAt.After(first.UpdatedAt) {
		t.Fatalf("expected updatedAt to move forward: %v <= %v", second.UpdatedAt, first.UpdatedAt)
	}
}
