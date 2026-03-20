package controllers

import (
	"strings"

	"gaussgo/internal/apperrors"
	"gaussgo/internal/contracts"
	"gaussgo/internal/types"
)

var _ contracts.LearningController = (*LearningController)(nil)

type LearningController struct {
	ctx *types.RuntimeContext
}

func NewLearningController(ctx *types.RuntimeContext) *LearningController {
	if ctx == nil {
		ctx = &types.RuntimeContext{}
	}
	return &LearningController{ctx: ctx}
}

func (c *LearningController) SelectMod(modID string) error {
	modID = strings.TrimSpace(modID)
	if modID == "" {
		return apperrors.New(apperrors.CodeValidation, apperrors.ErrorTypeInput, "modID is required")
	}
	c.ctx.CurrentModID = modID
	c.ctx.CurrentUnitID = ""
	c.ctx.CurrentConceptID = ""
	return nil
}

func (c *LearningController) SelectUnit(unitID string) error {
	unitID = strings.TrimSpace(unitID)
	if unitID == "" {
		return apperrors.New(apperrors.CodeValidation, apperrors.ErrorTypeInput, "unitID is required")
	}
	if strings.TrimSpace(c.ctx.CurrentModID) == "" {
		return apperrors.New(apperrors.CodeDependency, apperrors.ErrorTypeDomain, "cannot select unit without selected mod")
	}
	c.ctx.CurrentUnitID = unitID
	c.ctx.CurrentConceptID = ""
	return nil
}

func (c *LearningController) SelectConcept(conceptID string) error {
	conceptID = strings.TrimSpace(conceptID)
	if conceptID == "" {
		return apperrors.New(apperrors.CodeValidation, apperrors.ErrorTypeInput, "conceptID is required")
	}
	if strings.TrimSpace(c.ctx.CurrentModID) == "" || strings.TrimSpace(c.ctx.CurrentUnitID) == "" {
		return apperrors.New(apperrors.CodeDependency, apperrors.ErrorTypeDomain, "cannot select concept without selected mod and unit")
	}
	c.ctx.CurrentConceptID = types.JoinNamespacedID(c.ctx.CurrentUnitID, conceptID)
	return nil
}
