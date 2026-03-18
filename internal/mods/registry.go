package mods

import (
	"sync"

	"gaussgo/internal/contracts"
	"gaussgo/internal/registry"
)

type Registry struct {
	base *registry.Base[string, contracts.ModManifest]
}

type Repository struct {
	mu       sync.RWMutex
	modsDir  string
	enabled  map[string]bool
	statuses []contracts.ModStatus
}

func NewRegistry() *Registry {
	return &Registry{base: registry.NewBase[string, contracts.ModManifest]()}
}

func (r *Registry) Register(m contracts.ModManifest) error {
	return r.base.Register(m.ID, m)
}

func (r *Registry) Upsert(m contracts.ModManifest) {
	r.base.Upsert(m.ID, m)
}

func (r *Registry) Get(id string) (contracts.ModManifest, error) {
	return r.base.Get(id)
}

func (r *Registry) Delete(id string) error {
	return r.base.Delete(id)
}

func (r *Registry) List() map[string]contracts.ModManifest {
	return r.base.List()
}

func NewRepository(modsDir string, enabled map[string]bool) *Repository {
	copyEnabled := make(map[string]bool, len(enabled))
	for id, isEnabled := range enabled {
		copyEnabled[id] = isEnabled
	}

	return &Repository{modsDir: modsDir, enabled: copyEnabled}
}

func (r *Repository) Discover() ([]contracts.ModManifest, error) {
	return Discover(r.modsDir)
}

func (r *Repository) ResolveLoadOrder(manifests []contracts.ModManifest) ([]string, error) {
	return ResolveLoadOrder(manifests)
}

func (r *Repository) Refresh() ([]contracts.ModStatus, error) {
	manifests, err := Discover(r.modsDir)
	if err != nil {
		return nil, err
	}

	statuses, err := BuildStatuses(manifests, r.enabled)
	if err != nil {
		return nil, err
	}

	// Cache the latest computed statuses for read-mostly consumers.
	r.mu.Lock()
	r.statuses = append([]contracts.ModStatus(nil), statuses...)
	r.mu.Unlock()

	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]contracts.ModStatus(nil), r.statuses...), nil
}

func (r *Repository) LastStatuses() []contracts.ModStatus {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]contracts.ModStatus(nil), r.statuses...)
}
