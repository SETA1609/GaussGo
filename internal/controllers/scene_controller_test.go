package controllers

import (
	"testing"

	"gaussgo/internal/types"
)

func TestSceneControllerNavigateBack(t *testing.T) {
	c := NewSceneController("main_menu")
	c.Navigate("mod_select")
	c.Navigate("learning_hub")

	if got := c.Current(); got != "learning_hub" {
		t.Fatalf("expected current learning_hub, got %s", got)
	}

	if err := c.Back(); err != nil {
		t.Fatalf("back failed: %v", err)
	}
	if got := c.Current(); got != "mod_select" {
		t.Fatalf("expected current mod_select, got %s", got)
	}
}

func TestSceneControllerBackEmptyHistory(t *testing.T) {
	c := NewSceneController("main_menu")
	if err := c.Back(); err == nil {
		t.Fatal("expected error when history is empty")
	}
}

func TestSceneControllerWithContextSyncsOnNavigateAndBack(t *testing.T) {
	ctx := &types.RuntimeContext{}
	c := NewSceneControllerWithContext("main_menu", ctx)
	if ctx.CurrentSceneID != "main_menu" {
		t.Fatalf("expected context scene main_menu, got %s", ctx.CurrentSceneID)
	}

	c.Navigate("mod_select")
	if ctx.CurrentSceneID != "mod_select" {
		t.Fatalf("expected context scene mod_select, got %s", ctx.CurrentSceneID)
	}

	if err := c.Back(); err != nil {
		t.Fatalf("back failed: %v", err)
	}
	if ctx.CurrentSceneID != "main_menu" {
		t.Fatalf("expected context scene main_menu, got %s", ctx.CurrentSceneID)
	}
}

func TestSceneControllerWithContextBackDrainsHistory(t *testing.T) {
	ctx := &types.RuntimeContext{}
	c := NewSceneControllerWithContext("main_menu", ctx)
	c.Navigate("mod_select")
	c.Navigate("unit_select")
	c.Navigate("learning_hub")

	if err := c.Back(); err != nil {
		t.Fatalf("back 1 failed: %v", err)
	}
	if ctx.CurrentSceneID != "unit_select" {
		t.Fatalf("expected context scene unit_select, got %s", ctx.CurrentSceneID)
	}

	if err := c.Back(); err != nil {
		t.Fatalf("back 2 failed: %v", err)
	}
	if ctx.CurrentSceneID != "mod_select" {
		t.Fatalf("expected context scene mod_select, got %s", ctx.CurrentSceneID)
	}

	if err := c.Back(); err != nil {
		t.Fatalf("back 3 failed: %v", err)
	}
	if ctx.CurrentSceneID != "main_menu" {
		t.Fatalf("expected context scene main_menu, got %s", ctx.CurrentSceneID)
	}

	if err := c.Back(); err == nil {
		t.Fatal("expected error after draining history")
	}
}
