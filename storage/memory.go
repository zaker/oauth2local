package storage

import "sync"

type MemoryStorage struct {
	rw   *sync.RWMutex
	code string
}

func Memory() *MemoryStorage {
	return &MemoryStorage{rw: new(sync.RWMutex)}
}

func (m *MemoryStorage) GetCode() (string, error) {
	m.rw.RLock()
	defer m.rw.RUnlock()
	return m.code, nil
}

func (m *MemoryStorage) SetCode(code string) error {
	m.rw.Lock()
	defer m.rw.Unlock()
	m.code = code
	return nil
}
