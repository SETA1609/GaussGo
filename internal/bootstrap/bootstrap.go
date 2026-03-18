package bootstrap

import (
	"fmt"

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

func Bootstrap(modsDir string) (Diagnostics, error) {
	manifests, err := mods.Discover(modsDir)
	if err != nil {
		return Diagnostics{}, fmt.Errorf("discover mods: %w", err)
	}

	loadOrder, err := mods.ResolveLoadOrder(manifests)
	if err != nil {
		return Diagnostics{}, fmt.Errorf("resolve load order: %w", err)
	}

	discovered := make([]string, 0, len(manifests))
	for _, m := range manifests {
		discovered = append(discovered, m.ID)
	}

	return Diagnostics{
		DiscoveredMods: discovered,
		LoadOrder:      loadOrder,
		Warnings:       nil,
	}, nil
}
