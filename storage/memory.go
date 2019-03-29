package storage

import "sync"

type Memory struct {
	rw   *sync.RWMutex
	code string
}

func (m *Memory) GetCode() (string, error) {
	m.rw.RLock()
	defer m.rw.RUnlock()
	return m.code, nil
}

func (m *Memory) SetCode(code string) error {
	m.rw.Lock()
	defer m.rw.Unlock()
	m.code = code
	return nil
}
