package store

import (
	"errors"
	"sync"
)

var ErrKeyNotFound = errors.New("key not found")

type Store struct {
	mu    sync.RWMutex
	store map[string]any
}

func New() *Store {
	return &Store{
		store: make(map[string]any),
	}
}

// Get retrieves a value from the store.
func Get[T any](key string, s *Store) (T, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	values, ok := s.store[key]
	if !ok {
		return *new(T), ErrKeyNotFound
	}

	return values.(T), nil
}

// Set stores a value in the store.
func Set[T any](key string, value T, s *Store) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.store[key] = value
}
