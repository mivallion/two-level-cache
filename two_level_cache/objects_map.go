package two_level_cache

import "sync"

type objectsMap struct {
	mux     sync.RWMutex
	objects map[string]Marshaler
}

func NewObjectsMap() *objectsMap {
	return &objectsMap{objects: make(map[string]Marshaler)}
}

func (m *objectsMap) Get(key string) (Marshaler, bool) {
	m.mux.RLock()
	defer m.mux.RUnlock()
	obj, ok := m.objects[key]
	return obj, ok
}

func (m *objectsMap) Set(key string, obj Marshaler) {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.objects[key] = obj
}

func (m *objectsMap) Remove(key string) {
	m.mux.Lock()
	defer m.mux.Unlock()
	delete(m.objects, key)
}

func (m *objectsMap) Clear() {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.objects = make(map[string]Marshaler)
}
