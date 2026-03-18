package bootstrap

import (
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

	appController := controllers.NewAppController(store, ctx, services.Logger, services.EventBus)
	if _, err := appController.LoadActiveState(); err != nil {
		services.Logger.Error("bootstrap runtime failed loading active state", map[string]any{"error": err.Error()})
		return Runtime{}, err
	}

	runtimeControllers := RuntimeControllers{
		App:      appController,
		Scene:    controllers.NewSceneControllerWithContext(initialScene, ctx),
		Learning: controllers.NewLearningController(ctx),
		Mod:      controllers.NewModController(modsDir, store, ctx, services.Logger, services.EventBus),
		Locale:   controllers.NewLocaleController(store, ctx, services.Logger, services.EventBus, state.SupportedLocalesSet()),
	}

	return Runtime{
		Diagnostics: diagnostics,
		Context:     ctx,
		Controllers: runtimeControllers,
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
