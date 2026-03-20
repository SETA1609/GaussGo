package mods

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"gaussgo/internal/apperrors"
	"gaussgo/internal/contracts"
	"gaussgo/internal/types"
)

const (
	eventRuntimeInitStarted    = "mod.runtime.init.started"
	eventRuntimeInitSucceeded  = "mod.runtime.init.succeeded"
	eventRuntimeInitFailed     = "mod.runtime.init.failed"
	eventRuntimeSyncCompleted  = "mod.runtime.sync.completed"
	eventRuntimeSyncFailed     = "mod.runtime.sync.failed"
	eventRuntimeShutdownStart  = "mod.runtime.shutdown.started"
	eventRuntimeShutdownOK     = "mod.runtime.shutdown.succeeded"
	eventRuntimeShutdownFailed = "mod.runtime.shutdown.failed"
)

type RuntimeProgress struct {
	Phase string
	ModID string
	Index int
	Total int
}

type RuntimeLifecycle struct {
	mu       sync.RWMutex
	registry *RuntimeRegistry
	factory  contracts.ModFactory
	services contracts.RuntimeServices
	logger   contracts.LoggingService
	eventBus contracts.EventBusService
}

func NewRuntimeLifecycle(factory contracts.ModFactory, services contracts.RuntimeServices, logger contracts.LoggingService, eventBus contracts.EventBusService) *RuntimeLifecycle {
	return &RuntimeLifecycle{
		registry: NewRuntimeRegistry(),
		factory:  factory,
		services: services,
		logger:   logger,
		eventBus: eventBus,
	}
}

func (l *RuntimeLifecycle) Registry() contracts.RuntimeModRegistry {
	return l.registry
}

func (l *RuntimeLifecycle) LoadAll(ctx context.Context, manifests []contracts.ModManifest, loadOrder []string, runtimeCtx types.RuntimeContext, onProgress func(RuntimeProgress)) error {
	if l.factory == nil {
		return apperrors.New(apperrors.CodeValidation, apperrors.ErrorTypeDomain, "mod factory is required")
	}

	manifestByID := make(map[string]contracts.ModManifest, len(manifests))
	for _, manifest := range manifests {
		manifestByID[manifest.ID] = manifest
	}

	runtimeMods := make(map[string]contracts.ModRuntime, len(loadOrder))
	total := len(loadOrder)
	for idx, modID := range loadOrder {
		manifest, ok := manifestByID[modID]
		if !ok {
			return apperrors.New(apperrors.CodeNotFound, apperrors.ErrorTypeDomain, "manifest not found for load order mod: "+modID)
		}

		progress := RuntimeProgress{Phase: "build", ModID: modID, Index: idx + 1, Total: total}
		emitProgress(l.eventBus, eventRuntimeInitStarted, progress)
		if onProgress != nil {
			onProgress(progress)
		}

		runtimeMod, err := l.factory.Build(manifest)
		if err != nil {
			emitProgress(l.eventBus, eventRuntimeInitFailed, progress)
			return apperrors.Wrap(apperrors.CodeValidation, apperrors.ErrorTypeDomain, "build runtime mod "+modID, err)
		}
		if runtimeMod == nil {
			emitProgress(l.eventBus, eventRuntimeInitFailed, progress)
			return apperrors.New(apperrors.CodeValidation, apperrors.ErrorTypeDomain, "build runtime mod returned nil: "+modID)
		}

		runtimeMods[modID] = runtimeMod
	}

	loaded := make([]string, 0, len(loadOrder))
	for idx, modID := range loadOrder {
		runtimeMod := runtimeMods[modID]
		progress := RuntimeProgress{Phase: "init", ModID: modID, Index: idx + 1, Total: total}
		emitProgress(l.eventBus, eventRuntimeInitStarted, progress)
		if onProgress != nil {
			onProgress(progress)
		}

		initCtx := contracts.ModInitContext{RuntimeContext: runtimeCtx, Services: l.services}
		if err := runtimeMod.Init(ctx, initCtx); err != nil {
			emitProgress(l.eventBus, eventRuntimeInitFailed, progress)
			l.shutdownRegistered(ctx, loaded)
			return apperrors.Wrap(apperrors.CodeDependency, apperrors.ErrorTypeDomain, "init runtime mod "+modID, err)
		}

		if err := l.registry.Register(runtimeMod); err != nil {
			emitProgress(l.eventBus, eventRuntimeInitFailed, progress)
			_ = runtimeMod.Shutdown(ctx)
			l.shutdownRegistered(ctx, loaded)
			return err
		}

		loaded = append(loaded, modID)
		emitProgress(l.eventBus, eventRuntimeInitSucceeded, progress)
	}

	if l.logger != nil {
		l.logger.Info("runtime mods loaded", map[string]any{"count": len(loaded)})
	}

	return nil
}

