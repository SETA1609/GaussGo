package registry

import "sync"

type Base[K comparable, V any] struct {
	mu    sync.RWMutex
	items map[K]V
}

func NewBase[K comparable, V any]() *Base[K, V] {
	return &Base[K, V]{items: map[K]V{}}
}

func (r *Base[K, V]) Register(key K, value V) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.items[key]; exists {
		return ErrDuplicateKey
	}

	r.items[key] = value
	return nil
}

func (r *Base[K, V]) Upsert(key K, value V) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[key] = value
}

func (r *Base[K, V]) Get(key K) (V, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, exists := r.items[key]
	if !exists {
		var zero V
		return zero, ErrNotFound
	}

	return item, nil
}

func (r *Base[K, V]) Delete(key K) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.items[key]; !exists {
		return ErrNotFound
	}

	delete(r.items, key)
	return nil
}

func (r *Base[K, V]) List() map[K]V {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make(map[K]V, len(r.items))
	for k, v := range r.items {
		out[k] = v
	}

	return out
}
