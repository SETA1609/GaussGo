package mods

import (
	"sync"

	"gaussgo/internal/apperrors"
	"gaussgo/internal/contracts"
)

type RuntimeRegistry struct {
	mu   sync.RWMutex
	mods map[string]contracts.ModRuntime
}

func NewRuntimeRegistry() *RuntimeRegistry {
	return &RuntimeRegistry{mods: map[string]contracts.ModRuntime{}}
}

func (r *RuntimeRegistry) Register(mod contracts.ModRuntime) error {
	if mod == nil {
		return apperrors.New(apperrors.CodeValidation, apperrors.ErrorTypeInput, "mod runtime is required")
	}
	id := mod.ID()
	if id == "" {
		return apperrors.New(apperrors.CodeValidation, apperrors.ErrorTypeInput, "mod runtime id is required")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.mods[id]; exists {
		return apperrors.New(apperrors.CodeConflict, apperrors.ErrorTypeDomain, "runtime mod already registered: "+id)
	}

	r.mods[id] = mod
	return nil
}

func (r *RuntimeRegistry) Get(modID string) (contracts.ModRuntime, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	mod, ok := r.mods[modID]
	if !ok {
		return nil, apperrors.New(apperrors.CodeNotFound, apperrors.ErrorTypeInput, "runtime mod not found: "+modID)
	}

	return mod, nil
}

func (r *RuntimeRegistry) List() map[string]contracts.ModRuntime {
	r.mu.RLock()
	defer r.mu.RUnlock()

	copyMap := make(map[string]contracts.ModRuntime, len(r.mods))
	for id, mod := range r.mods {
		copyMap[id] = mod
	}

	return copyMap
}

func (r *RuntimeRegistry) Delete(modID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.mods[modID]; !exists {
		return apperrors.New(apperrors.CodeNotFound, apperrors.ErrorTypeInput, "runtime mod not found: "+modID)
	}

	delete(r.mods, modID)
	return nil
}