func (l *RuntimeLifecycle) SyncEnabled(ctx context.Context, manifests []contracts.ModManifest, runtimeCtx types.RuntimeContext, enabled map[string]bool, onProgress func(RuntimeProgress)) error {
	allManifests, err := filterEnabledManifests(manifests, enabled)
	if err != nil {
		emitProgress(l.eventBus, eventRuntimeSyncFailed, RuntimeProgress{Phase: "sync", ModID: "", Index: 0, Total: 0})
		return err
	}

	loadOrder, err := ResolveLoadOrder(allManifests)
	if err != nil {
		emitProgress(l.eventBus, eventRuntimeSyncFailed, RuntimeProgress{Phase: "sync", ModID: "", Index: 0, Total: 0})
		return err
	}

	current := l.registry.List()
	for modID, runtimeMod := range current {
		if enabled[modID] {
			continue
		}
		shutdownProgress := RuntimeProgress{Phase: "shutdown", ModID: modID, Index: 0, Total: 0}
		emitProgress(l.eventBus, eventRuntimeShutdownStart, shutdownProgress)
		if onProgress != nil {
			onProgress(shutdownProgress)
		}

		if err := runtimeMod.Shutdown(ctx); err != nil {
			emitProgress(l.eventBus, eventRuntimeShutdownFailed, shutdownProgress)
			return apperrors.Wrap(apperrors.CodeUnknown, apperrors.ErrorTypeInfrastructure, "shutdown runtime mod "+modID, err)
		}
		if err := l.registry.Delete(modID); err != nil {
			return err
		}
		emitProgress(l.eventBus, eventRuntimeShutdownOK, shutdownProgress)
	}

	manifestByID := make(map[string]contracts.ModManifest, len(allManifests))
	for _, manifest := range allManifests {
		manifestByID[manifest.ID] = manifest
	}

	total := len(loadOrder)
	for idx, modID := range loadOrder {
		if _, err := l.registry.Get(modID); err == nil {
			continue
		}

		manifest, ok := manifestByID[modID]
		if !ok {
			return apperrors.New(apperrors.CodeNotFound, apperrors.ErrorTypeDomain, "manifest not found for enabled mod: "+modID)
		}

		progress := RuntimeProgress{Phase: "init", ModID: modID, Index: idx + 1, Total: total}
		emitProgress(l.eventBus, eventRuntimeInitStarted, progress)
		if onProgress != nil {
			onProgress(progress)
		}

		runtimeMod, buildErr := l.factory.Build(manifest)
		if buildErr != nil {
			emitProgress(l.eventBus, eventRuntimeInitFailed, progress)
			return apperrors.Wrap(apperrors.CodeValidation, apperrors.ErrorTypeDomain, "build runtime mod "+modID, buildErr)
		}
		if runtimeMod == nil {
			emitProgress(l.eventBus, eventRuntimeInitFailed, progress)
			return apperrors.New(apperrors.CodeValidation, apperrors.ErrorTypeDomain, "build runtime mod returned nil: "+modID)
		}
		if err := runtimeMod.Init(ctx, contracts.ModInitContext{RuntimeContext: runtimeCtx, Services: l.services}); err != nil {
			emitProgress(l.eventBus, eventRuntimeInitFailed, progress)
			return apperrors.Wrap(apperrors.CodeDependency, apperrors.ErrorTypeDomain, "init runtime mod "+modID, err)
		}
		if err := l.registry.Register(runtimeMod); err != nil {
			emitProgress(l.eventBus, eventRuntimeInitFailed, progress)
			_ = runtimeMod.Shutdown(ctx)
			return err
		}

		emitProgress(l.eventBus, eventRuntimeInitSucceeded, progress)
	}

	emitProgress(l.eventBus, eventRuntimeSyncCompleted, RuntimeProgress{Phase: "sync", ModID: "", Index: len(loadOrder), Total: len(loadOrder)})
	return nil
}

func (l *RuntimeLifecycle) UnloadMod(ctx context.Context, modID string) error {
	runtimeMod, err := l.registry.Get(modID)
	if err != nil {
		return err
	}

	progress := RuntimeProgress{Phase: "shutdown", ModID: modID, Index: 0, Total: 0}
	emitProgress(l.eventBus, eventRuntimeShutdownStart, progress)
	if shutdownErr := runtimeMod.Shutdown(ctx); shutdownErr != nil {
		emitProgress(l.eventBus, eventRuntimeShutdownFailed, progress)
		return apperrors.Wrap(apperrors.CodeUnknown, apperrors.ErrorTypeInfrastructure, "shutdown runtime mod "+modID, shutdownErr)
	}

	if err := l.registry.Delete(modID); err != nil {
		return err
	}

	emitProgress(l.eventBus, eventRuntimeShutdownOK, progress)
	return nil
}

