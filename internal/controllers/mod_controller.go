package controllers

import (
	"strings"

	"gaussgo/internal/apperrors"
	"gaussgo/internal/contracts"
	"gaussgo/internal/mods"
	"gaussgo/internal/types"
)

var _ contracts.ModController = (*ModController)(nil)

type ModController struct {
	modsDir  string
	state    contracts.StateStore
	ctx      *types.RuntimeContext
	logger   contracts.LoggingService
	eventBus contracts.EventBusService
}

func NewModController(modsDir string, stateStore contracts.StateStore, ctx *types.RuntimeContext, logger contracts.LoggingService, eventBus contracts.EventBusService) *ModController {
	if ctx == nil {
		ctx = &types.RuntimeContext{}
	}
	return &ModController{modsDir: modsDir, state: stateStore, ctx: ctx, logger: logger, eventBus: eventBus}
}

func (c *ModController) Enable(modID string) error {
	modID = strings.TrimSpace(modID)
	if modID == "" {
		return apperrors.New(apperrors.CodeValidation, apperrors.ErrorTypeInput, "modID is required")
	}

	exists, err := c.modExists(modID)
	if err != nil {
		return err
	}
	if !exists {
		return apperrors.New(apperrors.CodeNotFound, apperrors.ErrorTypeInput, "mod not found: "+modID)
	}

	st, err := c.loadActiveState()
	if err != nil {
		return err
	}

	if !contains(st.EnabledMods, modID) {
		st.EnabledMods = append(st.EnabledMods, modID)
	}
	if err := c.state.Save(st); err != nil {
		return err
	}

	c.ctx.EnabledMods = append([]string(nil), st.EnabledMods...)
	emitEvent(c.eventBus, c.logger, EventModsToggled, "mod_controller", map[string]any{"modId": modID, "enabled": true})
	if c.logger != nil {
		c.logger.Info("mod enabled", map[string]any{"modId": modID})
	}

	return nil
}

func (c *ModController) Disable(modID string) error {
	modID = strings.TrimSpace(modID)
	if modID == "" {
		return apperrors.New(apperrors.CodeValidation, apperrors.ErrorTypeInput, "modID is required")
	}
	if modID == "core" {
		return apperrors.New(apperrors.CodeCoreRequired, apperrors.ErrorTypeDomain, "cannot disable core mod")
	}

	st, err := c.loadActiveState()
	if err != nil {
		return err
	}

	next := make([]string, 0, len(st.EnabledMods))
	for _, id := range st.EnabledMods {
		if id != modID {
			next = append(next, id)
		}
	}
	st.EnabledMods = next

	if err := c.state.Save(st); err != nil {
		return err
	}

	c.ctx.EnabledMods = append([]string(nil), st.EnabledMods...)
	if c.ctx.CurrentModID == modID {
		c.ctx.CurrentModID = ""
		c.ctx.CurrentUnitID = ""
		c.ctx.CurrentConceptID = ""
	}
	emitEvent(c.eventBus, c.logger, EventModsToggled, "mod_controller", map[string]any{"modId": modID, "enabled": false})
	if c.logger != nil {
		c.logger.Info("mod disabled", map[string]any{"modId": modID})
	}

	return nil
}

func (c *ModController) Refresh() ([]contracts.ModStatus, error) {
	st, err := c.loadActiveState()
	if err != nil {
		return nil, err
	}

	enabled := make(map[string]bool, len(st.EnabledMods))
	for _, id := range st.EnabledMods {
		enabled[id] = true
	}

	repo := mods.NewRepository(c.modsDir, enabled)
	statuses, err := repo.Refresh()
	if err != nil {
		if c.logger != nil {
			c.logger.Error("mod refresh failed", map[string]any{"error": err.Error()})
		}
		return nil, err
	}

	c.ctx.EnabledMods = enabledFromStatuses(statuses)
	validCount, invalidCount := summarizeRefreshCounts(statuses)
	emitEvent(c.eventBus, c.logger, EventModsRefreshed, "mod_controller", map[string]any{
		"total":   len(statuses),
		"valid":   validCount,
		"invalid": invalidCount,
	})
	if c.logger != nil {
		c.logger.Info("mods refreshed", map[string]any{"count": len(statuses)})
	}

	return statuses, nil
}

func (c *ModController) loadActiveState() (contracts.State, error) {
	activeID := strings.TrimSpace(c.ctx.ActiveStateID)
	if activeID == "" {
		id, err := c.state.GetActive()
		if err != nil {
			return contracts.State{}, err
		}
		activeID = id
		c.ctx.ActiveStateID = id
	}

	return c.state.Load(activeID)
}

func enabledFromStatuses(statuses []contracts.ModStatus) []string {
	out := make([]string, 0, len(statuses))
	for _, s := range statuses {
		if s.Enabled {
			out = append(out, s.ModID)
		}
	}
	return out
}

func contains(in []string, target string) bool {
	for _, v := range in {
		if v == target {
			return true
		}
	}
	return false
}

func (c *ModController) modExists(modID string) (bool, error) {
	manifests, err := mods.Discover(c.modsDir)
	if err != nil {
		return false, err
	}
	for _, manifest := range manifests {
		if manifest.ID == modID {
			return true, nil
		}
	}
	return false, nil
}

func summarizeRefreshCounts(statuses []contracts.ModStatus) (valid int, invalid int) {
	for _, status := range statuses {
		if status.Available {
			valid++
		} else {
			invalid++
		}
	}
	return valid, invalid
}
