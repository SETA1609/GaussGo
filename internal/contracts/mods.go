package contracts

import (
	"context"

	"gaussgo/internal/types"
)

type ModDependency struct {
	ID      string
	Version string
}

type ModLocalization struct {
	Default   string
	Supported []string
	Path      string
}

type ModManifest struct {
	ID           string
	Name         string
	Version      string
	Entry        string
	Localization ModLocalization
	Dependencies []ModDependency
	Provides     []string
}

type ModStatus struct {
	ModID          string
	Enabled        bool
	Available      bool
	ReasonDisabled string
}

type ModsRepository interface {
	Discover() ([]ModManifest, error)
	ResolveLoadOrder([]ModManifest) ([]string, error)
	Refresh() ([]ModStatus, error)
}

type CapabilityKey string

const (
	CapabilityHelpersBasicMath CapabilityKey = "helpers.basic_math"
)

type HelperOperationInput map[string]any
type HelperOperationOutput map[string]any

type ModInitContext struct {
	RuntimeContext types.RuntimeContext
	Services       RuntimeServices
}

type RuntimeServices interface {
	Logger() LoggingService
	EventBus() EventBusService
	StateStore() StateStore
}

type ModRuntime interface {
	ID() string
	Manifest() ModManifest
	Init(ctx context.Context, initCtx ModInitContext) error
	Shutdown(ctx context.Context) error
}

type CapabilityProvider interface {
	HasCapability(key CapabilityKey) bool
	Capability(key CapabilityKey) (any, bool)
}

type CapabilityCatalog interface {
	Capabilities() []CapabilityKey
}

type BasicMathHelpers interface {
	Add(ctx context.Context, a float64, b float64) (float64, error)
	Sub(ctx context.Context, a float64, b float64) (float64, error)
	Mul(ctx context.Context, a float64, b float64) (float64, error)
	Div(ctx context.Context, a float64, b float64) (float64, error)
}

type ModFactory interface {
	Build(manifest ModManifest) (ModRuntime, error)
}

type RuntimeModRegistry interface {
	Register(mod ModRuntime) error
	Get(modID string) (ModRuntime, error)
	List() map[string]ModRuntime
}
