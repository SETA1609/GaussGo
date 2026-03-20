package controllers

import (
	"sort"
	"strings"

	"gaussgo/internal/apperrors"
	"gaussgo/internal/contracts"
	"gaussgo/internal/state"
	"gaussgo/internal/types"
)

var _ contracts.LocaleController = (*LocaleController)(nil)

type LocaleController struct {
	stateStore       contracts.StateStore
	ctx              *types.RuntimeContext
	logger           contracts.LoggingService
	eventBus         contracts.EventBusService
	supportedLocales map[string]struct{}
}

func NewLocaleController(stateStore contracts.StateStore, ctx *types.RuntimeContext, logger contracts.LoggingService, eventBus contracts.EventBusService, supportedLocales map[string]struct{}) *LocaleController {
	if ctx == nil {
		ctx = &types.RuntimeContext{}
	}
	if len(supportedLocales) == 0 {
		supportedLocales = state.SupportedLocalesSet()
	}
	return &LocaleController{stateStore: stateStore, ctx: ctx, logger: logger, eventBus: eventBus, supportedLocales: cloneLocaleSet(supportedLocales)}
}

func (c *LocaleController) Current() string {
	locale := strings.TrimSpace(c.ctx.CurrentLocale)
	if locale == "" {
		return state.DefaultLocale
	}
	return locale
}

func (c *LocaleController) SetLocale(locale string) error {
	locale = strings.TrimSpace(locale)
	if locale == "" {
		return apperrors.New(apperrors.CodeValidation, apperrors.ErrorTypeInput, "locale is required")
	}
	if _, ok := c.supportedLocales[locale]; !ok {
		return apperrors.New(apperrors.CodeUnsupportedLocale, apperrors.ErrorTypeInput, "unsupported locale: "+locale)
	}

	activeID := strings.TrimSpace(c.ctx.ActiveStateID)
	if activeID == "" {
		id, err := c.stateStore.GetActive()
		if err != nil {
			return err
		}
		activeID = id
		c.ctx.ActiveStateID = id
	}

	st, err := c.stateStore.Load(activeID)
	if err != nil {
		return err
	}

	previous := st.UI.Locale
	st.UI.Locale = locale
	if err := c.stateStore.Save(st); err != nil {
		return err
	}

	c.ctx.CurrentLocale = locale
	emitEvent(c.eventBus, c.logger, EventLocaleChanged, "locale_controller", map[string]any{
		"stateId":  st.StateID,
		"previous": previous,
		"current":  locale,
	})
	if c.logger != nil {
		c.logger.Info("locale changed", map[string]any{"stateId": st.StateID, "previous": previous, "current": locale})
	}

	return nil
}

func SupportedLocales() []string {
	set := state.SupportedLocalesSet()
	locales := make([]string, 0, len(set))
	for locale := range set {
		locales = append(locales, locale)
	}
	sort.Strings(locales)
	return locales
}

func cloneLocaleSet(in map[string]struct{}) map[string]struct{} {
	out := make(map[string]struct{}, len(in))
	for locale := range in {
		out[locale] = struct{}{}
	}
	return out
}
