package bootstrap

import (
	"context"
	"fmt"
	"gaussgo/internal/apperrors"
	"gaussgo/internal/concurrency"
	"gaussgo/internal/contracts"
	"gaussgo/internal/controllers"
	"gaussgo/internal/mods"
	"gaussgo/internal/state"
	"gaussgo/internal/types"
)

type Diagnostics struct {
	DiscoveredMods []string
	LoadOrder      []string
	Warnings       []string
}

type RuntimeServices struct {
	Logger      contracts.LoggingService
	EventBus    contracts.EventBusService
	Concurrency contracts.ConcurrencyService
}

type RuntimeControllers struct {
	App      *controllers.AppController
	Scene    contracts.SceneController
	Learning contracts.LearningController
	Mod      contracts.ModController
	Locale   contracts.LocaleController
}

type Runtime struct {
	Diagnostics Diagnostics
	Context     *types.RuntimeContext
	Controllers RuntimeControllers
	ModRegistry contracts.RuntimeModRegistry
	ModRuntime  *mods.RuntimeLifecycle
}

var newStoreForBootstrapRuntime = func(statesDir string) contracts.StateStore {
	return state.NewStore(statesDir)
}

func DefaultServices() RuntimeServices {
	return RuntimeServices{
		Logger:      resolveDefaultLogger(),
		EventBus:    NewInMemoryEventBus(),
		Concurrency: concurrency.DefaultService(),
	}
}

func Bootstrap(modsDir string, services RuntimeServices) (Diagnostics, error) {
	services = normalizeServices(services)

	// Bootstrap accepts optional injected services so callers can provide
	// production implementations while tests use defaults or stubs.
	services.Logger.Info("bootstrap started", map[string]any{"modsDir": modsDir})

	manifests, err := mods.Discover(modsDir)
	if err != nil {
		services.Logger.Error("bootstrap failed discovering mods", map[string]any{"error": err.Error()})
		return Diagnostics{}, apperrors.Wrap(apperrors.CodeValidation, apperrors.ErrorTypeDomain, "discover mods", err)
	}

	loadOrder, err := mods.ResolveLoadOrder(manifests)
	if err != nil {
		services.Logger.Error("bootstrap failed resolving load order", map[string]any{"error": err.Error()})
		return Diagnostics{}, apperrors.Wrap(apperrors.CodeDependency, apperrors.ErrorTypeDomain, "resolve load order", err)
	}

	discovered := make([]string, 0, len(manifests))
	for _, m := range manifests {
		discovered = append(discovered, m.ID)
	}

	diagnostics := Diagnostics{
		DiscoveredMods: discovered,
		LoadOrder:      loadOrder,
		Warnings:       nil,
	}

	services.Logger.Info("bootstrap completed", map[string]any{"mods": len(discovered)})

	return diagnostics, nil
}

func BootstrapRuntime(modsDir string, statesDir string, initialScene string, services RuntimeServices) (Runtime, error) {
	services = normalizeServices(services)

	diagnostics, err := Bootstrap(modsDir, services)
	if err != nil {
		return Runtime{}, err
	}

	store := newStoreForBootstrapRuntime(statesDir)
	ctx := &types.RuntimeContext{CurrentSceneID: initialScene, CurrentLocale: state.DefaultLocale}
	runtimeServices := newRuntimeServicesBridge(services.Logger, services.EventBus, store)
	runtimeLifecycle := mods.NewRuntimeLifecycle(mods.ManifestFactory{}, runtimeServices, services.Logger, services.EventBus)

	appController := controllers.NewAppController(store, ctx, services.Logger, services.EventBus)
	if _, err := appController.LoadActiveState(); err != nil {
		services.Logger.Error("bootstrap runtime failed loading active state", map[string]any{"error": err.Error()})
		return Runtime{}, err
	}

	manifests, err := mods.Discover(modsDir)
	if err != nil {
		return Runtime{}, err
	}

	loadOrder, err := mods.ResolveLoadOrder(manifests)
	if err != nil {
		return Runtime{}, err
	}

	if err := runtimeLifecycle.LoadAll(context.Background(), manifests, loadOrder, *ctx, func(progress mods.RuntimeProgress) {
		if services.Logger != nil {
			services.Logger.Info("bootstrap runtime progress", map[string]any{
				"step":  formatProgressLabel(progress),
				"phase": progress.Phase,
				"modId": progress.ModID,
				"index": progress.Index,
				"total": progress.Total,
			})
		}
	}); err != nil {
		return Runtime{}, apperrors.Wrap(apperrors.CodeDependency, apperrors.ErrorTypeDomain, "load runtime mods", err)
	}

	runtimeControllers := RuntimeControllers{
		App:      appController,
		Scene:    controllers.NewSceneControllerWithContext(initialScene, ctx),
		Learning: controllers.NewLearningController(ctx),
		Mod:      controllers.NewModController(modsDir, store, ctx, services.Logger, services.EventBus, runtimeLifecycle),
		Locale:   controllers.NewLocaleController(store, ctx, services.Logger, services.EventBus, state.SupportedLocalesSet()),
	}

	return Runtime{
		Diagnostics: diagnostics,
		Context:     ctx,
		Controllers: runtimeControllers,
		ModRegistry: runtimeLifecycle.Registry(),
		ModRuntime:  runtimeLifecycle,
	}, nil
}

func normalizeServices(services RuntimeServices) RuntimeServices {
	if services.Logger == nil {
		services.Logger = NewStdLogger()
	}
	if services.EventBus == nil {
		services.EventBus = NewInMemoryEventBus()
	}
	if services.Concurrency == nil {
		services.Concurrency = concurrency.DefaultService()
	}
	return services
}

type runtimeServicesBridge struct {
	logger     contracts.LoggingService
	eventBus   contracts.EventBusService
	stateStore contracts.StateStore
}

func newRuntimeServicesBridge(logger contracts.LoggingService, eventBus contracts.EventBusService, stateStore contracts.StateStore) contracts.RuntimeServices {
	return runtimeServicesBridge{logger: logger, eventBus: eventBus, stateStore: stateStore}
}

func (b runtimeServicesBridge) Logger() contracts.LoggingService {
	return b.logger
}

func (b runtimeServicesBridge) EventBus() contracts.EventBusService {
	return b.eventBus
}

func (b runtimeServicesBridge) StateStore() contracts.StateStore {
	return b.stateStore
}

func formatProgressLabel(progress mods.RuntimeProgress) string {
	if progress.Total <= 0 {
		return progress.ModID
	}
	return fmt.Sprintf("%d/%d %s", progress.Index, progress.Total, progress.ModID)
}
