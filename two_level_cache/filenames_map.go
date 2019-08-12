package two_level_cache

import "sync"

type filenamesMap struct {
	mux       sync.RWMutex
	filenames map[string]string
}

func NewFilenamesMap() *filenamesMap {
	return &filenamesMap{filenames: make(map[string]string)}
}

func (m *filenamesMap) Get(key string) (string, bool) {
	m.mux.RLock()
	defer m.mux.RUnlock()
	filename, ok := m.filenames[key]
	return filename, ok
}

func (m *filenamesMap) Set(key string, filename string) {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.filenames[key] = filename
}

func (m *filenamesMap) Remove(key string) {
	m.mux.Lock()
	defer m.mux.Unlock()
	delete(m.filenames, key)
}

func (m *filenamesMap) Clear(key string) {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.filenames = make(map[string]string)
}
