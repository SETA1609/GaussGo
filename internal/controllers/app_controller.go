package controllers

import (
	"strings"

	"gaussgo/internal/contracts"
	"gaussgo/internal/types"
)

var _ contracts.AppController = (*AppController)(nil)

type AppController struct {
	stateStore contracts.StateStore
	ctx        *types.RuntimeContext
	logger     contracts.LoggingService
	eventBus   contracts.EventBusService
}

func NewAppController(stateStore contracts.StateStore, ctx *types.RuntimeContext, logger contracts.LoggingService, eventBus contracts.EventBusService) *AppController {
	if ctx == nil {
		ctx = &types.RuntimeContext{}
	}
	return &AppController{stateStore: stateStore, ctx: ctx, logger: logger, eventBus: eventBus}
}

func (c *AppController) LoadActiveState() (contracts.State, error) {
	activeID := strings.TrimSpace(c.ctx.ActiveStateID)
	if activeID == "" {
		id, err := c.stateStore.GetActive()
		if err != nil {
			return contracts.State{}, err
		}
		activeID = id
	}

	st, err := c.stateStore.Load(activeID)
	if err != nil {
		return contracts.State{}, err
	}

	c.ctx.ActiveStateID = st.StateID
	c.ctx.EnabledMods = append([]string(nil), st.EnabledMods...)
	c.ctx.CurrentLocale = st.UI.Locale
	if c.ctx.CurrentLocale == "" {
		c.ctx.CurrentLocale = "en"
	}

	emitEvent(c.eventBus, c.logger, EventStateLoaded, "app_controller", map[string]any{
		"stateId":    st.StateID,
		"locale":     c.ctx.CurrentLocale,
		"normalized": true,
	})
	if c.logger != nil {
		c.logger.Info("state loaded", map[string]any{"stateId": st.StateID, "locale": c.ctx.CurrentLocale})
	}

	return st, nil
}

func (c *AppController) CreateAndActivateState(profileName string) (contracts.State, error) {
	st, err := c.stateStore.Create(profileName)
	if err != nil {
		return contracts.State{}, err
	}
	if err := c.stateStore.SetActive(st.StateID); err != nil {
		return contracts.State{}, err
	}
	emitEvent(c.eventBus, c.logger, EventStateCreated, "app_controller", map[string]any{
		"stateId": st.StateID,
		"profile": profileName,
	})
	if c.logger != nil {
		c.logger.Info("state created and activated", map[string]any{"stateId": st.StateID})
	}
	c.ctx.ActiveStateID = "" // let LoadActiveState pick up the pointer
	return c.LoadActiveState()
}

func (c *AppController) ActivateState(stateID string) error {
	if err := c.stateStore.SetActive(stateID); err != nil {
		return err
	}
	c.ctx.ActiveStateID = "" // force re-read
	_, err := c.LoadActiveState()
	return err
}