func (l *RuntimeLifecycle) shutdownRegistered(ctx context.Context, loaded []string) {
	for i := len(loaded) - 1; i >= 0; i-- {
		modID := loaded[i]
		runtimeMod, err := l.registry.Get(modID)
		if err != nil {
			continue
		}
		_ = runtimeMod.Shutdown(ctx)
		_ = l.registry.Delete(modID)
	}
}

func filterEnabledManifests(manifests []contracts.ModManifest, enabled map[string]bool) ([]contracts.ModManifest, error) {
	enabledManifests := make([]contracts.ModManifest, 0, len(manifests))
	for _, manifest := range manifests {
		if enabled[manifest.ID] {
			enabledManifests = append(enabledManifests, manifest)
		}
	}

	hasCore := false
	for _, manifest := range enabledManifests {
		if manifest.ID == "core" {
			hasCore = true
			break
		}
	}
	if !hasCore {
		return nil, apperrors.New(apperrors.CodeCoreRequired, apperrors.ErrorTypeDomain, "enabled mod set must include core")
	}

	return enabledManifests, nil
}

func emitProgress(bus contracts.EventBusService, eventName string, progress RuntimeProgress) {
	if bus == nil {
		return
	}
	_ = bus.Emit(contracts.Event{
		Name:      eventName,
		Timestamp: time.Now().UTC(),
		Source:    "runtime_lifecycle",
		Data: map[string]any{
			"phase": progress.Phase,
			"modId": progress.ModID,
			"index": progress.Index,
			"total": progress.Total,
		},
	})
}

type ManifestFactory struct{}

func (ManifestFactory) Build(manifest contracts.ModManifest) (contracts.ModRuntime, error) {
	if manifest.ID == "core" {
		return &coreRuntimeMod{manifest: manifest, mathHelpers: coreBasicMathHelpers{}}, nil
	}
	return &noopRuntimeMod{manifest: manifest}, nil
}

type coreRuntimeMod struct {
	manifest    contracts.ModManifest
	mathHelpers contracts.BasicMathHelpers
	mu          sync.Mutex
	active      bool
}

func (m *coreRuntimeMod) ID() string {
	return m.manifest.ID
}

func (m *coreRuntimeMod) Manifest() contracts.ModManifest {
	return m.manifest
}

func (m *coreRuntimeMod) Init(context.Context, contracts.ModInitContext) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.active = true
	return nil
}

func (m *coreRuntimeMod) Shutdown(context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.active = false
	return nil
}

func (m *coreRuntimeMod) HasCapability(key contracts.CapabilityKey) bool {
	return key == contracts.CapabilityHelpersBasicMath
}

func (m *coreRuntimeMod) Capability(key contracts.CapabilityKey) (any, bool) {
	if key == contracts.CapabilityHelpersBasicMath {
		return m.mathHelpers, true
	}
	return nil, false
}

func (m *coreRuntimeMod) Capabilities() []contracts.CapabilityKey {
	return []contracts.CapabilityKey{contracts.CapabilityHelpersBasicMath}
}

type coreBasicMathHelpers struct{}

func (coreBasicMathHelpers) Add(ctx context.Context, float64A float64, float64B float64) (float64, error) {
	_ = ctx
	return float64A + float64B, nil
}

func (coreBasicMathHelpers) Sub(ctx context.Context, float64A float64, float64B float64) (float64, error) {
	_ = ctx
	return float64A - float64B, nil
}

func (coreBasicMathHelpers) Mul(ctx context.Context, float64A float64, float64B float64) (float64, error) {
	_ = ctx
	return float64A * float64B, nil
}

func (coreBasicMathHelpers) Div(ctx context.Context, float64A float64, float64B float64) (float64, error) {
	_ = ctx
	if math.Abs(float64B) < 1e-12 {
		return 0, apperrors.New(apperrors.CodeValidation, apperrors.ErrorTypeInput, "division by zero")
	}
	return float64A / float64B, nil
}

type noopRuntimeMod struct {
	manifest contracts.ModManifest
	mu       sync.Mutex
	active   bool
}

func (m *noopRuntimeMod) ID() string {
	return m.manifest.ID
}

func (m *noopRuntimeMod) Manifest() contracts.ModManifest {
	return m.manifest
}

func (m *noopRuntimeMod) Init(context.Context, contracts.ModInitContext) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.active = true
	return nil
}

func (m *noopRuntimeMod) Shutdown(context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.active = false
	return nil
}

func (m *noopRuntimeMod) String() string {
	return fmt.Sprintf("noopRuntimeMod(%s)", m.manifest.ID)
}
