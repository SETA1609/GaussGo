package controllers

import (
	"context"
	"strings"

	"gaussgo/internal/apperrors"
	"gaussgo/internal/contracts"
	"gaussgo/internal/mods"
	"gaussgo/internal/types"
)

var _ contracts.ModController = (*ModController)(nil)

type ModController struct {
	modsDir   string
	state     contracts.StateStore
	ctx       *types.RuntimeContext
	logger    contracts.LoggingService
	eventBus  contracts.EventBusService
	lifecycle *mods.RuntimeLifecycle
}

func NewModController(modsDir string, stateStore contracts.StateStore, ctx *types.RuntimeContext, logger contracts.LoggingService, eventBus contracts.EventBusService, lifecycle *mods.RuntimeLifecycle) *ModController {
	if ctx == nil {
		ctx = &types.RuntimeContext{}
	}
	return &ModController{modsDir: modsDir, state: stateStore, ctx: ctx, logger: logger, eventBus: eventBus, lifecycle: lifecycle}
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

	if c.lifecycle != nil {
		manifests, discoverErr := mods.Discover(c.modsDir)
		if discoverErr != nil {
			st.EnabledMods = removeValue(st.EnabledMods, modID)
			_ = c.state.Save(st)
			return discoverErr
		}

		enabled := enabledSet(st.EnabledMods)
		runtimeCopy := copyRuntimeContext(c.ctx, st.EnabledMods)
		if syncErr := c.lifecycle.SyncEnabled(context.Background(), manifests, runtimeCopy, enabled, nil); syncErr != nil {
			st.EnabledMods = removeValue(st.EnabledMods, modID)
			_ = c.state.Save(st)
			c.ctx.EnabledMods = append([]string(nil), st.EnabledMods...)
			emitEvent(c.eventBus, c.logger, EventModsToggled, "mod_controller", map[string]any{"modId": modID, "enabled": false, "reason": "runtime_init_failed"})
			if c.logger != nil {
				c.logger.Warn("mod enable rolled back after runtime init failure", map[string]any{"modId": modID, "error": syncErr.Error()})
			}
			return syncErr
		}
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

	if c.lifecycle != nil {
		if unloadErr := c.lifecycle.UnloadMod(context.Background(), modID); unloadErr != nil {
			if c.logger != nil {
				c.logger.Warn("runtime unload failed for disabled mod", map[string]any{"modId": modID, "error": unloadErr.Error()})
			}
		}
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

	if c.lifecycle != nil {
		manifests, discoverErr := mods.Discover(c.modsDir)
		if discoverErr != nil {
			return nil, discoverErr
		}
		runtimeCopy := copyRuntimeContext(c.ctx, st.EnabledMods)
		if syncErr := c.lifecycle.SyncEnabled(context.Background(), manifests, runtimeCopy, enabled, nil); syncErr != nil {
			if c.logger != nil {
				c.logger.Warn("runtime mod sync failed after refresh", map[string]any{"error": syncErr.Error()})
			}
			return nil, syncErr
		}
	}

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

func enabledSet(enabledMods []string) map[string]bool {
	out := make(map[string]bool, len(enabledMods))
	for _, modID := range enabledMods {
		out[modID] = true
	}
	return out
}

func removeValue(values []string, target string) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		if value != target {
			result = append(result, value)
		}
	}
	return result
}

func copyRuntimeContext(ctx *types.RuntimeContext, enabledMods []string) types.RuntimeContext {
	if ctx == nil {
		return types.RuntimeContext{EnabledMods: append([]string(nil), enabledMods...)}
	}
	return types.RuntimeContext{
		ActiveStateID:    ctx.ActiveStateID,
		EnabledMods:      append([]string(nil), enabledMods...),
		CurrentSceneID:   ctx.CurrentSceneID,
		CurrentModID:     ctx.CurrentModID,
		CurrentUnitID:    ctx.CurrentUnitID,
		CurrentConceptID: ctx.CurrentConceptID,
		CurrentLocale:    ctx.CurrentLocale,
	}
}
