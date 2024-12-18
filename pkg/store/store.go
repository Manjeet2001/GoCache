package store

import (
	"sync"
)

type KeyValueStore struct {
	data sync.Map
	mu   sync.RWMutex
}

func NewKeyValueStore() *KeyValueStore {
	return &KeyValueStore{}
}

func (s *KeyValueStore) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data.Store(key, value)
}

func (s *KeyValueStore) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if value, ok := s.data.Load(key); ok {
		return value.(string), true
	}
	return "", false
}

func (s *KeyValueStore) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data.Delete(key)
}
