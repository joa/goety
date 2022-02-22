package sync

import (
	zync "sync"
)

type Map[K, V any] struct {
	m zync.Map
}

func (m *Map[K, V]) Load(key K) (value V, loaded bool) {
	if m == nil {
		return
	}

	if v, ok := m.m.Load(key); ok {
		return v.(V), ok
	}

	return
}

func (m *Map[K, V]) Store(key K, value V) {
	if m == nil {
		return
	}

	m.m.Store(key, value)
}

func (m *Map[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	if m == nil {
		return
	}

	if v, ok := m.m.LoadOrStore(key, value); ok {
		return v.(V), ok
	}

	return
}

func (m *Map[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	if m == nil {
		return
	}

	if v, ok := m.m.LoadAndDelete(key); ok {
		return v.(V), ok
	}

	return
}

func (m *Map[K, V]) Delete(key K) {
	if m == nil {
		return
	}

	m.m.Delete(key)
}

func (m *Map[K, V]) Range(f func(key K, value V) (shouldContinue bool)) {
	if m == nil {
		return
	}

	m.m.Range(func(k, v any) (shouldContinue bool) {
		return f(k.(K), v.(V))
	})
}
