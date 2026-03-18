package bootstrap

import (
	"gaussgo/internal/apperrors"
	"gaussgo/internal/contracts"
	"gaussgo/internal/mods"
)

type Diagnostics struct {
	DiscoveredMods []string
	LoadOrder      []string
	Warnings       []string
}

type RuntimeServices struct {
	Logger   contracts.LoggingService
	EventBus contracts.EventBusService
}

func DefaultServices() RuntimeServices {
	return RuntimeServices{
		Logger:   NewStdLogger(),
		EventBus: NewInMemoryEventBus(),
	}
}

func Bootstrap(modsDir string, services RuntimeServices) (Diagnostics, error) {
	if services.Logger == nil {
		services.Logger = NewStdLogger()
	}
	if services.EventBus == nil {
		services.EventBus = NewInMemoryEventBus()
	}

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
