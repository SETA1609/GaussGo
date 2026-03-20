package controllers

import (
	"testing"

	"gaussgo/internal/types"
)

func TestLearningControllerSelectionFlow(t *testing.T) {
	ctx := &types.RuntimeContext{}
	c := NewLearningController(ctx)

	if err := c.SelectMod("linearAlgebra"); err != nil {
		t.Fatalf("select mod failed: %v", err)
	}
	if err := c.SelectUnit("vectors"); err != nil {
		t.Fatalf("select unit failed: %v", err)
	}
	if err := c.SelectConcept("dot-product"); err != nil {
		t.Fatalf("select concept failed: %v", err)
	}

	if ctx.CurrentModID != "linearAlgebra" {
		t.Fatalf("unexpected mod id: %s", ctx.CurrentModID)
	}
	if ctx.CurrentUnitID != "vectors" {
		t.Fatalf("unexpected unit id: %s", ctx.CurrentUnitID)
	}
	if ctx.CurrentConceptID != "vectors.dot-product" {
		t.Fatalf("unexpected concept id: %s", ctx.CurrentConceptID)
	}
}

func TestLearningControllerUnitRequiresMod(t *testing.T) {
	c := NewLearningController(&types.RuntimeContext{})
	if err := c.SelectUnit("vectors"); err == nil {
		t.Fatal("expected error selecting unit before mod")
	}
}

func TestLearningControllerRejectsEmptyInputs(t *testing.T) {
	ctx := &types.RuntimeContext{}
	c := NewLearningController(ctx)

	if err := c.SelectMod("  "); err == nil {
		t.Fatal("expected select mod validation error")
	}
	if err := c.SelectUnit(""); err == nil {
		t.Fatal("expected select unit validation error")
	}
	if err := c.SelectConcept("  "); err == nil {
		t.Fatal("expected select concept validation error")
	}
}

func TestLearningControllerSelectConceptRequiresUnit(t *testing.T) {
	ctx := &types.RuntimeContext{}
	c := NewLearningController(ctx)
	if err := c.SelectMod("linearAlgebra"); err != nil {
		t.Fatalf("select mod failed: %v", err)
	}

	if err := c.SelectConcept("dot-product"); err == nil {
		t.Fatal("expected dependency error selecting concept before unit")
	}
}
